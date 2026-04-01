package resolver

import (
	"testing"
)

func TestCleanHostname(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		ip   string
		want string
	}{
		{"trailing dot", "myhost.lan.", "10.0.0.1", "myhost"},
		{"strip .lan", "myphone.lan", "10.0.0.2", "myphone"},
		{"strip .local", "MacBook-Pro.local", "10.0.0.3", "MacBook-Pro"},
		{"strip .home", "laptop.home", "10.0.0.4", "laptop"},
		{"strip .localdomain", "server.localdomain", "10.0.0.5", "server"},
		{"strip .internal", "nas.internal", "10.0.0.6", "nas"},
		{"strip .private", "cam.private", "10.0.0.7", "cam"},
		{"strip .fritz.box", "phone.fritz.box", "10.0.0.8", "phone"},
		{"strip .home.arpa", "printer.home.arpa", "10.0.0.9", "printer"},
		{"no suffix", "pi.hole", "10.0.0.10", "pi.hole"},
		{"empty after strip", ".lan", "10.0.0.11", ".lan"},
		{"returns empty for ip match", "10.0.0.1", "10.0.0.1", ""},
		{"empty raw", "", "10.0.0.1", ""},
		{"dot only", ".", "10.0.0.1", ""},
		{"case insensitive suffix", "MyHost.LAN", "10.0.0.1", "MyHost"},
		{"preserves case in name", "MyLaptop.local", "10.0.0.1", "MyLaptop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanHostname(tt.raw, tt.ip)
			if got != tt.want {
				t.Errorf("cleanHostname(%q, %q) = %q, want %q", tt.raw, tt.ip, got, tt.want)
			}
		})
	}
}

func TestGatewayRegex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"macOS netstat default route",
			"default            192.168.10.1       UGScg             en0",
			"192.168.10.1",
		},
		{
			"linux netstat default route",
			"0.0.0.0         192.168.1.1     0.0.0.0         UG    100    0        0 eth0",
			"192.168.1.1",
		},
		{
			"no match",
			"192.168.1.0     *               255.255.255.0   U     0      0        0 eth0",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := gwRe.FindStringSubmatch(tt.input)
			got := ""
			if len(m) >= 2 {
				got = m[1]
			}
			if got != tt.want {
				t.Errorf("gwRe match = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNetBIOSRegex(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		match bool // whether nbRe matches at all
		want  string
	}{
		{
			"unique entry",
			"\tMYPC            <00> -         B <ACTIVE>",
			true,
			"MYPC",
		},
		{
			"group entry matches regex but code skips via GROUP check",
			"\tWORKGROUP       <00> - <GROUP> B <ACTIVE>",
			false, // regex won't match because <GROUP> breaks the pattern
			"",
		},
		{
			"no match",
			"Looking up status of 192.168.1.5",
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := nbRe.FindStringSubmatch(tt.line)
			got := ""
			if len(m) >= 2 {
				got = m[1]
			}
			if got != tt.want {
				t.Errorf("nbRe match = %q, want %q", got, tt.want)
			}
		})
	}
}

