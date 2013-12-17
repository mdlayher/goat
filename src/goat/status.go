package goat

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
)

// Struct to be serialized, containing information about the system running goat
type ServerStatus struct {
	Pid          int     `json:"pid"`
	Hostname     string  `json:"hostname"`
	Platform     string  `json:"platform"`
	Architecture string  `json:"architecture"`
	NumCpu       int     `json:"numCpu"`
	NumGoroutine int     `json:"numGoroutine"`
	MemoryMb     float64 `json:"memoryMb"`
	HttpTotal    int64   `json:"httpTotal"`
	HttpCurrent  int64   `json:"httpCurrent"`
}

// Tracker status request
func GetServerStatus() ServerStatus {
	// Get system hostname
	hostname, _ := os.Hostname()

	// Get current memory profile
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	// Report memory usage in MB
	memMb := float64((float64(mem.Alloc) / 1000) / 1000)

	// Build status struct
	status := ServerStatus{
		os.Getpid(),
		hostname,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		memMb,
		atomic.LoadInt64(&Static.Http.Total),
		atomic.LoadInt64(&Static.Http.Current),
	}

	// Return status struct
	return status
}

// Return JSON representation of server status
func GetStatusJson(resChan chan []byte) {
	// Marshal into JSON from request
	res, err := json.Marshal(GetServerStatus())
	if err != nil {
		resChan <- nil
	}

	// Return status
	resChan <- res
}

// Log the startup status banner
func PrintStatusBanner() {
	// Grab initial server status
	stat := GetServerStatus()

	// Startup banner
	Static.LogChan <- fmt.Sprintf("%s - %s_%s (%d CPU) [pid: %d]", stat.Hostname, stat.Platform, stat.Architecture, stat.NumCpu, stat.Pid)
}

// Log the regular status check banner
func PrintCurrentStatus() {
	// Grab server status
	stat := GetServerStatus()

	// Regular status banner
	Static.LogChan <- fmt.Sprintf("status - [goroutines: %d] [memory: %02.3f MB]", stat.NumGoroutine, stat.MemoryMb)

	// HTTP stats
	if Static.Config.Http {
		Static.LogChan <- fmt.Sprintf("  http - [current: %d] [total: %d]", stat.HttpCurrent, stat.HttpTotal)

		// Reset current HTTP counter
		atomic.StoreInt64(&Static.Http.Current, 0)
	}
}
