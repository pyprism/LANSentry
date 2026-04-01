package notifier

import (
	"fmt"
	"strings"
	"time"

	"lansentry/device"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramNotifier handles sending notifications via Telegram.
type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

// NewTelegram creates a new Telegram notifier.
func NewTelegram(botToken string, chatID int64) (*TelegramNotifier, error) {
	if botToken == "" {
		return nil, fmt.Errorf("telegram bot token is required")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

// NotifyNewDevice sends a notification about a new device.
func (t *TelegramNotifier) NotifyNewDevice(d device.Device) error {
	msg := t.formatNewDeviceMessage(d)
	return t.send(msg)
}

// NotifyRejoin sends a notification about a rejoining device.
func (t *TelegramNotifier) NotifyRejoin(d device.Device, offlineFor time.Duration) error {
	msg := t.formatRejoinMessage(d, offlineFor)
	return t.send(msg)
}

// NotifyEvent sends a notification based on the event type.
func (t *TelegramNotifier) NotifyEvent(event device.DeviceEvent) error {
	switch event.Type {
	case device.EventNew:
		return t.NotifyNewDevice(event.Device)
	case device.EventRejoin:
		return t.NotifyRejoin(event.Device, event.OfflineFor)
	default:
		return nil
	}
}

func (t *TelegramNotifier) formatNewDeviceMessage(d device.Device) string {
	var sb strings.Builder

	sb.WriteString("🆕 *New Device Joined Network*\n\n")
	sb.WriteString(fmt.Sprintf("📍 *IP:* `%s`\n", d.IP))
	sb.WriteString(fmt.Sprintf("🔗 *MAC:* `%s`\n", d.MAC))

	if d.Manufacturer != "" {
		sb.WriteString(fmt.Sprintf("🏭 *Manufacturer:* %s\n", escapeMarkdown(d.Manufacturer)))
	}
	if d.Hostname != "" {
		sb.WriteString(fmt.Sprintf("💻 *Hostname:* %s\n", escapeMarkdown(d.Hostname)))
	}

	sb.WriteString(fmt.Sprintf("\n⏰ *First Seen:* %s", d.FirstSeen.Format("2006-01-02 3:04 PM")))

	return sb.String()
}

func (t *TelegramNotifier) formatRejoinMessage(d device.Device, offlineFor time.Duration) string {
	var sb strings.Builder

	sb.WriteString("🔄 *Device Rejoined Network*\n\n")
	sb.WriteString(fmt.Sprintf("📍 *IP:* `%s`\n", d.IP))
	sb.WriteString(fmt.Sprintf("🔗 *MAC:* `%s`\n", d.MAC))

	if d.Manufacturer != "" {
		sb.WriteString(fmt.Sprintf("🏭 *Manufacturer:* %s\n", escapeMarkdown(d.Manufacturer)))
	}
	if d.Hostname != "" {
		sb.WriteString(fmt.Sprintf("💻 *Hostname:* %s\n", escapeMarkdown(d.Hostname)))
	}

	sb.WriteString(fmt.Sprintf("\n⏱ *Offline For:* %s\n", formatDuration(offlineFor)))
	sb.WriteString(fmt.Sprintf("📊 *Times Seen:* %d\n", d.TimesSeen))
	sb.WriteString(fmt.Sprintf("📅 *First Seen:* %s", d.FirstSeen.Format("2006-01-02 3:04 PM")))

	return sb.String()
}

func (t *TelegramNotifier) send(message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	return nil
}

// escapeMarkdown escapes special characters for Telegram Markdown v1.
// Only _ * ` [ are special in v1 mode.
func escapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
		"[", "\\[",
	)
	return replacer.Replace(s)
}

// formatDuration formats a duration in a human-readable way.
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	return strings.Join(parts, " ")
}
