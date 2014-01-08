package goat

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"sync/atomic"
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
	HTTP         HTTPStats `json:"http"`
	UDP          UDPStats  `json:"udp"`
}

// HTTPStats represents statistics regarding HTTP server
type HTTPStats struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

// UDPStats represents statistics regarding UDP server
type UDPStats struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
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
		httpStatus,
		udpStatus,
	}

	// Return status struct
	return status
}

// GetStatusJSON returns a JSON representation of server status
func GetStatusJSON(resChan chan []byte) {
	// Marshal into JSON from request
	res, err := json.Marshal(GetServerStatus())
	if err != nil {
		log.Println(err.Error())
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
		log.Println("Could not print startup status banner")
		return
	}

	// Startup banner
	log.Printf("%s - %s_%s (%d CPU) [pid: %d]", stat.Hostname, stat.Platform, stat.Architecture, stat.NumCPU, stat.PID)
}

// PrintCurrentStatus logs the regular status check banner
func PrintCurrentStatus() {
	// Grab server status
	stat := GetServerStatus()
	if stat == (ServerStatus{}) {
		log.Println("Could not print current status")
		return
	}

	// Regular status banner
	log.Printf("status - [goroutines: %d] [memory: %02.3f MB]", stat.NumGoroutine, stat.MemoryMB)

	// HTTP stats
	if Static.Config.HTTP {
		log.Printf("  http - [current: %d] [total: %d]", stat.HTTP.Current, stat.HTTP.Total)

		// Reset current HTTP counter
		atomic.StoreInt64(&Static.HTTP.Current, 0)
	}

	// UDP stats
	if Static.Config.UDP {
		log.Printf("   udp - [current: %d] [total: %d]", stat.UDP.Current, stat.UDP.Total)

		// Reset current UDP counter
		atomic.StoreInt64(&Static.UDP.Current, 0)
	}
}
