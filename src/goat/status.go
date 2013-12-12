package goat

import (
	"encoding/json"
	"os"
	"runtime"
)

// Struct to be serialized, containing information about the system running goat
type ServerStatus struct {
	Pid int
	Hostname string
	Platform string
	Architecture string
	NumCpu int
	NumGoroutine int
	MemoryMb float64
}

// Tracker status request
func GetServerStatus(resChan chan []byte) {
	// Get system hostname
	hostname, _ := os.Hostname()

	// Get current memory profile
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	// Report memory usage in MB
	memMb := float64((float64(mem.Alloc) / 1000) / 1000)

	res, err := json.Marshal(ServerStatus{
		os.Getpid(),
		hostname,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		memMb,
	})
	if err != nil {
		resChan <- nil
	}

	// Return JSON
	resChan <- res
}
