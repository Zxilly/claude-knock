//go:build windows

package window

import (
	"fmt"
	"os"
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
)

const (
	swRestore = 9
)

// FindAncestorWindow walks up the process tree from the current process
// to find the first visible top-level window.
func FindAncestorWindow() (uintptr, error) {
	pid := uint32(os.Getpid())
	chain, err := BuildParentChain(pid)
	if err != nil {
		return 0, fmt.Errorf("build parent chain: %w", err)
	}

	for _, p := range chain {
		hwnd := findVisibleWindowForPID(p)
		if hwnd != 0 {
			return uintptr(hwnd), nil
		}
	}
	return 0, fmt.Errorf("no visible ancestor window found")
}

// findVisibleWindowForPID finds a visible top-level window owned by the given PID.
func findVisibleWindowForPID(targetPID uint32) windows.HWND {
	var found windows.HWND

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
		found = hwnd
		return 0 // stop
	})

	procEnumWindows.Call(cb, 0)
	return found
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
