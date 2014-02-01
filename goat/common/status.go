package common

import (
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// ServerStatus represents a struct to be serialized, containing information about the system running goat
type ServerStatus struct {
	PID          int       `json:"pid"`
	Hostname     string    `json:"hostname"`
	Platform     string    `json:"platform"`
	Architecture string    `json:"architecture"`
	NumCPU       int       `json:"numCpu"`
	NumGoroutine int       `json:"numGoroutine"`
	MemoryMB     float64   `json:"memoryMb"`
	Uptime       int64     `json:"uptime"`
	HTTP         HTTPStats `json:"http"`
	UDP          UDPStats  `json:"udp"`
}

// GetServerStatus represents a tracker status request
func GetServerStatus() ServerStatus {
	// Get system hostname
	var hostname string
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err.Error())
		return ServerStatus{}
	}

	// Get current memory profile
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	// Report memory usage in MB
	memMb := float64((float64(mem.Alloc) / 1000) / 1000)

	// Current uptime
	uptime := time.Now().Unix() - Static.StartTime

	// HTTP status
	httpStatus := HTTPStats{
		atomic.LoadInt64(&Static.HTTP.Current),
		atomic.LoadInt64(&Static.HTTP.Total),
	}

	// UDP status
	udpStatus := UDPStats{
		atomic.LoadInt64(&Static.UDP.Current),
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
		uptime,
		httpStatus,
		udpStatus,
	}

	// Return status struct
	return status
}
