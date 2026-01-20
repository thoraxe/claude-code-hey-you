# claude-ntfy

Minimal cross-platform bridge from Claude Code hooks to [ntfy](https://ntfy.sh) notifications.

## Installation

Download the binary for your platform from the [releases page](https://github.com/user/claude-ntfy/releases), or build from source:

```bash
go build -ldflags "-s -w" -o claude-ntfy
```

## Usage

```bash
# Basic - topic as argument
claude-ntfy --topic my-alerts

# With custom server
claude-ntfy --topic my-alerts --server https://ntfy.example.com

# With priority
claude-ntfy --topic my-alerts --priority high

# Via environment variables
NTFY_TOPIC=my-alerts claude-ntfy
```

## Configuration

Configuration priority (highest to lowest):
1. CLI flags (`--topic`, `--server`, `--priority`)
2. Environment variables (`NTFY_TOPIC`, `NTFY_SERVER`, `NTFY_PRIORITY`)
3. Config file

Config file location:
- Windows: `%APPDATA%\claude-ntfy\config.toml`
- Unix: `~/.config/claude-ntfy/config.toml`

```toml
topic = "my-alerts"
server = "https://ntfy.sh"  # default
priority = "default"        # min, low, default, high, urgent
```

## Claude Code Integration

Add to `~/.claude/settings.json` (or `%USERPROFILE%\.claude\settings.json` on Windows):

```json
{
  "hooks": {
    "PermissionRequest": [
      {
        "matcher": "",
        "hooks": [{"type": "command", "command": "claude-ntfy --topic my-alerts"}]
      }
    ],
    "Notification": [
      {
        "matcher": "",
        "hooks": [{"type": "command", "command": "claude-ntfy --topic my-alerts"}]
      }
    ],
    "Stop": [
      {
        "matcher": "",
        "hooks": [{"type": "command", "command": "claude-ntfy --topic my-alerts"}]
      }
    ]
  }
}
```

## How It Works

Claude Code passes JSON to hook commands via stdin. This tool parses the JSON and sends contextual notifications to ntfy.

| Event | Title | Body |
|-------|-------|------|
| PermissionRequest (Bash) | "Permission: Bash" | The command |
| PermissionRequest (Write/Edit) | "Permission: Write" | The file path |
| Notification | "Claude Code" | The message |
| Stop | "Task Complete" | Working directory |

## Building

```bash
# Current platform
make build

# All platforms
make build-all
```

## License

MIT
