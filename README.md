# claude-knock

[![Go Reference](https://pkg.go.dev/badge/github.com/Zxilly/claude-knock.svg)](https://pkg.go.dev/github.com/Zxilly/claude-knock)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)

[简体中文](README.ZH.md)

A cross-platform desktop notification tool for Claude Code hook events.

## Platform support

| Platform | Notification method                                             | Click-to-activate terminal                                          | Optional dependencies                      |
| -------- | --------------------------------------------------------------- | ------------------------------------------------------------------- | ------------------------------------------ |
| Windows  | Windows Toast                                                   | Yes (on click)                                                      | None                                       |
| macOS    | `alerter` / `terminal-notifier` / `osascript` (first available) | `alerter`/`terminal-notifier`: on click; `osascript`: not supported | Optional: `alerter` or `terminal-notifier` |
| Linux    | `notify-send`                                                   | No                                                                  | `notify-send`                              |

## Install

```bash
go install github.com/Zxilly/claude-knock@latest
```

## Build

```bash
go build
```

## Configure

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "Notification": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/claude-knock"
          }
        ]
      }
    ],
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/claude-knock"
          }
        ]
      }
    ]
  }
}
```

## Test

```bash
echo '{"session_id":"test","hook_event_name":"Notification","message":"Hello"}' | ./claude-knock
echo '{"session_id":"test","hook_event_name":"Stop"}' | ./claude-knock
```
