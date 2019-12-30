package utils

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"os"
	"syscall"
)

// GenID returns a unique ID
// TODO: maybe distinguish between types of IDs, runs, processors, etc
func GenID() string {
	return ksuid.New().String()
}

// From https://github.com/shirou/gopsutil/blob/c141152a7b8f59b63e060fa8450f5cd5e7196dfb/process/process_posix.go#L73
func PidExists(pid int32) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("invalid pid %v", pid)
	}
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return false, err
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil
	}
	if err.Error() == "os: process already finished" {
		return false, nil
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false, err
	}
	switch errno {
	case syscall.ESRCH:
		return false, nil
	case syscall.EPERM:
		return true, nil
	}
	return false, err
}
