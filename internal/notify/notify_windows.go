//go:build windows

package notify

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"

	toast "git.sr.ht/~jackmordaunt/go-toast/v2"
	"git.sr.ht/~jackmordaunt/go-toast/v2/wintoast"
	"github.com/Zxilly/claude-knock/internal/window"
	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	modcombase              = windows.NewLazySystemDLL("combase.dll")
	procRegisterClassObject = modcombase.NewProc("CoRegisterClassObject")
)

// Send sends a toast notification with a "Go to Terminal" action button.
// The hwnd is encoded in the action arguments so the COM relaunch can use it.
// This function returns immediately after pushing the notification.
func Send(title, message string) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	hwnd, _ := window.FindAncestorWindow()
	args := fmt.Sprintf("hwnd:%d", hwnd)

	if err := toast.SetAppData(toast.AppData{
		AppID:         "Claude Code",
		ActivationExe: exe,
	}); err != nil {
		return fmt.Errorf("set app data: %w", err)
	}

	noti := toast.Notification{
		AppID:               "Claude Code",
		Title:               title,
		Body:                message,
		ActivationExe:       exe,
		ActivationArguments: args,
		Actions: []toast.Action{
			{
				Type:      toast.Foreground,
				Content:   "Go to Terminal",
				Arguments: args,
			},
		},
	}
	return noti.Push()
}

// HandleCOMActivation is called when Windows relaunches the exe via -Embedding.
// It registers the COM class factory, sets up the activation callback to
// activate the target window, then waits for the callback to fire.
func HandleCOMActivation(timeout time.Duration) {
	if err := ole.RoInitialize(1); err != nil {
		return
	}

	done := make(chan struct{}, 1)

	wintoast.SetActivationCallback(func(appUserModelId, invokedArgs string, userData []wintoast.UserData) {
		hwnd := parseHWND(invokedArgs)
		if hwnd != 0 {
			window.Activate(hwnd)
		}
		select {
		case done <- struct{}{}:
		default:
		}
	})

	var cookie int64
	procRegisterClassObject.Call(
		uintptr(unsafe.Pointer(wintoast.GUID_ImplNotificationActivationCallback)),
		uintptr(unsafe.Pointer(wintoast.ClassFactory)),
		uintptr(ole.CLSCTX_LOCAL_SERVER),
		1, // REGCLS_MULTIPLEUSE
		uintptr(unsafe.Pointer(&cookie)),
	)

	go func() {
		for range time.NewTicker(time.Second).C {
		}
	}()

	select {
	case <-done:
	case <-time.After(timeout):
	}
}

func parseHWND(args string) uintptr {
	if len(args) > 5 && args[:5] == "hwnd:" {
		v, err := strconv.ParseUint(args[5:], 10, 64)
		if err == nil {
			return uintptr(v)
		}
	}
	return 0
}
