//go:build darwin

package window

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// BuildParentChain returns PIDs from startPID up to the root process
// by using ps to look up parent PIDs.
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
	out, err := exec.Command("ps", "-o", "ppid=", "-p", fmt.Sprint(pid)).Output()
	if err != nil {
		return 0, err
	}
	ppid, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(ppid), nil
}
