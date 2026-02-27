//go:build windows

package window

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	moduser32              = windows.NewLazySystemDLL("user32.dll")
	procEnumWindows        = moduser32.NewProc("EnumWindows")
	procGetWindowThreadPID = moduser32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible    = moduser32.NewProc("IsWindowVisible")
	procSetForegroundWnd   = moduser32.NewProc("SetForegroundWindow")
	procShowWindow         = moduser32.NewProc("ShowWindow")
	procIsIconic           = moduser32.NewProc("IsIconic")
	procGetWindowTextW     = moduser32.NewProc("GetWindowTextW")
)

const (
	swRestore = 9
)

// FindAncestorWindow walks up the process tree from the current process
// to find the best matching visible top-level window.
// When multiple windows exist for the same process (e.g., multiple VSCode/Cursor windows),
// it tries to match the window title with the current working directory.
func FindAncestorWindow() (uintptr, error) {
	pid := uint32(os.Getpid())
	chain, err := BuildParentChain(pid)
	if err != nil {
		return 0, fmt.Errorf("build parent chain: %w", err)
	}

	// Get current working directory name for matching
	cwd, _ := os.Getwd()
	cwdName := filepath.Base(cwd)

	for _, p := range chain {
		hwnd := findBestWindowForPID(p, cwdName)
		if hwnd != 0 {
			return uintptr(hwnd), nil
		}
	}
	return 0, fmt.Errorf("no visible ancestor window found")
}

// findBestWindowForPID finds the best matching visible window for the given PID.
// If cwdName is provided, it prefers windows whose title contains the cwd name.
// Falls back to any visible window if no match is found.
func findBestWindowForPID(targetPID uint32, cwdName string) windows.HWND {
	var matchedWindow windows.HWND
	var fallbackWindow windows.HWND

	cb := syscall.NewCallback(func(hwnd windows.HWND, lparam uintptr) uintptr {
		var pid uint32
		procGetWindowThreadPID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
		if pid != targetPID {
			return 1 // continue
		}
		visible, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
		if visible == 0 {
			return 1 // continue
		}

		// Get window title
		title := getWindowText(hwnd)
		if title == "" {
			return 1 // continue, skip windows without title
		}

		// Record as fallback
		if fallbackWindow == 0 {
			fallbackWindow = hwnd
		}

		// Check if title contains the cwd name (case-insensitive)
		if cwdName != "" && strings.Contains(strings.ToLower(title), strings.ToLower(cwdName)) {
			matchedWindow = hwnd
			return 0 // stop, found the best match
		}

		return 1 // continue looking for a better match
	})

	procEnumWindows.Call(cb, 0)

	if matchedWindow != 0 {
		return matchedWindow
	}
	return fallbackWindow
}

// getWindowText retrieves the title of a window.
func getWindowText(hwnd windows.HWND) string {
	buf := make([]uint16, 256)
	procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return windows.UTF16ToString(buf)
}

// Activate brings a window to the foreground. If minimized, it restores first.
func Activate(hwnd uintptr) error {
	iconic, _, _ := procIsIconic.Call(hwnd)
	if iconic != 0 {
		procShowWindow.Call(hwnd, swRestore)
	}
	r1, _, err := procSetForegroundWnd.Call(hwnd)
	if r1 == 0 {
		return fmt.Errorf("SetForegroundWindow: %w", err)
	}
	return nil
}
