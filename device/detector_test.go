package device

import (
	"testing"
	"time"
)

func TestDetectEvents_NewDevice(t *testing.T) {
	d := NewDetector(7)
	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10"}}
	known := map[string]Device{}

	events := d.DetectEvents(scanned, known)
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
	if events[0].Type != EventNew {
		t.Errorf("event type = %q, want %q", events[0].Type, EventNew)
	}
	if events[0].Device.TimesSeen != 1 {
		t.Errorf("TimesSeen = %d, want 1", events[0].Device.TimesSeen)
	}
}

func TestDetectEvents_RejoinAfterThreshold(t *testing.T) {
	d := NewDetector(7)
	now := time.Now()
	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10"}}
	known := map[string]Device{
		"aa:bb:cc:dd:ee:ff": {
			MAC:       "aa:bb:cc:dd:ee:ff",
			IP:        "192.168.1.10",
			FirstSeen: now.Add(-30 * 24 * time.Hour),
			LastSeen:  now.Add(-8 * 24 * time.Hour), // 8 days ago > 7 day threshold
			TimesSeen: 5,
		},
	}

	events := d.DetectEvents(scanned, known)
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
	if events[0].Type != EventRejoin {
		t.Errorf("event type = %q, want %q", events[0].Type, EventRejoin)
	}
	if events[0].Device.TimesSeen != 6 {
		t.Errorf("TimesSeen = %d, want 6", events[0].Device.TimesSeen)
	}
	if events[0].OfflineFor < 7*24*time.Hour {
		t.Errorf("OfflineFor = %v, want >= 7 days", events[0].OfflineFor)
	}
}

func TestDetectEvents_NoEventWhenRecent(t *testing.T) {
	d := NewDetector(7)
	now := time.Now()
	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10"}}
	known := map[string]Device{
		"aa:bb:cc:dd:ee:ff": {
			MAC:       "aa:bb:cc:dd:ee:ff",
			IP:        "192.168.1.10",
			FirstSeen: now.Add(-24 * time.Hour),
			LastSeen:  now.Add(-1 * time.Hour), // 1 hour ago — well under threshold
			TimesSeen: 3,
		},
	}

	events := d.DetectEvents(scanned, known)
	if len(events) != 0 {
		t.Errorf("got %d events, want 0 (device seen recently)", len(events))
	}
}

func TestDetectEvents_ExactlyAtThreshold(t *testing.T) {
	d := NewDetector(7)
	now := time.Now()
	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10"}}
	known := map[string]Device{
		"aa:bb:cc:dd:ee:ff": {
			MAC:       "aa:bb:cc:dd:ee:ff",
			IP:        "192.168.1.10",
			FirstSeen: now.Add(-14 * 24 * time.Hour),
			LastSeen:  now.Add(-7 * 24 * time.Hour), // exactly at threshold
			TimesSeen: 2,
		},
	}

	events := d.DetectEvents(scanned, known)
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1 (>= threshold should trigger)", len(events))
	}
	if events[0].Type != EventRejoin {
		t.Errorf("event type = %q, want %q", events[0].Type, EventRejoin)
	}
}

func TestDetectEvents_JustUnderThreshold(t *testing.T) {
	d := NewDetector(7)
	now := time.Now()
	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10"}}
	known := map[string]Device{
		"aa:bb:cc:dd:ee:ff": {
			MAC:       "aa:bb:cc:dd:ee:ff",
			IP:        "192.168.1.10",
			FirstSeen: now.Add(-14 * 24 * time.Hour),
			LastSeen:  now.Add(-7*24*time.Hour + time.Minute), // just under threshold
			TimesSeen: 2,
		},
	}

	events := d.DetectEvents(scanned, known)
	if len(events) != 0 {
		t.Errorf("got %d events, want 0 (just under threshold)", len(events))
	}
}

func TestSetRejoinThreshold(t *testing.T) {
	d := NewDetector(1)
	d.SetRejoinThreshold(14)

	now := time.Now()
	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10"}}
	known := map[string]Device{
		"aa:bb:cc:dd:ee:ff": {
			MAC:      "aa:bb:cc:dd:ee:ff",
			IP:       "192.168.1.10",
			LastSeen: now.Add(-10 * 24 * time.Hour), // 10 days < new 14 day threshold
		},
	}

	events := d.DetectEvents(scanned, known)
	if len(events) != 0 {
		t.Errorf("got %d events, want 0 (10d < 14d threshold)", len(events))
	}
}

func TestUpdateDevices_NewDevice(t *testing.T) {
	d := NewDetector(7)
	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10", Hostname: "phone"}}
	known := map[string]Device{}

	updates := d.UpdateDevices(scanned, known)
	if len(updates) != 1 {
		t.Fatalf("got %d updates, want 1", len(updates))
	}
	if updates[0].TimesSeen != 1 {
		t.Errorf("TimesSeen = %d, want 1", updates[0].TimesSeen)
	}
	if updates[0].Hostname != "phone" {
		t.Errorf("Hostname = %q, want %q", updates[0].Hostname, "phone")
	}
}

func TestUpdateDevices_ExistingDevice_MergesFields(t *testing.T) {
	d := NewDetector(7)
	now := time.Now()
	firstSeen := now.Add(-48 * time.Hour)

	scanned := []Device{{MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.20", Hostname: "", Manufacturer: "NewMfr"}}
	known := map[string]Device{
		"aa:bb:cc:dd:ee:ff": {
			MAC:          "aa:bb:cc:dd:ee:ff",
			IP:           "192.168.1.10",
			Hostname:     "old-host",
			Manufacturer: "OldMfr",
			FirstSeen:    firstSeen,
			LastSeen:     now.Add(-1 * time.Hour),
			TimesSeen:    5,
		},
	}

	updates := d.UpdateDevices(scanned, known)
	if len(updates) != 1 {
		t.Fatalf("got %d updates, want 1", len(updates))
	}
	u := updates[0]
	if u.TimesSeen != 6 {
		t.Errorf("TimesSeen = %d, want 6", u.TimesSeen)
	}
	// Hostname empty in scanned → keep known
	if u.Hostname != "old-host" {
		t.Errorf("Hostname = %q, want %q (should keep known)", u.Hostname, "old-host")
	}
	// Manufacturer non-empty in scanned → use scanned
	if u.Manufacturer != "NewMfr" {
		t.Errorf("Manufacturer = %q, want %q (should use scanned)", u.Manufacturer, "NewMfr")
	}
	// IP should be updated
	if u.IP != "192.168.1.20" {
		t.Errorf("IP = %q, want %q", u.IP, "192.168.1.20")
	}
	// FirstSeen should be preserved
	if !u.FirstSeen.Equal(firstSeen) {
		t.Errorf("FirstSeen changed, want preserved")
	}
}

func TestMultipleDevices(t *testing.T) {
	d := NewDetector(7)
	now := time.Now()

	scanned := []Device{
		{MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.1.1"},  // new
		{MAC: "aa:bb:cc:dd:ee:02", IP: "192.168.1.2"},  // rejoin
		{MAC: "aa:bb:cc:dd:ee:03", IP: "192.168.1.3"},  // recent
	}
	known := map[string]Device{
		"aa:bb:cc:dd:ee:02": {MAC: "aa:bb:cc:dd:ee:02", LastSeen: now.Add(-10 * 24 * time.Hour), TimesSeen: 1},
		"aa:bb:cc:dd:ee:03": {MAC: "aa:bb:cc:dd:ee:03", LastSeen: now.Add(-1 * time.Hour), TimesSeen: 2},
	}

	events := d.DetectEvents(scanned, known)

	var newCount, rejoinCount int
	for _, e := range events {
		switch e.Type {
		case EventNew:
			newCount++
		case EventRejoin:
			rejoinCount++
		}
	}
	if newCount != 1 {
		t.Errorf("new events = %d, want 1", newCount)
	}
	if rejoinCount != 1 {
		t.Errorf("rejoin events = %d, want 1", rejoinCount)
	}
}

