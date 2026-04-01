package config

import (
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ScanIntervalMinutes != 1 {
		t.Errorf("ScanIntervalMinutes = %d, want 1", cfg.ScanIntervalMinutes)
	}
	if cfg.RejoinThresholdDays != 7 {
		t.Errorf("RejoinThresholdDays = %d, want 7", cfg.RejoinThresholdDays)
	}
	if !cfg.NotifyNewDevice {
		t.Error("NotifyNewDevice should be true by default")
	}
	if !cfg.NotifyRejoin {
		t.Error("NotifyRejoin should be true by default")
	}
	if cfg.DBPath == "" {
		t.Error("DBPath should not be empty")
	}
}

func TestDefaultDBFilePath(t *testing.T) {
	path := DefaultDBFilePath()
	if !strings.HasSuffix(path, "lansentry/lansentry.db") {
		t.Errorf("DBPath = %q, want suffix lansentry/lansentry.db", path)
	}
}

func TestValidate_ClampsLowValues(t *testing.T) {
	cfg := &Config{ScanIntervalMinutes: 0, RejoinThresholdDays: 0}
	_ = cfg.Validate()

	if cfg.ScanIntervalMinutes != DefaultScanInterval {
		t.Errorf("ScanIntervalMinutes = %d, want %d", cfg.ScanIntervalMinutes, DefaultScanInterval)
	}
	if cfg.RejoinThresholdDays != DefaultRejoinThresholdDays {
		t.Errorf("RejoinThresholdDays = %d, want %d", cfg.RejoinThresholdDays, DefaultRejoinThresholdDays)
	}
}

func TestValidate_KeepsValidValues(t *testing.T) {
	cfg := &Config{ScanIntervalMinutes: 5, RejoinThresholdDays: 14}
	_ = cfg.Validate()

	if cfg.ScanIntervalMinutes != 5 {
		t.Errorf("ScanIntervalMinutes = %d, want 5", cfg.ScanIntervalMinutes)
	}
	if cfg.RejoinThresholdDays != 14 {
		t.Errorf("RejoinThresholdDays = %d, want 14", cfg.RejoinThresholdDays)
	}
}

func TestIsTelegramConfigured(t *testing.T) {
	tests := []struct {
		name  string
		token string
		chat  string
		want  bool
	}{
		{"both set", "tok", "123", true},
		{"token only", "tok", "", false},
		{"chat only", "", "123", false},
		{"neither", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{TelegramBotToken: tt.token, TelegramChatID: tt.chat}
			if got := cfg.IsTelegramConfigured(); got != tt.want {
				t.Errorf("IsTelegramConfigured() = %v, want %v", got, tt.want)
			}
		})
	}
}

