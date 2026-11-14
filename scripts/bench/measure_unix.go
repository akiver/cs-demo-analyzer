//go:build !windows

package main

import (
	"errors"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

func runAndMeasure(cmd *exec.Cmd) (measure, error) {
	start := time.Now()
	if err := cmd.Run(); err != nil {
		return measure{}, err
	}
	wall := time.Since(start)

	usage, ok := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	if !ok {
		return measure{}, errors.New("peak memory is not available on this OS")
	}

	// Maxrss is in bytes on macOS and in kilobytes on Linux and the BSDs.
	peak := uint64(usage.Maxrss)
	if runtime.GOOS != "darwin" {
		peak *= 1024
	}

	return measure{wall: wall, peak: peak}, nil
}
