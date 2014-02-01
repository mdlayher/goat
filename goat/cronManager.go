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
	files := new(data.FileRecordRepository).All()
	if len(files) == 0 {
		log.Println("cronPeerReaper: no files found")
		return
	}

	// Iterate all files
	for _, f := range files {
		// Reap peers on each
		if !f.PeerReaper() {
			log.Println("cronPeerReaper: failed to reap peers on file ID:", f.ID)
		}
	}

	log.Println("cronPeerReaper: complete")
}

// cronPrintCurrentStatus logs the regular status check banner
func cronPrintCurrentStatus() {
	// Grab server status
	stat := common.GetServerStatus()
	if stat == (common.ServerStatus{}) {
		log.Println("Could not print current status")
		return
	}

	// Regular status banner
	log.Printf("status - [goroutines: %d] [memory: %02.3f MB]", stat.NumGoroutine, stat.MemoryMB)

	// HTTP stats
	if common.Static.Config.HTTP {
		log.Printf("  http - [current: %d] [total: %d]", stat.HTTP.Current, stat.HTTP.Total)

		// Reset current HTTP counter
		atomic.StoreInt64(&common.Static.HTTP.Current, 0)
	}

	// UDP stats
	if common.Static.Config.UDP {
		log.Printf("   udp - [current: %d] [total: %d]", stat.UDP.Current, stat.UDP.Total)

		// Reset current UDP counter
		atomic.StoreInt64(&common.Static.UDP.Current, 0)
	}
}
