package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"lansentry/config"
	"lansentry/device"

	_ "modernc.org/sqlite"
)

// Store handles all database operations.
type Store struct {
	db *sql.DB
}

// New creates a new Store and initializes the database.
func New(dbPath string) (*Store, error) {
	if err := config.EnsureDBDir(dbPath); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// GetDevice retrieves a device by MAC address.
func (s *Store) GetDevice(mac string) (*device.Device, error) {
	mac = device.NormalizeMac(mac)
	row := s.db.QueryRow(`
		SELECT mac, ip, hostname, manufacturer, first_seen, last_seen, times_seen
		FROM devices WHERE mac = ?
	`, mac)

	var d device.Device
	var firstSeen, lastSeen int64
	err := row.Scan(&d.MAC, &d.IP, &d.Hostname, &d.Manufacturer, &firstSeen, &lastSeen, &d.TimesSeen)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	d.FirstSeen = time.Unix(firstSeen, 0)
	d.LastSeen = time.Unix(lastSeen, 0)
	return &d, nil
}

// GetAllDevices retrieves all known devices.
func (s *Store) GetAllDevices() (map[string]device.Device, error) {
	rows, err := s.db.Query(`
		SELECT mac, ip, hostname, manufacturer, first_seen, last_seen, times_seen
		FROM devices
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := make(map[string]device.Device)
	for rows.Next() {
		var d device.Device
		var firstSeen, lastSeen int64
		if err := rows.Scan(&d.MAC, &d.IP, &d.Hostname, &d.Manufacturer, &firstSeen, &lastSeen, &d.TimesSeen); err != nil {
			return nil, err
		}
		d.FirstSeen = time.Unix(firstSeen, 0)
		d.LastSeen = time.Unix(lastSeen, 0)
		devices[d.MAC] = d
	}

	return devices, rows.Err()
}

// UpsertDevice inserts or updates a device.
func (s *Store) UpsertDevice(d device.Device) error {
	d.MAC = device.NormalizeMac(d.MAC)
	_, err := s.db.Exec(`
		INSERT INTO devices (mac, ip, hostname, manufacturer, first_seen, last_seen, times_seen)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(mac) DO UPDATE SET
			ip = excluded.ip,
			hostname = CASE WHEN excluded.hostname != '' THEN excluded.hostname ELSE hostname END,
			manufacturer = CASE WHEN excluded.manufacturer != '' THEN excluded.manufacturer ELSE manufacturer END,
			last_seen = excluded.last_seen,
			times_seen = excluded.times_seen
	`, d.MAC, d.IP, d.Hostname, d.Manufacturer, d.FirstSeen.Unix(), d.LastSeen.Unix(), d.TimesSeen)
	return err
}

// Config operations

// GetConfig retrieves all configuration values and merges with defaults.
func (s *Store) GetConfig() (*config.Config, error) {
	cfg := config.DefaultConfig()

	rows, err := s.db.Query(`SELECT key, value FROM config`)
	if err != nil {
		return cfg, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}

		switch config.ConfigKey(key) {
		case config.KeyScanInterval:
			if v, err := strconv.Atoi(value); err == nil {
				cfg.ScanIntervalMinutes = v
			}
		case config.KeyRejoinThreshold:
			if v, err := strconv.Atoi(value); err == nil {
				cfg.RejoinThresholdDays = v
			}
		case config.KeyNotifyNewDevice:
			cfg.NotifyNewDevice = value == "true" || value == "1"
		case config.KeyNotifyRejoin:
			cfg.NotifyRejoin = value == "true" || value == "1"
		case config.KeyTelegramBotToken:
			cfg.TelegramBotToken = value
		case config.KeyTelegramChatID:
			cfg.TelegramChatID = value
		case config.KeyInterface:
			cfg.Interface = value
		}
	}

	return cfg, rows.Err()
}

// SetConfig saves a configuration value to the database.
func (s *Store) SetConfig(key config.ConfigKey, value string) error {
	_, err := s.db.Exec(`
		INSERT INTO config (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP
	`, string(key), value)
	return err
}

// SaveConfig persists all configuration values to the database.
func (s *Store) SaveConfig(cfg *config.Config) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updates := map[config.ConfigKey]string{
		config.KeyScanInterval:     strconv.Itoa(cfg.ScanIntervalMinutes),
		config.KeyRejoinThreshold:  strconv.Itoa(cfg.RejoinThresholdDays),
		config.KeyNotifyNewDevice:  strconv.FormatBool(cfg.NotifyNewDevice),
		config.KeyNotifyRejoin:     strconv.FormatBool(cfg.NotifyRejoin),
		config.KeyTelegramBotToken: cfg.TelegramBotToken,
		config.KeyTelegramChatID:   cfg.TelegramChatID,
		config.KeyInterface:        cfg.Interface,
	}

	for key, value := range updates {
		if _, err := tx.Exec(`
			INSERT INTO config (key, value) VALUES (?, ?)
			ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP
		`, string(key), value); err != nil {
			return err
		}
	}

	return tx.Commit()
}
