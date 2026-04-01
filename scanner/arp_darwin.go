package scanner

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"lansentry/device"
)

// Scan performs an ARP scan on the local network.
// Uses arp-scan or falls back to arp table parsing.
func (s *ARPScanner) Scan() ([]device.Device, error) {
	// Try using arp-scan first (most reliable)
	devices, err := s.scanWithArpScan()
	if err == nil {
		return devices, nil
	}

	// Fall back to ping sweep + arp table
	return s.scanWithPingAndArp()
}

// scanWithArpScan uses the arp-scan tool if available.
func (s *ARPScanner) scanWithArpScan() ([]device.Device, error) {
	// Check if arp-scan is available
	_, err := exec.LookPath("arp-scan")
	if err != nil {
		return nil, fmt.Errorf("arp-scan not found")
	}

	cmd := exec.Command("arp-scan", "-I", s.interfaceName, "-l", "-q")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("arp-scan failed: %w", err)
	}

	return parseArpScanOutput(string(output))
}

// parseArpScanOutput parses the output of arp-scan.
func parseArpScanOutput(output string) ([]device.Device, error) {
	var devices []device.Device
	scanner := bufio.NewScanner(strings.NewReader(output))

	// Match lines like: 192.168.1.1	00:11:22:33:44:55	Manufacturer
	re := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+)\s+([0-9a-fA-F:]+)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			devices = append(devices, device.Device{
				IP:  matches[1],
				MAC: device.NormalizeMac(matches[2]),
			})
		}
	}

	return devices, nil
}

// scanWithPingAndArp performs a ping sweep then reads the ARP table.
func (s *ARPScanner) scanWithPingAndArp() ([]device.Device, error) {
	ipnet, localIP, err := s.getLocalIPAndMask()
	if err != nil {
		return nil, err
	}

	// Generate IP range
	ips := generateIPRange(ipnet)

	// Limit concurrent pings
	const maxConcurrent = 50
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	// Ping all IPs to populate ARP table
	for _, ip := range ips {
		if ip.Equal(localIP) {
			continue // Skip self
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(targetIP net.IP) {
			defer wg.Done()
			defer func() { <-sem }()
			pingHost(targetIP.String())
		}(ip)
	}

	wg.Wait()

	// Wait a moment for ARP table to populate
	time.Sleep(500 * time.Millisecond)

	// Read ARP table
	return s.readArpTable()
}

// pingHost sends a single ping to the host.
func pingHost(ip string) {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	_ = cmd.Run() // Ignore errors - we just want to populate ARP cache
}

// readArpTable reads the system ARP table.
func (s *ARPScanner) readArpTable() ([]device.Device, error) {
	cmd := exec.Command("arp", "-an")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read ARP table: %w", err)
	}

	return parseArpOutput(string(output), s.interfaceName)
}

// parseArpOutput parses the output of `arp -an`.
func parseArpOutput(output, interfaceName string) ([]device.Device, error) {
	var devices []device.Device
	scanner := bufio.NewScanner(strings.NewReader(output))

	// Match lines like:
	//   ? (192.168.1.1) at 00:11:22:33:44:55 on en0 ifscope [ethernet]
	//   myhost.local (192.168.1.5) at 00:11:22:33:44:55 on en0 ifscope [ethernet]
	re := regexp.MustCompile(`^(\S+)\s+\((\d+\.\d+\.\d+\.\d+)\)\s+at\s+([0-9a-fA-F:]+)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip incomplete entries
		if strings.Contains(line, "(incomplete)") {
			continue
		}

		// Filter by interface if specified
		if interfaceName != "" && !strings.Contains(line, interfaceName) {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 4 {
			hostname := matches[1]
			ip := matches[2]
			mac := matches[3]

			// "?" means no hostname known
			if hostname == "?" {
				hostname = ""
			}

			// Skip broadcast and invalid MACs
			if mac == "ff:ff:ff:ff:ff:ff" || mac == "(incomplete)" {
				continue
			}

			// Skip multicast IPs (224.0.0.0 - 239.255.255.255)
			if isMulticastIP(ip) {
				continue
			}

			// Skip multicast MACs (first octet has LSB set to 1)
			if isMulticastMAC(mac) {
				continue
			}

			devices = append(devices, device.Device{
				IP:       ip,
				MAC:      device.NormalizeMac(mac),
				Hostname: hostname,
			})
		}
	}

	return devices, nil
}

// isMulticastIP checks if an IP is a multicast address (224.0.0.0 - 239.255.255.255).
func isMulticastIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	firstOctet := 0
	fmt.Sscanf(parts[0], "%d", &firstOctet)
	return firstOctet >= 224 && firstOctet <= 239
}

// isMulticastMAC checks if a MAC address is multicast (first octet LSB is 1).
func isMulticastMAC(mac string) bool {
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	if len(mac) < 2 {
		return false
	}
	// Parse first octet
	var firstOctet uint8
	fmt.Sscanf(mac[:2], "%02x", &firstOctet)
	return (firstOctet & 0x01) != 0
}
