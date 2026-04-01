package scanner

import (
	"fmt"
	"net"
	"time"

	"lansentry/device"
)

// Scanner defines the interface for network scanning.
type Scanner interface {
	Scan() ([]device.Device, error)
}

// ARPScanner scans the local network using ARP.
type ARPScanner struct {
	interfaceName string
	timeout       time.Duration
}

// NewARPScanner creates a new ARP scanner.
func NewARPScanner(interfaceName string) (*ARPScanner, error) {
	if interfaceName == "" {
		// Auto-detect interface
		iface, err := getDefaultInterface()
		if err != nil {
			return nil, fmt.Errorf("failed to detect network interface: %w", err)
		}
		interfaceName = iface.Name
	}

	return &ARPScanner{
		interfaceName: interfaceName,
		timeout:       10 * time.Second,
	}, nil
}

// Interface returns the network interface name being used.
func (s *ARPScanner) Interface() string {
	return s.interfaceName
}

// getDefaultInterface returns the first suitable network interface.
func getDefaultInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// Skip interfaces without hardware address
		if len(iface.HardwareAddr) == 0 {
			continue
		}

		// Check if interface has IPv4 address
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() {
					return &iface, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no suitable network interface found")
}

// getLocalIPAndMask returns the local IP and network mask for the interface.
func (s *ARPScanner) getLocalIPAndMask() (*net.IPNet, net.IP, error) {
	iface, err := net.InterfaceByName(s.interfaceName)
	if err != nil {
		return nil, nil, fmt.Errorf("interface %s not found: %w", s.interfaceName, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil && !ip4.IsLoopback() {
				return ipnet, ip4, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("no IPv4 address found on interface %s", s.interfaceName)
}

// getHardwareAddr returns the hardware address of the interface.
func (s *ARPScanner) getHardwareAddr() (net.HardwareAddr, error) {
	iface, err := net.InterfaceByName(s.interfaceName)
	if err != nil {
		return nil, err
	}
	return iface.HardwareAddr, nil
}

// generateIPRange generates all IPs in the subnet.
func generateIPRange(ipnet *net.IPNet) []net.IP {
	var ips []net.IP
	ip := ipnet.IP.Mask(ipnet.Mask)

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		// Skip network address and broadcast address
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		ips = append(ips, ipCopy)
	}

	// Remove first (network) and last (broadcast) addresses
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}

	return ips
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
