package device

import (
	"strings"
	"time"
)

// Device represents a network device detected on the local network.
type Device struct {
	MAC          string    `json:"mac"`
	IP           string    `json:"ip"`
	Hostname     string    `json:"hostname,omitempty"`
	Manufacturer string    `json:"manufacturer,omitempty"`
	FirstSeen    time.Time `json:"first_seen"`
	LastSeen     time.Time `json:"last_seen"`
	TimesSeen    int       `json:"times_seen"`
}

// EventType represents the type of device event.
type EventType string

const (
	EventNew    EventType = "new"
	EventRejoin EventType = "rejoin"
)

// DeviceEvent represents a device state change event.
type DeviceEvent struct {
	Device     Device
	Type       EventType
	OfflineFor time.Duration // Only set for rejoin events
}

// NormalizeMac ensures MAC address is in a consistent format (lowercase, colon-separated).
func NormalizeMac(mac string) string {
	// Convert to lowercase and ensure colon separators
	result := ""
	mac = strings.ToLower(mac)

	// Remove any existing separators
	cleaned := ""
	for _, c := range mac {
		if c != ':' && c != '-' && c != '.' {
			cleaned += string(c)
		}
	}

	// Add colons every 2 characters
	for i, c := range cleaned {
		if i > 0 && i%2 == 0 {
			result += ":"
		}
		result += string(c)
	}

	return result
}
