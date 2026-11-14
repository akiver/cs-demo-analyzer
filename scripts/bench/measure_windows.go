//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
	"unsafe"
)

const (
	processQueryLimitedInformation = 0x1000
	processVMRead                  = 0x0010
)

// https://learn.microsoft.com/en-us/windows/win32/api/psapi/ns-psapi-process_memory_counters
type processMemoryCounters struct {
	cb                         uint32
	pageFaultCount             uint32
	peakWorkingSetSize         uintptr
	workingSetSize             uintptr
	quotaPeakPagedPoolUsage    uintptr
	quotaPagedPoolUsage        uintptr
	quotaPeakNonPagedPoolUsage uintptr
	quotaNonPagedPoolUsage     uintptr
	pagefileUsage              uintptr
	peakPagefileUsage          uintptr
}

var procGetProcessMemoryInfo = syscall.NewLazyDLL("psapi.dll").NewProc("GetProcessMemoryInfo")

func runAndMeasure(cmd *exec.Cmd) (measure, error) {
	start := time.Now()
	if err := cmd.Start(); err != nil {
		return measure{}, err
	}

	// Our own handle keeps the process object, and so its memory counters, alive after it exits.
	// The PID can't be reused before Wait because os.Process holds a handle to the process too.
	handle, err := syscall.OpenProcess(processQueryLimitedInformation|processVMRead, false, uint32(cmd.Process.Pid))
	if err != nil {
		_ = cmd.Wait()
		return measure{}, fmt.Errorf("OpenProcess: %w", err)
	}
	defer syscall.CloseHandle(handle)

	if err := cmd.Wait(); err != nil {
		return measure{}, err
	}
	wall := time.Since(start)

	var counters processMemoryCounters
	counters.cb = uint32(unsafe.Sizeof(counters))
	ok, _, err := procGetProcessMemoryInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&counters)), uintptr(counters.cb))
	if ok == 0 {
		return measure{}, fmt.Errorf("GetProcessMemoryInfo: %w", err)
	}

	return measure{wall: wall, peak: uint64(counters.peakWorkingSetSize)}, nil
}
