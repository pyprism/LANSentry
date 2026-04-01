package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const serviceName = "lansentry"
const systemdUnitTemplate = `[Unit]
Description=LANSentry Network Device Monitor
After=network.target

[Service]
Type=simple
ExecStart={{.ExecPath}}
Restart=on-failure
RestartSec=10
StandardOutput=append:{{.LogDir}}/lansentry.log
StandardError=append:{{.LogDir}}/lansentry.error.log
Environment="PATH=/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"

[Install]
WantedBy=default.target
`

type systemdConfig struct {
	ExecPath string
	LogDir   string
}

// Install installs the service to start on boot (user-level systemd).
func Install() error {
	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Make path absolute
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Get home directory for user-level service
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create user systemd directory if needed
	systemdDir := filepath.Join(homeDir, ".config", "systemd", "user")
	if err := os.MkdirAll(systemdDir, 0755); err != nil {
		return fmt.Errorf("failed to create systemd directory: %w", err)
	}

	// Create log directory
	logDir := filepath.Join(homeDir, ".local", "share", "lansentry", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate unit file
	unitPath := filepath.Join(systemdDir, serviceName+".service")

	config := systemdConfig{
		ExecPath: execPath,
		LogDir:   logDir,
	}

	tmpl, err := template.New("unit").Parse(systemdUnitTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(unitPath)
	if err != nil {
		return fmt.Errorf("failed to create unit file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to write unit file: %w", err)
	}

	// Reload systemd and enable service
	if err := runSystemctl("--user", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	if err := runSystemctl("--user", "enable", serviceName+".service"); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	// Enable lingering so user services can run without active session
	if err := enableLingering(); err != nil {
		fmt.Printf("⚠️  Warning: Could not enable lingering. Service may not start at boot without active session.\n")
		fmt.Printf("   Run: loginctl enable-linger %s\n", os.Getenv("USER"))
	}

	fmt.Printf("✅ LANSentry installed successfully!\n")
	fmt.Printf("   Unit: %s\n", unitPath)
	fmt.Printf("   Logs: %s\n", logDir)
	fmt.Printf("\nTo start the service now:\n")
	fmt.Printf("   systemctl --user start %s\n", serviceName)
	fmt.Printf("\nTo view logs:\n")
	fmt.Printf("   journalctl --user -u %s -f\n", serviceName)

	return nil
}

// Uninstall removes the service.
func Uninstall() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	unitPath := filepath.Join(homeDir, ".config", "systemd", "user", serviceName+".service")

	// Stop and disable the service (ignore errors if not running/enabled)
	_ = runSystemctl("--user", "stop", serviceName+".service")
	_ = runSystemctl("--user", "disable", serviceName+".service")

	// Remove unit file
	if err := os.Remove(unitPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove unit file: %w", err)
	}

	// Reload systemd
	_ = runSystemctl("--user", "daemon-reload")

	fmt.Printf("✅ LANSentry uninstalled successfully!\n")
	fmt.Printf("   Removed: %s\n", unitPath)

	return nil
}

// Status returns the current status of the service.
func Status() (string, error) {
	cmd := exec.Command("systemctl", "--user", "status", serviceName+".service")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("not running or not installed\n%s", string(output)), nil
	}
	return string(output), nil
}

func runSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}
	return nil
}

func enableLingering() error {
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("USER environment variable not set")
	}
	cmd := exec.Command("loginctl", "enable-linger", user)
	return cmd.Run()
}
