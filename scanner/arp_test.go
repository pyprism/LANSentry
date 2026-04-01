package scanner

import (
	"net"
	"testing"
)

func TestGenerateIPRange_Slash24(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.1.0/24")
	ips := generateIPRange(ipnet)

	// /24 has 256 addresses, minus network and broadcast = 254 hosts
	if len(ips) != 254 {
		t.Errorf("got %d IPs, want 254", len(ips))
	}

	first := ips[0].String()
	last := ips[len(ips)-1].String()
	if first != "192.168.1.1" {
		t.Errorf("first IP = %s, want 192.168.1.1", first)
	}
	if last != "192.168.1.254" {
		t.Errorf("last IP = %s, want 192.168.1.254", last)
	}
}

func TestGenerateIPRange_Slash30(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/30")
	ips := generateIPRange(ipnet)

	// /30 has 4 addresses, minus network and broadcast = 2 hosts
	if len(ips) != 2 {
		t.Errorf("got %d IPs, want 2", len(ips))
	}
}

func TestGenerateIPRange_Slash31(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/31")
	ips := generateIPRange(ipnet)

	// /31 has 2 addresses — no network/broadcast by convention
	// With current code logic (remove first and last if >2), len(ips) would be 2 and not stripped
	if len(ips) != 2 {
		t.Errorf("got %d IPs, want 2 (point-to-point)", len(ips))
	}
}

func TestIncrementIP(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want string
	}{
		{"simple", net.ParseIP("192.168.1.1").To4(), "192.168.1.2"},
		{"rollover last octet", net.ParseIP("192.168.1.255").To4(), "192.168.2.0"},
		{"rollover two octets", net.ParseIP("192.168.255.255").To4(), "192.169.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incrementIP(tt.ip)
			if tt.ip.String() != tt.want {
				t.Errorf("after increment: %s, want %s", tt.ip, tt.want)
			}
		})
	}
}

func TestParseArpScanOutput(t *testing.T) {
	output := `Interface: en0, type: EN10MB, MAC: aa:bb:cc:dd:ee:ff, IPv4: 192.168.1.100
Starting arp-scan 1.10.0 with 256 hosts
192.168.1.1	00:11:22:33:44:55	Vendor A
192.168.1.2	AA:BB:CC:DD:EE:FF	Vendor B

3 packets received by filter, 0 packets dropped by kernel
Ending arp-scan: 256 hosts scanned in 1.234 seconds
`
	devices, err := parseArpScanOutput(output)
	if err != nil {
		t.Fatalf("parseArpScanOutput error: %v", err)
	}
	if len(devices) != 2 {
		t.Fatalf("got %d devices, want 2", len(devices))
	}
	if devices[0].IP != "192.168.1.1" {
		t.Errorf("device[0].IP = %q, want 192.168.1.1", devices[0].IP)
	}
	if devices[1].MAC != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("device[1].MAC = %q, want aa:bb:cc:dd:ee:ff", devices[1].MAC)
	}
}

func TestParseArpOutput(t *testing.T) {
	output := `? (192.168.1.1) at 00:11:22:33:44:55 on en0 ifscope [ethernet]
myhost.local (192.168.1.5) at aa:bb:cc:dd:ee:ff on en0 ifscope [ethernet]
? (192.168.1.99) at (incomplete) on en0 ifscope [ethernet]
? (192.168.1.3) at 00:11:22:33:44:66 on en1 ifscope [ethernet]
? (224.0.0.1) at 01:00:5e:00:00:01 on en0 ifscope [ethernet]
? (192.168.1.255) at ff:ff:ff:ff:ff:ff on en0 ifscope [ethernet]
`

	devices, err := parseArpOutput(output, "en0")
	if err != nil {
		t.Fatalf("parseArpOutput error: %v", err)
	}

	// Should include: 192.168.1.1 and 192.168.1.5
	// Should exclude: incomplete, en1 (wrong iface), multicast IP, broadcast MAC
	if len(devices) != 2 {
		t.Fatalf("got %d devices, want 2", len(devices))
	}

	// First device: no hostname (was "?")
	if devices[0].Hostname != "" {
		t.Errorf("device[0].Hostname = %q, want empty", devices[0].Hostname)
	}
	if devices[0].IP != "192.168.1.1" {
		t.Errorf("device[0].IP = %q, want 192.168.1.1", devices[0].IP)
	}

	// Second device: has hostname
	if devices[1].Hostname != "myhost.local" {
		t.Errorf("device[1].Hostname = %q, want myhost.local", devices[1].Hostname)
	}
}

func TestParseArpOutput_NoInterfaceFilter(t *testing.T) {
	output := `? (192.168.1.1) at 00:11:22:33:44:55 on en0 ifscope [ethernet]
? (192.168.1.3) at 00:11:22:33:44:66 on en1 ifscope [ethernet]
`
	devices, err := parseArpOutput(output, "")
	if err != nil {
		t.Fatalf("parseArpOutput error: %v", err)
	}
	if len(devices) != 2 {
		t.Errorf("got %d devices, want 2 (no interface filter)", len(devices))
	}
}

func TestIsMulticastIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"224.0.0.1", true},
		{"239.255.255.250", true},
		{"192.168.1.1", false},
		{"223.255.255.255", false},
		{"240.0.0.1", false},
		{"invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			if got := isMulticastIP(tt.ip); got != tt.want {
				t.Errorf("isMulticastIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsMulticastMAC(t *testing.T) {
	tests := []struct {
		mac  string
		want bool
	}{
		{"01:00:5e:00:00:01", true},  // multicast (LSB of first octet = 1)
		{"ff:ff:ff:ff:ff:ff", true},  // broadcast is also multicast
		{"00:11:22:33:44:55", false}, // unicast
		{"02:00:00:00:00:00", false}, // locally administered but unicast
		{"", false},
		{"x", false},
	}
	for _, tt := range tests {
		t.Run(tt.mac, func(t *testing.T) {
			if got := isMulticastMAC(tt.mac); got != tt.want {
				t.Errorf("isMulticastMAC(%q) = %v, want %v", tt.mac, got, tt.want)
			}
		})
	}
}

