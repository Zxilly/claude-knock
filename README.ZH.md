# claude-knock

[![Go Reference](https://pkg.go.dev/badge/github.com/Zxilly/claude-knock.svg)](https://pkg.go.dev/github.com/Zxilly/claude-knock)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)

[English](README.md)

跨平台桌面通知工具，用于 Claude Code hook 事件。

## 平台支持

| 平台    | 通知方式                                                                | 点击激活终端                                                   | 可选依赖                               |
| ------- | ----------------------------------------------------------------------- | -------------------------------------------------------------- | -------------------------------------- |
| Windows | Windows Toast                                                           | 支持（点击通知时）                                             | 无                                     |
| macOS   | `alerter` / `terminal-notifier` / `osascript`（按优先级选择第一个可用） | `alerter`/`terminal-notifier`：点击时激活；`osascript`：不支持 | 可选：`alerter` 或 `terminal-notifier` |
| Linux   | `notify-send`                                                           | 不支持                                                         | `notify-send`                          |

## 安装

```bash
go install github.com/Zxilly/claude-knock@latest
```

## 构建

```bash
go build
```

## 配置

添加到 `~/.claude/settings.json`：

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

## 测试

```bash
echo '{"session_id":"test","hook_event_name":"Notification","message":"Hello"}' | ./claude-knock
echo '{"session_id":"test","hook_event_name":"Stop"}' | ./claude-knock
```
