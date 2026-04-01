package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"lansentry/config"
	"lansentry/device"
	"lansentry/internal/service"
	"lansentry/notifier"
	"lansentry/oui"
	"lansentry/resolver"
	"lansentry/scanner"
	"lansentry/store"

	"github.com/spf13/pflag"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Parse CLI flags
	cfg := config.DefaultConfig()
	parseFlags(cfg)

	// Handle uninstall command first (does not require DB access)
	if cfg.Uninstall {
		if err := service.Uninstall(); err != nil {
			log.Fatalf("Failed to uninstall service: %v", err)
		}
		return
	}

	// Initialize store
	db, err := store.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Load config from database and merge with CLI flags
	dbConfig, err := db.GetConfig()
	if err != nil {
		log.Printf("Warning: Failed to load config from database: %v", err)
	} else {
		mergeConfig(cfg, dbConfig)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Persist final merged config so CLI overrides are remembered
	if err := db.SaveConfig(cfg); err != nil {
		log.Printf("Warning: Failed to save config to database: %v", err)
	}

	// Install service after config is merged/saved so runtime flags (e.g. Telegram creds)
	// are persisted and the service can be installed with the same DB path.
	if cfg.Install {
		if err := service.Install(cfg); err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		return
	}

	// Initialize scanner
	arpScanner, err := scanner.NewARPScanner(cfg.Interface)
	if err != nil {
		log.Fatalf("Failed to initialize scanner: %v", err)
	}

	// Initialize OUI database
	ouiDB := oui.NewOUIDatabase()

	// Initialize device name resolver
	deviceResolver := resolver.New(cfg.Verbose)

	// Initialize notifier if configured
	var telegramNotifier *notifier.TelegramNotifier
	if cfg.IsTelegramConfigured() {
		chatID, err := strconv.ParseInt(cfg.TelegramChatID, 10, 64)
		if err != nil {
			log.Printf("Warning: Invalid Telegram chat ID: %v", err)
		} else {
			telegramNotifier, err = notifier.NewTelegram(cfg.TelegramBotToken, chatID)
			if err != nil {
				log.Printf("Warning: Failed to initialize Telegram notifier: %v", err)
			}
		}
	}

	// Initialize detector
	detector := device.NewDetector(cfg.RejoinThresholdDays)

	// Print startup info
	printStartupInfo(cfg, arpScanner)

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping...")
		cancel()
	}()

	// Run main loop
	if cfg.OneShot {
		runScan(db, arpScanner, detector, ouiDB, deviceResolver, telegramNotifier, cfg)
	} else {
		runDaemon(ctx, db, arpScanner, detector, ouiDB, deviceResolver, telegramNotifier, cfg)
	}
}

func parseFlags(cfg *config.Config) {
	pflag.IntVar(&cfg.ScanIntervalMinutes, "scan-interval", config.DefaultScanInterval, "Scan interval in minutes")
	pflag.IntVar(&cfg.RejoinThresholdDays, "rejoin-days", config.DefaultRejoinThresholdDays, "Days before a device is considered rejoining")
	pflag.BoolVar(&cfg.NotifyNewDevice, "notify-new", config.DefaultNotifyNewDevice, "Send notifications for new devices")
	pflag.BoolVar(&cfg.NotifyRejoin, "notify-rejoin", config.DefaultNotifyRejoin, "Send notifications for rejoining devices")
	pflag.StringVar(&cfg.TelegramBotToken, "telegram-token", "", "Telegram bot token")
	pflag.StringVar(&cfg.TelegramChatID, "telegram-chat", "", "Telegram chat ID")
	pflag.StringVar(&cfg.Interface, "interface", "", "Network interface to scan (auto-detect if empty)")
	pflag.StringVar(&cfg.DBPath, "db-path", config.DefaultDBFilePath(), "Path to SQLite database")
	pflag.BoolVar(&cfg.Install, "install", false, "Install as system service")
	pflag.BoolVar(&cfg.Uninstall, "uninstall", false, "Uninstall system service")
	pflag.BoolVar(&cfg.OneShot, "one-shot", false, "Run a single scan and exit")
	pflag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")

	showVersion := pflag.Bool("version", false, "Show version information")

	pflag.Parse()

	if *showVersion {
		fmt.Printf("LANSentry %s\n", version)
		fmt.Printf("  Commit: %s\n", commit)
		fmt.Printf("  Built:  %s\n", date)
		os.Exit(0)
	}
}

// mergeConfig merges CLI flags with database config.
// CLI flags take precedence if explicitly set.
func mergeConfig(cli *config.Config, db *config.Config) {
	// Only use DB values if CLI flags weren't explicitly set
	if !pflag.CommandLine.Changed("scan-interval") {
		cli.ScanIntervalMinutes = db.ScanIntervalMinutes
	}
	if !pflag.CommandLine.Changed("rejoin-days") {
		cli.RejoinThresholdDays = db.RejoinThresholdDays
	}
	if !pflag.CommandLine.Changed("notify-new") {
		cli.NotifyNewDevice = db.NotifyNewDevice
	}
	if !pflag.CommandLine.Changed("notify-rejoin") {
		cli.NotifyRejoin = db.NotifyRejoin
	}
	if !pflag.CommandLine.Changed("telegram-token") && db.TelegramBotToken != "" {
		cli.TelegramBotToken = db.TelegramBotToken
	}
	if !pflag.CommandLine.Changed("telegram-chat") && db.TelegramChatID != "" {
		cli.TelegramChatID = db.TelegramChatID
	}
	if !pflag.CommandLine.Changed("interface") && db.Interface != "" {
		cli.Interface = db.Interface
	}
}

func printStartupInfo(cfg *config.Config, s *scanner.ARPScanner) {
	log.Println("╔═══════════════════════════════════════════╗")
	log.Println("║         LANSentry Network Monitor         ║")
	log.Println("╚═══════════════════════════════════════════╝")
	log.Printf("Interface:        %s", s.Interface())
	log.Printf("Scan interval:    %d minute(s)", cfg.ScanIntervalMinutes)
	log.Printf("Rejoin threshold: %d day(s)", cfg.RejoinThresholdDays)
	log.Printf("Notify new:       %v", cfg.NotifyNewDevice)
	log.Printf("Notify rejoin:    %v", cfg.NotifyRejoin)
	log.Printf("Telegram:         %v", cfg.IsTelegramConfigured())
	log.Printf("Database:         %s", cfg.DBPath)
	log.Println("───────────────────────────────────────────")
}

func runDaemon(ctx context.Context, db *store.Store, s *scanner.ARPScanner, detector *device.Detector, ouiDB *oui.OUIDatabase, res *resolver.Resolver, notif *notifier.TelegramNotifier, cfg *config.Config) {
	ticker := time.NewTicker(time.Duration(cfg.ScanIntervalMinutes) * time.Minute)
	defer ticker.Stop()

	// Run first scan immediately
	runScan(db, s, detector, ouiDB, res, notif, cfg)

	for {
		select {
		case <-ctx.Done():
			log.Println("Daemon stopped")
			return
		case <-ticker.C:
			runScan(db, s, detector, ouiDB, res, notif, cfg)
		}
	}
}

func runScan(db *store.Store, s *scanner.ARPScanner, detector *device.Detector, ouiDB *oui.OUIDatabase, res *resolver.Resolver, notif *notifier.TelegramNotifier, cfg *config.Config) {
	log.Printf("Starting network scan...")
	startTime := time.Now()

	// Perform scan
	scannedDevices, err := s.Scan()
	if err != nil {
		log.Printf("Scan failed: %v", err)
		return
	}

	log.Printf("Found %d device(s) in %v", len(scannedDevices), time.Since(startTime).Round(time.Millisecond))

	// Enrich with manufacturer info
	for i := range scannedDevices {
		if scannedDevices[i].Manufacturer == "" {
			scannedDevices[i].Manufacturer = ouiDB.LookupManufacturer(scannedDevices[i].MAC)
		}
	}

	// Get known devices from database
	knownDevices, err := db.GetAllDevices()
	if err != nil {
		log.Printf("Failed to get known devices: %v", err)
		return
	}

	// Pre-populate hostname from known data so the resolver
	// only needs to look up genuinely unknown devices.
	for i := range scannedDevices {
		mac := device.NormalizeMac(scannedDevices[i].MAC)
		if known, ok := knownDevices[mac]; ok {
			if scannedDevices[i].Hostname == "" && known.Hostname != "" {
				scannedDevices[i].Hostname = known.Hostname
			}
		}
	}

	// Resolve remaining hostnames
	if cfg.Verbose {
		log.Printf("Resolving device names...")
	}
	res.EnrichDevices(scannedDevices)

	// Detect events
	events := detector.DetectEvents(scannedDevices, knownDevices)
	updates := detector.UpdateDevices(scannedDevices, knownDevices)

	// Count events by type
	var newCount, rejoinCount int
	for _, event := range events {
		switch event.Type {
		case device.EventNew:
			newCount++
		case device.EventRejoin:
			rejoinCount++
		}
	}

	// Log events
	for _, event := range events {
		name := deviceDisplayName(event.Device)
		switch event.Type {
		case device.EventNew:
			log.Printf("🆕 NEW DEVICE: %s (%s) %s %s",
				event.Device.IP, event.Device.MAC, event.Device.Manufacturer, name)
		case device.EventRejoin:
			log.Printf("🔄 REJOINED: %s (%s) %s - offline for %s",
				event.Device.IP, event.Device.MAC, name, formatDuration(event.OfflineFor))
		}
	}

	// Send notifications
	if notif != nil {
		for _, event := range events {
			shouldNotify := false
			switch event.Type {
			case device.EventNew:
				shouldNotify = cfg.NotifyNewDevice
			case device.EventRejoin:
				shouldNotify = cfg.NotifyRejoin
			}

			if shouldNotify {
				if err := notif.NotifyEvent(event); err != nil {
					log.Printf("Failed to send notification: %v", err)
				}
			}
		}
	}

	// Update database
	for _, dev := range updates {
		if err := db.UpsertDevice(dev); err != nil {
			log.Printf("Failed to update device %s: %v", dev.MAC, err)
		}
	}

	// Record scan history
	if err := db.RecordScan(len(scannedDevices), newCount, rejoinCount); err != nil {
		log.Printf("Failed to record scan history: %v", err)
	}

	if cfg.Verbose {
		log.Printf("Scan complete: %d total, %d new, %d rejoined", len(scannedDevices), newCount, rejoinCount)
	}
}

// deviceDisplayName returns a human-readable label from the hostname.
func deviceDisplayName(d device.Device) string {
	if d.Hostname != "" {
		return "- " + d.Hostname
	}
	return ""
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
