//go:build windows

package window

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type processEntry struct {
	Size              uint32
	Usage             uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	Threads           uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [windows.MAX_PATH]uint16
}

// BuildParentChain returns PIDs from startPID up to root process.
func BuildParentChain(startPID uint32) ([]uint32, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot: %w", err)
	}
	defer windows.CloseHandle(snap)

	// Build a map of PID -> parent PID.
	parentMap := make(map[uint32]uint32)
	var entry processEntry
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = process32First(snap, &entry)
	if err != nil {
		return nil, fmt.Errorf("Process32First: %w", err)
	}

	for {
		parentMap[entry.ProcessID] = entry.ParentProcessID
		entry.Size = uint32(unsafe.Sizeof(entry))
		err = process32Next(snap, &entry)
		if err != nil {
			break
		}
	}

	// Walk up from startPID.
	var chain []uint32
	visited := make(map[uint32]bool)
	pid := startPID
	for {
		if visited[pid] {
			break // avoid cycles
		}
		visited[pid] = true
		chain = append(chain, pid)
		parent, ok := parentMap[pid]
		if !ok || parent == 0 || parent == pid {
			break
		}
		pid = parent
	}
	return chain, nil
}

var (
	modkernel32    = windows.NewLazySystemDLL("kernel32.dll")
	procProcess32F = modkernel32.NewProc("Process32FirstW")
	procProcess32N = modkernel32.NewProc("Process32NextW")
)

func process32First(snap windows.Handle, entry *processEntry) error {
	r1, _, e1 := procProcess32F.Call(uintptr(snap), uintptr(unsafe.Pointer(entry)))
	if r1 == 0 {
		return e1
	}
	return nil
}

func process32Next(snap windows.Handle, entry *processEntry) error {
	r1, _, e1 := procProcess32N.Call(uintptr(snap), uintptr(unsafe.Pointer(entry)))
	if r1 == 0 {
		return e1
	}
	return nil
}
