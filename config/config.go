package config

import (
	"os"
	"path/filepath"
)

// Default configuration values.
const (
	DefaultScanInterval        = 1 // minutes
	DefaultRejoinThresholdDays = 7 // days
	DefaultNotifyNewDevice     = true
	DefaultNotifyRejoin        = true
	DefaultDBPath              = "" // will be set to user config dir
)

// Config holds all configuration values for the application.
type Config struct {
	// Scan settings
	ScanIntervalMinutes int    `json:"scan_interval_minutes"`
	RejoinThresholdDays int    `json:"rejoin_threshold_days"`
	Interface           string `json:"interface,omitempty"` // network interface to scan

	// Notification settings
	NotifyNewDevice  bool   `json:"notify_new_device"`
	NotifyRejoin     bool   `json:"notify_rejoin"`
	TelegramBotToken string `json:"telegram_bot_token,omitempty"`
	TelegramChatID   string `json:"telegram_chat_id,omitempty"`

	// Storage
	DBPath string `json:"db_path"`

	// Runtime flags (not stored in DB)
	Install   bool `json:"-"`
	Uninstall bool `json:"-"`
	OneShot   bool `json:"-"` // Run one scan and exit
	Verbose   bool `json:"-"`
}

// DefaultConfig returns a new Config with default values.
func DefaultConfig() *Config {
	return &Config{
		ScanIntervalMinutes: DefaultScanInterval,
		RejoinThresholdDays: DefaultRejoinThresholdDays,
		NotifyNewDevice:     DefaultNotifyNewDevice,
		NotifyRejoin:        DefaultNotifyRejoin,
		DBPath:              DefaultDBFilePath(),
	}
}

// DefaultDBFilePath returns the default path for the SQLite database.
func DefaultDBFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	appDir := filepath.Join(configDir, "lansentry")
	return filepath.Join(appDir, "lansentry.db")
}

// EnsureDBDir ensures the database directory exists.
func EnsureDBDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	return os.MkdirAll(dir, 0755)
}

// ConfigKey represents a configuration key stored in the database.
type ConfigKey string

const (
	KeyScanInterval     ConfigKey = "scan_interval_minutes"
	KeyRejoinThreshold  ConfigKey = "rejoin_threshold_days"
	KeyNotifyNewDevice  ConfigKey = "notify_new_device"
	KeyNotifyRejoin     ConfigKey = "notify_rejoin"
	KeyTelegramBotToken ConfigKey = "telegram_bot_token"
	KeyTelegramChatID   ConfigKey = "telegram_chat_id"
	KeyInterface        ConfigKey = "interface"
)

// IsTelegramConfigured returns true if Telegram credentials are set.
func (c *Config) IsTelegramConfigured() bool {
	return c.TelegramBotToken != "" && c.TelegramChatID != ""
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.ScanIntervalMinutes < 1 {
		c.ScanIntervalMinutes = DefaultScanInterval
	}
	if c.RejoinThresholdDays < 1 {
		c.RejoinThresholdDays = DefaultRejoinThresholdDays
	}
	return nil
}
