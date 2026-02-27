//go:build darwin

package window

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Known macOS terminal application names and their bundle identifiers.
var knownTerminals = map[string]string{
	"Terminal":     "Terminal",
	"iTerm2":      "iTerm",
	"iTerm.app":   "iTerm",
	"Alacritty":   "Alacritty",
	"kitty":       "kitty",
	"WezTerm":     "WezTerm",
	"Hyper":       "Hyper",
	"Warp":        "Warp",
	"ghostty":     "Ghostty",
	"Ghostty":     "Ghostty",
}

// FindAncestorWindow walks up the process tree to find the terminal application.
// On macOS, the returned uintptr is unused (always 0); activation is done by app name.
func FindAncestorWindow() (uintptr, error) {
	pid := uint32(os.Getpid())
	chain, err := BuildParentChain(pid)
	if err != nil {
		return 0, fmt.Errorf("build parent chain: %w", err)
	}

	for _, p := range chain {
		name, err := getProcessName(p)
		if err != nil {
			continue
		}
		for keyword := range knownTerminals {
			if strings.Contains(name, keyword) {
				return 0, nil
			}
		}
	}
	return 0, fmt.Errorf("no ancestor terminal found")
}

// Activate brings the terminal application to the foreground via osascript.
// On macOS the hwnd parameter is unused.
func Activate(hwnd uintptr) error {
	appName := FindTerminalAppName()
	if appName == "" {
		appName = "Terminal"
	}
	script := fmt.Sprintf(`tell application "%s" to activate`, appName)
	return exec.Command("osascript", "-e", script).Run()
}

// FindTerminalAppName walks the process tree to find which terminal we're running in.
func FindTerminalAppName() string {
	pid := uint32(os.Getpid())
	chain, err := BuildParentChain(pid)
	if err != nil {
		return ""
	}

	for _, p := range chain {
		name, err := getProcessName(p)
		if err != nil {
			continue
		}
		for keyword, appName := range knownTerminals {
			if strings.Contains(name, keyword) {
				return appName
			}
		}
	}
	return ""
}

func getProcessName(pid uint32) (string, error) {
	out, err := exec.Command("ps", "-o", "comm=", "-p", fmt.Sprint(pid)).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
