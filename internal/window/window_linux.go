//go:build linux

package window

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// FindAncestorWindow walks up the process tree and uses xdotool to find
// a window belonging to an ancestor process.
func FindAncestorWindow() (uintptr, error) {
	pid := uint32(os.Getpid())
	chain, err := BuildParentChain(pid)
	if err != nil {
		return 0, fmt.Errorf("build parent chain: %w", err)
	}

	xdotool, _ := exec.LookPath("xdotool")
	if xdotool == "" {
		return 0, fmt.Errorf("xdotool not found")
	}

	for _, p := range chain {
		out, err := exec.Command(xdotool, "search", "--pid", fmt.Sprint(p)).Output()
		if err != nil {
			continue
		}
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				return 0, nil // found a window — we don't need the ID, just confirm existence
			}
		}
	}
	return 0, fmt.Errorf("no ancestor window found")
}

// Activate brings a window to the foreground using xdotool or wmctrl.
// On Linux the hwnd parameter is unused; we activate by searching for the
// terminal process window.
func Activate(hwnd uintptr) error {
	pid := uint32(os.Getpid())
	chain, _ := BuildParentChain(pid)

	// Try xdotool first.
	if xdotool, _ := exec.LookPath("xdotool"); xdotool != "" {
		for _, p := range chain {
			out, err := exec.Command(xdotool, "search", "--pid", fmt.Sprint(p)).Output()
			if err != nil {
				continue
			}
			for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					if err := exec.Command(xdotool, "windowactivate", line).Run(); err == nil {
						return nil
					}
				}
			}
		}
	}

	// Fallback to wmctrl.
	if wmctrl, _ := exec.LookPath("wmctrl"); wmctrl != "" {
		for _, p := range chain {
			out, err := exec.Command("xdotool", "search", "--pid", fmt.Sprint(p)).Output()
			if err != nil {
				continue
			}
			for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					if err := exec.Command(wmctrl, "-i", "-a", line).Run(); err == nil {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("no window activation tool available")
}
