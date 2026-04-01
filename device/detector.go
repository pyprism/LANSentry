package device

import (
	"time"
)

// Detector handles device state changes and event detection.
type Detector struct {
	rejoinThreshold time.Duration
}

// NewDetector creates a new device detector with the specified rejoin threshold.
func NewDetector(rejoinThresholdDays int) *Detector {
	return &Detector{
		rejoinThreshold: time.Duration(rejoinThresholdDays) * 24 * time.Hour,
	}
}

// SetRejoinThreshold updates the rejoin threshold.
func (d *Detector) SetRejoinThreshold(days int) {
	d.rejoinThreshold = time.Duration(days) * 24 * time.Hour
}

// DetectEvents compares scanned devices with known devices and returns events.
// knownDevices is a map of MAC address to Device.
func (d *Detector) DetectEvents(scannedDevices []Device, knownDevices map[string]Device) []DeviceEvent {
	var events []DeviceEvent
	now := time.Now()

	for _, scanned := range scannedDevices {
		mac := NormalizeMac(scanned.MAC)
		known, exists := knownDevices[mac]

		if !exists {
			// New device - never seen before
			events = append(events, DeviceEvent{
				Device: Device{
					MAC:          mac,
					IP:           scanned.IP,
					Hostname:     scanned.Hostname,
					Manufacturer: scanned.Manufacturer,
					FirstSeen:    now,
					LastSeen:     now,
					TimesSeen:    1,
				},
				Type: EventNew,
			})
		} else {
			// Known device - check if it's a rejoin
			offlineFor := now.Sub(known.LastSeen)

			if offlineFor >= d.rejoinThreshold {
				// Device has been offline longer than threshold - rejoin event
				events = append(events, DeviceEvent{
					Device: Device{
						MAC:          mac,
						IP:           scanned.IP,
						Hostname:     scanned.Hostname,
						Manufacturer: scanned.Manufacturer,
						FirstSeen:    known.FirstSeen,
						LastSeen:     now,
						TimesSeen:    known.TimesSeen + 1,
					},
					Type:       EventRejoin,
					OfflineFor: offlineFor,
				})
			}
			// If device was seen recently, no event is generated (just update last_seen)
		}
	}

	return events
}

// UpdateDevices returns a list of all devices that need their last_seen updated.
// This includes both new devices and existing devices that were seen again.
func (d *Detector) UpdateDevices(scannedDevices []Device, knownDevices map[string]Device) []Device {
	var updates []Device
	now := time.Now()

	for _, scanned := range scannedDevices {
		mac := NormalizeMac(scanned.MAC)
		known, exists := knownDevices[mac]

		if !exists {
			// New device
			updates = append(updates, Device{
				MAC:          mac,
				IP:           scanned.IP,
				Hostname:     scanned.Hostname,
				Manufacturer: scanned.Manufacturer,
				FirstSeen:    now,
				LastSeen:     now,
				TimesSeen:    1,
			})
		} else {
			// Update existing device
			updates = append(updates, Device{
				MAC:          mac,
				IP:           scanned.IP,
				Hostname:     ternaryStr(scanned.Hostname != "", scanned.Hostname, known.Hostname),
				Manufacturer: ternaryStr(scanned.Manufacturer != "", scanned.Manufacturer, known.Manufacturer),
				FirstSeen:    known.FirstSeen,
				LastSeen:     now,
				TimesSeen:    known.TimesSeen + 1,
			})
		}
	}

	return updates
}

func ternaryStr(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
