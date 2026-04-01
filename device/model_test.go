package device

import "testing"

func TestNormalizeMac(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"colon-separated lowercase", "aa:bb:cc:dd:ee:ff", "aa:bb:cc:dd:ee:ff"},
		{"colon-separated uppercase", "AA:BB:CC:DD:EE:FF", "aa:bb:cc:dd:ee:ff"},
		{"dash-separated", "AA-BB-CC-DD-EE-FF", "aa:bb:cc:dd:ee:ff"},
		{"dot-separated cisco", "aabb.ccdd.eeff", "aa:bb:cc:dd:ee:ff"},
		{"no separators", "aabbccddeeff", "aa:bb:cc:dd:ee:ff"},
		{"mixed case colons", "Aa:Bb:Cc:Dd:Ee:Ff", "aa:bb:cc:dd:ee:ff"},
		{"empty", "", ""},
		{"short", "aa:bb", "aa:bb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeMac(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeMac(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

