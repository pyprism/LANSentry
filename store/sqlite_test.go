package store

import (
	"path/filepath"
	"testing"
	"time"

	"lansentry/config"
	"lansentry/device"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, err := New(dbPath)
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// --- Device CRUD ---

func TestUpsertAndGetDevice(t *testing.T) {
	s := newTestStore(t)
	now := time.Now().Truncate(time.Second)

	d := device.Device{
		MAC:          "aa:bb:cc:dd:ee:ff",
		IP:           "192.168.1.10",
		Hostname:     "my-phone",
		Manufacturer: "Apple",
		FirstSeen:    now.Add(-24 * time.Hour),
		LastSeen:     now,
		TimesSeen:    3,
	}

	if err := s.UpsertDevice(d); err != nil {
		t.Fatalf("UpsertDevice: %v", err)
	}

	got, err := s.GetDevice("aa:bb:cc:dd:ee:ff")
	if err != nil {
		t.Fatalf("GetDevice: %v", err)
	}
	if got == nil {
		t.Fatal("GetDevice returned nil")
	}
	if got.IP != "192.168.1.10" {
		t.Errorf("IP = %q, want 192.168.1.10", got.IP)
	}
	if got.Hostname != "my-phone" {
		t.Errorf("Hostname = %q, want my-phone", got.Hostname)
	}
	if got.Manufacturer != "Apple" {
		t.Errorf("Manufacturer = %q, want Apple", got.Manufacturer)
	}
	if got.TimesSeen != 3 {
		t.Errorf("TimesSeen = %d, want 3", got.TimesSeen)
	}
}

func TestUpsertDevice_UpdateKeepsNonEmpty(t *testing.T) {
	s := newTestStore(t)
	now := time.Now().Truncate(time.Second)

	// Insert with hostname
	d := device.Device{
		MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.10",
		Hostname: "original", Manufacturer: "Mfr",
		FirstSeen: now, LastSeen: now, TimesSeen: 1,
	}
	if err := s.UpsertDevice(d); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	// Update with empty hostname — should keep "original"
	d2 := device.Device{
		MAC: "aa:bb:cc:dd:ee:ff", IP: "192.168.1.20",
		Hostname: "", Manufacturer: "",
		FirstSeen: now, LastSeen: now.Add(time.Hour), TimesSeen: 2,
	}
	if err := s.UpsertDevice(d2); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := s.GetDevice("aa:bb:cc:dd:ee:ff")
	if got.Hostname != "original" {
		t.Errorf("Hostname = %q, want %q (should keep non-empty)", got.Hostname, "original")
	}
	if got.Manufacturer != "Mfr" {
		t.Errorf("Manufacturer = %q, want %q (should keep non-empty)", got.Manufacturer, "Mfr")
	}
	if got.IP != "192.168.1.20" {
		t.Errorf("IP = %q, want 192.168.1.20 (should update)", got.IP)
	}
}

func TestGetDevice_NotFound(t *testing.T) {
	s := newTestStore(t)
	got, err := s.GetDevice("ff:ff:ff:ff:ff:ff")
	if err != nil {
		t.Fatalf("GetDevice: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for unknown MAC, got %+v", got)
	}
}

func TestGetAllDevices(t *testing.T) {
	s := newTestStore(t)
	now := time.Now().Truncate(time.Second)

	for _, mac := range []string{"aa:bb:cc:00:00:01", "aa:bb:cc:00:00:02", "aa:bb:cc:00:00:03"} {
		_ = s.UpsertDevice(device.Device{
			MAC: mac, IP: "10.0.0.1", FirstSeen: now, LastSeen: now, TimesSeen: 1,
		})
	}

	all, err := s.GetAllDevices()
	if err != nil {
		t.Fatalf("GetAllDevices: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("got %d devices, want 3", len(all))
	}
}

// --- Config CRUD ---

func TestSaveAndGetConfig(t *testing.T) {
	s := newTestStore(t)

	cfg := &config.Config{
		ScanIntervalMinutes: 5,
		RejoinThresholdDays: 14,
		NotifyNewDevice:     false,
		NotifyRejoin:        true,
		TelegramBotToken:    "tok123",
		TelegramChatID:      "-999",
		Interface:           "eth0",
	}

	if err := s.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	got, err := s.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig: %v", err)
	}

	if got.ScanIntervalMinutes != 5 {
		t.Errorf("ScanIntervalMinutes = %d, want 5", got.ScanIntervalMinutes)
	}
	if got.RejoinThresholdDays != 14 {
		t.Errorf("RejoinThresholdDays = %d, want 14", got.RejoinThresholdDays)
	}
	if got.NotifyNewDevice != false {
		t.Errorf("NotifyNewDevice = %v, want false", got.NotifyNewDevice)
	}
	if got.NotifyRejoin != true {
		t.Errorf("NotifyRejoin = %v, want true", got.NotifyRejoin)
	}
	if got.TelegramBotToken != "tok123" {
		t.Errorf("TelegramBotToken = %q, want tok123", got.TelegramBotToken)
	}
	if got.TelegramChatID != "-999" {
		t.Errorf("TelegramChatID = %q, want -999", got.TelegramChatID)
	}
	if got.Interface != "eth0" {
		t.Errorf("Interface = %q, want eth0", got.Interface)
	}
}

func TestSetConfig(t *testing.T) {
	s := newTestStore(t)

	if err := s.SetConfig(config.KeyScanInterval, "10"); err != nil {
		t.Fatalf("SetConfig: %v", err)
	}

	cfg, _ := s.GetConfig()
	if cfg.ScanIntervalMinutes != 10 {
		t.Errorf("ScanIntervalMinutes = %d, want 10", cfg.ScanIntervalMinutes)
	}
}

// --- Scan history ---

func TestRecordScan(t *testing.T) {
	s := newTestStore(t)

	if err := s.RecordScan(10, 2, 1); err != nil {
		t.Fatalf("RecordScan: %v", err)
	}
	// Verify no error on second call
	if err := s.RecordScan(8, 0, 0); err != nil {
		t.Fatalf("RecordScan (2nd): %v", err)
	}
}

