package notifier

import (
	"strings"
	"testing"
	"time"

	"lansentry/device"
)

func TestEscapeMarkdown(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"hello_world", "hello\\_world"},
		{"*bold*", "\\*bold\\*"},
		{"`code`", "\\`code\\`"},
		{"[link]", "\\[link]"},
		{"no-escape.here!", "no-escape.here!"},
		{"mix_of*chars`and[brackets", "mix\\_of\\*chars\\`and\\[brackets"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := escapeMarkdown(tt.input)
			if got != tt.want {
				t.Errorf("escapeMarkdown(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		dur  time.Duration
		want string
	}{
		{"zero", 0, "0m"},
		{"minutes only", 45 * time.Minute, "45m"},
		{"hours and minutes", 3*time.Hour + 15*time.Minute, "3h 15m"},
		{"days hours minutes", 2*24*time.Hour + 5*time.Hour + 30*time.Minute, "2d 5h 30m"},
		{"exact days", 7 * 24 * time.Hour, "7d"},
		{"exact hours", 2 * time.Hour, "2h"},
		{"1 minute", time.Minute, "1m"},
		{"sub-minute rounds to 0m", 30 * time.Second, "0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.dur)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.dur, got, tt.want)
			}
		})
	}
}

func TestFormatNewDeviceMessage(t *testing.T) {
	tn := &TelegramNotifier{chatID: 123} // bot is nil, fine for formatting
	d := device.Device{
		MAC:          "aa:bb:cc:dd:ee:ff",
		IP:           "192.168.1.10",
		Hostname:     "my-phone",
		Manufacturer: "Apple",
		FirstSeen:    time.Date(2026, 3, 31, 14, 30, 0, 0, time.UTC),
		TimesSeen:    1,
	}

	msg := tn.formatNewDeviceMessage(d)

	checks := []string{
		"New Device Joined",
		"192.168.1.10",
		"aa:bb:cc:dd:ee:ff",
		"Apple",
		"my-phone",
		"PM",
	}
	for _, sub := range checks {
		if !strings.Contains(msg, sub) {
			t.Errorf("message missing %q:\n%s", sub, msg)
		}
	}
}

func TestFormatRejoinMessage(t *testing.T) {
	tn := &TelegramNotifier{chatID: 123}
	d := device.Device{
		MAC:          "aa:bb:cc:dd:ee:ff",
		IP:           "192.168.1.10",
		Hostname:     "laptop",
		Manufacturer: "Dell",
		FirstSeen:    time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC),
		TimesSeen:    5,
	}

	msg := tn.formatRejoinMessage(d, 3*24*time.Hour+2*time.Hour)

	checks := []string{
		"Rejoined",
		"192.168.1.10",
		"aa:bb:cc:dd:ee:ff",
		"Dell",
		"laptop",
		"3d 2h",
		"Times Seen:* 5",
	}
	for _, sub := range checks {
		if !strings.Contains(msg, sub) {
			t.Errorf("message missing %q:\n%s", sub, msg)
		}
	}
}

func TestFormatNewDeviceMessage_NoOptionalFields(t *testing.T) {
	tn := &TelegramNotifier{chatID: 123}
	d := device.Device{
		MAC:       "aa:bb:cc:dd:ee:ff",
		IP:        "192.168.1.10",
		FirstSeen: time.Now(),
	}

	msg := tn.formatNewDeviceMessage(d)

	if strings.Contains(msg, "Manufacturer") {
		t.Error("should not contain Manufacturer when empty")
	}
	if strings.Contains(msg, "Hostname") {
		t.Error("should not contain Hostname when empty")
	}
}

