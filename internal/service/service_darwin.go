package service

import (
	"fmt"
	"lansentry/config"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const serviceName = "com.lansentry.agent"
const launchdPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.ServiceName}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.ExecPath}}</string>
		<string>--db-path</string>
		<string>{{.DBPath}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.LogDir}}/lansentry.log</string>
    <key>StandardErrorPath</key>
    <string>{{.LogDir}}/lansentry.error.log</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>
`

type launchdConfig struct {
	ServiceName string
	ExecPath    string
	DBPath      string
	LogDir      string
}

// Install installs the service to start on boot (user-level launchd agent).
func Install(cfg *config.Config) error {
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

	// Get home directory for user-level agent
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create LaunchAgents directory if needed
	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	// Create log directory
	logDir := filepath.Join(homeDir, "Library", "Logs", "LANSentry")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate plist
	plistPath := filepath.Join(launchAgentsDir, serviceName+".plist")

	config := launchdConfig{
		ServiceName: serviceName,
		ExecPath:    execPath,
		DBPath:      cfg.DBPath,
		LogDir:      logDir,
	}

	tmpl, err := template.New("plist").Parse(launchdPlistTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to write plist: %w", err)
	}

	// Load the service
	cmd := exec.Command("launchctl", "load", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load service: %w, output: %s", err, string(output))
	}

	fmt.Printf("✅ LANSentry installed successfully!\n")
	fmt.Printf("   Plist: %s\n", plistPath)
	fmt.Printf("   Logs:  %s\n", logDir)
	fmt.Printf("\nThe service will start automatically at login.\n")
	fmt.Printf("To start it now: launchctl start %s\n", serviceName)

	return nil
}

// Uninstall removes the service.
func Uninstall() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", serviceName+".plist")

	// Unload the service (ignore errors if not loaded)
	cmd := exec.Command("launchctl", "unload", plistPath)
	_ = cmd.Run()

	// Remove plist file
	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	fmt.Printf("✅ LANSentry uninstalled successfully!\n")
	fmt.Printf("   Removed: %s\n", plistPath)

	return nil
}

// Status returns the current status of the service.
func Status() (string, error) {
	cmd := exec.Command("launchctl", "list", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "not installed", nil
	}
	return fmt.Sprintf("installed and running\n%s", string(output)), nil
}
