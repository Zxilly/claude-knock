//go:build linux

package notify

import (
	"fmt"
	"os/exec"
	"time"
)

// Send sends a desktop notification using notify-send.
func Send(title, message string) error {
	notifySend, err := exec.LookPath("notify-send")
	if err != nil {
		return fmt.Errorf("notify-send not found: %w", err)
	}
	return exec.Command(notifySend, title, message).Run()
}

// HandleCOMActivation is a no-op on Linux (Windows COM activation only).
func HandleCOMActivation(timeout time.Duration) {}
