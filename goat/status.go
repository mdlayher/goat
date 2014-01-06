package goat

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
)

// ServerStatus represents a struct to be serialized, containing information about the system running goat
type ServerStatus struct {
	PID          int     `json:"pid"`
	Hostname     string  `json:"hostname"`
	Platform     string  `json:"platform"`
	Architecture string  `json:"architecture"`
	NumCPU       int     `json:"numCpu"`
	NumGoroutine int     `json:"numGoroutine"`
	MemoryMB     float64 `json:"memoryMb"`
	HTTPTotal    int64   `json:"httpTotal"`
	HTTPCurrent  int64   `json:"httpCurrent"`
}

// GetServerStatus represents a tracker status request
func GetServerStatus() ServerStatus {
	// Get system hostname
	var hostname string
	hostname, err := os.Hostname()
	if err != nil {
		Static.LogChan<- err.Error()
		return ServerStatus{}
	}

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
		atomic.LoadInt64(&Static.HTTP.Total),
		atomic.LoadInt64(&Static.HTTP.Current),
	}

	// Return status struct
	return status
}

// GetStatusJSON returns a JSON representation of server status
func GetStatusJSON(resChan chan []byte) {
	// Marshal into JSON from request
	res, err := json.Marshal(GetServerStatus())
	if err != nil {
		resChan <- nil
	}

	// Return status
	resChan <- res
}

// PrintStatusBanner logs the startup status banner
func PrintStatusBanner() {
	// Grab initial server status
	stat := GetServerStatus()
	if stat == (ServerStatus{}) {
		Static.LogChan <- "Could not print startup status banner"
		return
	}

	// Startup banner
	Static.LogChan <- fmt.Sprintf("%s - %s_%s (%d CPU) [pid: %d]", stat.Hostname, stat.Platform, stat.Architecture, stat.NumCPU, stat.PID)
}

// PrintCurrentStatus logs the regular status check banner
func PrintCurrentStatus() {
	// Grab server status
	stat := GetServerStatus()
	if stat == (ServerStatus{}) {
		Static.LogChan <- "Could not print current status"
		return
	}

	// Regular status banner
	Static.LogChan <- fmt.Sprintf("status - [goroutines: %d] [memory: %02.3f MB]", stat.NumGoroutine, stat.MemoryMB)

	// HTTP stats
	if Static.Config.HTTP {
		Static.LogChan <- fmt.Sprintf("  http - [current: %d] [total: %d]", stat.HTTPCurrent, stat.HTTPTotal)

		// Reset current HTTP counter
		atomic.StoreInt64(&Static.HTTP.Current, 0)
	}
}
