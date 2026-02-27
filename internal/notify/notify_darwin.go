//go:build darwin

package notify

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/Zxilly/claude-knock/internal/window"
)

// Send sends a desktop notification on macOS.
// It tries alerter, terminal-notifier, and osascript in order.
func Send(title, message string) error {
	if p, _ := exec.LookPath("alerter"); p != "" {
		return sendAlerter(p, title, message)
	}
	if p, _ := exec.LookPath("terminal-notifier"); p != "" {
		return sendTerminalNotifier(p, title, message)
	}
	return sendOsascript(title, message)
}

func sendAlerter(bin, title, message string) error {
	bundleID := findTerminalBundleID()
	return exec.Command(bin,
		"--title", title,
		"--message", message,
		"--sender", bundleID,
	).Run()
}

func sendTerminalNotifier(bin, title, message string) error {
	bundleID := findTerminalBundleID()
	return exec.Command(bin,
		"-title", title,
		"-message", message,
		"-sender", bundleID,
		"-activate", bundleID,
	).Run()
}

func sendOsascript(title, message string) error {
	script := fmt.Sprintf(`display notification %q with title %q`, message, title)
	return exec.Command("osascript", "-e", script).Run()
}

var terminalBundleIDs = map[string]string{
	"Terminal":   "com.apple.Terminal",
	"iTerm":     "com.googlecode.iterm2",
	"Alacritty": "org.alacritty",
	"kitty":     "net.kovidgoyal.kitty",
	"WezTerm":   "com.github.wez.wezterm",
	"Hyper":     "co.zeit.hyper",
	"Warp":      "dev.warp.Warp-Stable",
	"Ghostty":   "com.mitchellh.ghostty",
}

func findTerminalBundleID() string {
	appName := window.FindTerminalAppName()
	if id, ok := terminalBundleIDs[appName]; ok {
		return id
	}
	return "com.apple.Terminal"
}

// HandleCOMActivation is a no-op on macOS (Windows COM activation only).
func HandleCOMActivation(timeout time.Duration) {}
