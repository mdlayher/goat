package goat

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// cronManager spawns and triggers events at regular intervals
func cronManager() {
	// Run on startup
	go cronPeerReaper()

	// cronPeerReaper - run at regular announce interval
	peerReaper := time.NewTicker(time.Duration(common.Static.Config.Interval) * time.Second)

	// cronPrintCurrentStatus - run every 5 minutes
	status := time.NewTicker(5 * time.Minute)

	// Start cronStatsReset, which maintains its own timers
	go cronStatsReset()

	// Trigger events via ticker
	for {
		select {
		case <-peerReaper.C:
			go cronPeerReaper()
		case <-status.C:
			go cronPrintCurrentStatus()
		}
	}
}

// cronPeerReaper checks for inactive peers, and marks them as such in the database
func cronPeerReaper() {
	log.Println("cronPeerReaper: starting")

	// Load all files
	files, err := new(data.FileRecordRepository).All()
	if err != nil {
		log.Println(err.Error())
		log.Println("cronPeerReaper: failed to load list of files")
		return
	}

	if len(files) == 0 {
		log.Println("cronPeerReaper: no files found")
		return
	}

	// Sum of peers reaped
	total := 0

	// Iterate all files
	for _, f := range files {
		// Reap peers on each
		count, err := f.PeerReaper()
		if err != nil {
			log.Println("cronPeerReaper: failed to reap peers on file ID:", f.ID)
		}

		// Sum peers reaped
		if count > 0 {
			total += count
			log.Printf("cronPeerReaper: reaped %d peers on file ID: %d", f.ID)
		}
	}

	log.Printf("cronPeerReaper: complete, reaped %d peers on %d files", total, len(files))
}

// cronPrintCurrentStatus logs the regular status check banner
func cronPrintCurrentStatus() {
	// Grab server status
	stat, err := common.GetServerStatus()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Regular status banner
	log.Printf("status - [goroutines: %d] [memory: %02.3f MB]", stat.NumGoroutine, stat.MemoryMB)

	// API stats
	if common.Static.Config.API {
		log.Printf("   api - [1 min: %03d | 30 min: %03d | 60 min: %03d] [total: %03d]", stat.API.Minute, stat.API.HalfHour, stat.API.Hour, stat.API.Total)
	}

	// HTTP stats
	if common.Static.Config.HTTP {
		log.Printf("  http - [1 min: %03d | 30 min: %03d | 60 min: %03d] [total: %03d]", stat.HTTP.Minute, stat.HTTP.HalfHour, stat.HTTP.Hour, stat.HTTP.Total)
	}

	// UDP stats
	if common.Static.Config.UDP {
		log.Printf("   udp - [1 min: %03d | 30 min: %03d | 60 min: %03d] [total: %03d]", stat.UDP.Minute, stat.UDP.HalfHour, stat.UDP.Hour, stat.UDP.Total)
	}
}

// cronStatsReset triggers a reset of certain statistic counters at regular intervals
func cronStatsReset() {
	// Trigger events every hour
	hour := time.NewTicker(1 * time.Hour)

	// Trigger events every half hour
	halfHour := time.NewTicker(30 * time.Minute)

	// Trigger events every minute
	minute := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-hour.C:
			// Reset hourly stats counters
			atomic.StoreInt64(&common.Static.API.Hour, 0)
			atomic.StoreInt64(&common.Static.HTTP.Hour, 0)
			atomic.StoreInt64(&common.Static.UDP.Hour, 0)
		case <-halfHour.C:
			// Reset half-hourly stats counters
			atomic.StoreInt64(&common.Static.API.HalfHour, 0)
			atomic.StoreInt64(&common.Static.HTTP.HalfHour, 0)
			atomic.StoreInt64(&common.Static.UDP.HalfHour, 0)
		case <-minute.C:
			// Reset minute stats counters
			atomic.StoreInt64(&common.Static.API.Minute, 0)
			atomic.StoreInt64(&common.Static.HTTP.Minute, 0)
			atomic.StoreInt64(&common.Static.UDP.Minute, 0)
		}
	}
}
