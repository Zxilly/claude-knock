//go:build linux

package window

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// BuildParentChain returns PIDs from startPID up to the root process
// by reading /proc/<pid>/stat.
func BuildParentChain(startPID uint32) ([]uint32, error) {
	var chain []uint32
	visited := make(map[uint32]bool)
	pid := startPID
	for {
		if visited[pid] {
			break
		}
		visited[pid] = true
		chain = append(chain, pid)
		ppid, err := getParentPID(pid)
		if err != nil || ppid == 0 || ppid == pid {
			break
		}
		pid = ppid
	}
	return chain, nil
}

func getParentPID(pid uint32) (uint32, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, err
	}
	// Format: pid (comm) state ppid ...
	// Find the closing ')' to skip over the comm field which may contain spaces.
	s := string(data)
	idx := strings.LastIndex(s, ")")
	if idx < 0 || idx+2 >= len(s) {
		return 0, fmt.Errorf("unexpected /proc/%d/stat format", pid)
	}
	fields := strings.Fields(s[idx+2:])
	if len(fields) < 2 {
		return 0, fmt.Errorf("unexpected /proc/%d/stat format", pid)
	}
	ppid, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(ppid), nil
}
