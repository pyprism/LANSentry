# LANSentry [![codecov](https://codecov.io/gh/pyprism/LANSentry/graph/badge.svg?token=3ZMjhNXiAK)](https://codecov.io/gh/pyprism/LANSentry)

A lightweight, cross platform network device monitor that detects when devices join or rejoin your local network and sends Telegram notifications.

## Features

- **Periodic Network Scanning** - Scans your network every minute (configurable) for connected devices
- **New Device Detection** - Alerts when a device never seen before joins the network
- **Rejoin Detection** - Alerts when a known device returns after being offline
- **Telegram Notifications** - Real time notifications to your Telegram
- **SQLite Persistence** - Stores device history and configuration
- **Manufacturer Detection** - Identifies device manufacturers via MAC OUI lookup
- **Flexible Configuration** - CLI flags with database backed defaults
- **Cross Platform** - Works on Linux and macOS
- **Service Installation** - Auto starts on boot (user level, no root required)

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/pyprism/LANSentry.git
cd lansentry

# Build
go build -o lansentry ./cmd/netwatcher

# Optional: Install to PATH
sudo mv lansentry /usr/local/bin/
```

### Dependencies

For best results, install `arp-scan`:

**macOS:**
```bash
brew install arp-scan
```

**Linux (Debian/Ubuntu):**
```bash
sudo apt install arp-scan
```

**Linux (Fedora/RHEL):**
```bash
sudo dnf install arp-scan
```

> Note: LANSentry works without `arp-scan` using a fallback ping sweep + ARP table method, but `arp-scan` provides faster and more reliable results.

## Usage

### Basic Usage

```bash
# Run a single scan
./lansentry --one-shot

# Run as daemon (foreground)
./lansentry

# Run with verbose output
./lansentry --verbose
```

### Configuration

```bash
# Set scan interval (in minutes)
./lansentry --scan-interval 5

# Set rejoin threshold (in days)
./lansentry --rejoin-days 2

# Specify network interface
./lansentry --interface en0

# Configure Telegram notifications
./lansentry --telegram-token "BOT_TOKEN" --telegram-chat "CHAT_ID"

# Disable new device notifications
./lansentry --notify-new=false

# Disable rejoin notifications
./lansentry --notify-rejoin=false
```

### Service Installation

Install as a user level service (auto-starts on login/boot):

```bash
# Install the service
./lansentry --install

# Uninstall the service
./lansentry --uninstall
```

**macOS:** Creates a launchd agent in `~/Library/LaunchAgents/`

**Linux:** Creates a systemd user service in `~/.config/systemd/user/`

### All Options

| Flag | Description | Default |
|------|-------------|---------|
| `--scan-interval` | Scan interval in minutes | 1 |
| `--rejoin-days` | Days before device is considered rejoining | 1 |
| `--notify-new` | Enable new device notifications | true |
| `--notify-rejoin` | Enable rejoin notifications | true |
| `--telegram-token` | Telegram bot token | (none) |
| `--telegram-chat` | Telegram chat ID (private: `123...`, group/supergroup: `-100...`) | (none) |
| `--interface` | Network interface to scan | (auto) |
| `--db-path` | Path to SQLite database | ~/.config/lansentry/lansentry.db |
| `--install` | Install as system service | - |
| `--uninstall` | Uninstall system service | - |
| `--one-shot` | Run single scan and exit | false |
| `--verbose` | Enable verbose logging | false |
| `--version` | Show version info | - |

## Telegram Setup

1. Create a bot with [@BotFather](https://t.me/botfather) and copy the bot token.
2. Open your bot chat in Telegram and send a message (for groups, add the bot to the group and send a message there).
3. Query updates and inspect `message.chat.id`:

```bash
curl "https://api.telegram.org/bot<TOKEN>/getUpdates"
```

4. Use the returned `chat.id` as `--telegram-chat`:
   - Private chat IDs are usually positive (example: `123456789`)
   - Group/supergroup IDs are usually negative (example: `-1001234567890`)

5. Configure LANSentry:

```bash
./lansentry --telegram-token "123456:ABC..." --telegram-chat "123456789"
```

If `getUpdates` returns an empty `result`:

- Send a fresh message to the bot or group and run the command again.
- For groups, confirm the bot is added and allowed to read messages.
- If you already consumed updates in another script, clear offsets there or retry after a new message.

## Data Storage

LANSentry stores all data in SQLite:

- **Default location:** `~/.config/lansentry/lansentry.db`
- **Device records:** MAC, IP, hostname, manufacturer, timestamps
- **Configuration:** Persisted settings
- **Scan history:** Historical scan results

## How It Works

1. **Startup:**
   - Load configuration from database
   - Apply CLI flag overrides
   - Initialize network scanner

2. **Each Scan Interval:**
   - Send ARP requests to all IPs in local subnet
   - Collect responses to build device list
   - Compare with stored device database
   - Detect new devices and rejoins
   - Send notifications (if configured)
   - Update database

3. **Device States:**
   - **New:** MAC address never seen before
   - **Rejoined:** Known device offline longer than threshold

## Permissions

LANSentry uses system ARP commands that may require elevated privileges for raw socket access:

**Option 1:** Run with sudo (not recommended for daemon)
```bash
sudo ./lansentry
```

**Option 2:** Set capabilities (Linux only)
```bash
sudo setcap cap_net_raw+ep ./lansentry
```

**Option 3:** Use arp-scan (recommended)
Install `arp-scan` and run it with appropriate permissions.

## Logs

**macOS:** `~/Library/Logs/LANSentry/`

**Linux:** `~/.local/share/lansentry/logs/`

View systemd logs (Linux):
```bash
journalctl --user -u lansentry -f
```

## License

MIT

