package common

import (
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// ServerStatus represents a struct to be serialized, containing information about the system running goat
type ServerStatus struct {
	PID          int        `json:"pid"`
	Hostname     string     `json:"hostname"`
	Platform     string     `json:"platform"`
	Architecture string     `json:"architecture"`
	NumCPU       int        `json:"numCpu"`
	NumGoroutine int        `json:"numGoroutine"`
	MemoryMB     float64    `json:"memoryMb"`
	Maintenance  bool       `json:"maintenance"`
	Status       string     `json:"status"`
	Uptime       int64      `json:"uptime"`
	API          TimedStats `json:"api"`
	HTTP         TimedStats `json:"http"`
	UDP          TimedStats `json:"udp"`
}

// GetServerStatus returns the tracker's current status in a ServerStatus struct
func GetServerStatus() (ServerStatus, error) {
	// Get system hostname
	var hostname string
	hostname, err := os.Hostname()
	if err != nil {
		return ServerStatus{}, err
	}

	// Get current memory profile
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	// Report memory usage in MB
	memMb := float64((float64(mem.Alloc) / 1000) / 1000)

	// Current uptime
	uptime := time.Now().Unix() - Static.StartTime

	// API status
	apiStatus := TimedStats{
		atomic.LoadInt64(&Static.API.Minute),
		atomic.LoadInt64(&Static.API.HalfHour),
		atomic.LoadInt64(&Static.API.Hour),
		atomic.LoadInt64(&Static.API.Total),
	}

	// HTTP status
	httpStatus := TimedStats{
		atomic.LoadInt64(&Static.HTTP.Minute),
		atomic.LoadInt64(&Static.HTTP.HalfHour),
		atomic.LoadInt64(&Static.HTTP.Hour),
		atomic.LoadInt64(&Static.HTTP.Total),
	}

	// UDP status
	udpStatus := TimedStats{
		atomic.LoadInt64(&Static.UDP.Minute),
		atomic.LoadInt64(&Static.UDP.HalfHour),
		atomic.LoadInt64(&Static.UDP.Hour),
		atomic.LoadInt64(&Static.UDP.Total),
	}

	// Build status struct
	status := ServerStatus{
		os.Getpid(),
		hostname,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		memMb,
		Static.Maintenance,
		Static.StatusMessage,
		uptime,
		apiStatus,
		httpStatus,
		udpStatus,
	}

	// Return status struct
	return status, nil
}
