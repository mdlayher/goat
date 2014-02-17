package goat

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// cronManager spawns and triggers events at regular intervals
func cronManager() {
	// Run on startup
	go cronAPIKeyReaper()
	go cronPeerReaper()

	// cronAPIKeyReaper - run once per hour
	apiKeyReaper := time.NewTicker(1 * time.Hour)

	// cronPeerReaper - run at regular announce interval
	peerReaper := time.NewTicker(time.Duration(common.Static.Config.Interval) * time.Second)

	// cronPrintCurrentStatus - run every 5 minutes
	status := time.NewTicker(5 * time.Minute)

	// Start cronStatsReset, which maintains its own timers
	go cronStatsReset()

	// Trigger events via ticker
	for {
		select {
		case <-apiKeyReaper.C:
			go cronAPIKeyReaper()
		case <-peerReaper.C:
			go cronPeerReaper()
		case <-status.C:
			go cronPrintCurrentStatus()
		}
	}
}

// cronAPIKeyReaper checks for expired API keys, and deletes them from the database
func cronAPIKeyReaper() {
	log.Println("cronAPIKeyReaper: starting")

	// Load all API keys
	keys, err := new(data.APIKeyRepository).All()
	if err != nil {
		log.Println(err.Error())
		log.Println("cronAPIKeyReaper: failed to list list of API keys")
		return
	}

	if len(keys) == 0 {
		log.Println("cronAPIKeyReaper: no API keys found")
		return
	}

	// Sum of keys reaped
	var total int64
	atomic.StoreInt64(&total, 0)

	// WaitGroup to wait for all keys to finish being reaped
	var wg sync.WaitGroup
	wg.Add(len(keys))

	// Iterate all keys in parallel
	for _, k := range keys {
		go func(k data.APIKey, count *int64, wg *sync.WaitGroup) {
			// Check for expired key
			if k.Expire <= time.Now().Unix() {
				// Delete expired keys
				if err := k.Delete(); err != nil {
					log.Println(err.Error())
				}

				// Increment reap counter
				atomic.AddInt64(count, 1)
			}

			// Inform WaitGroup this goroutine is done
			wg.Done()
		}(k, &total, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Printf("cronAPIKeyReaper: complete, reaped %d/%d keys", total, len(keys))
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
	var total int64
	atomic.StoreInt64(&total, 0)

	// WaitGroup to wait for all peers to finish being reaped
	var wg sync.WaitGroup
	wg.Add(len(files))

	// Iterate all files in parallel
	for _, f := range files {
		go func(f data.FileRecord, count *int64, wg *sync.WaitGroup) {
			// Reap peers on each file
			reaped, err := f.PeerReaper()
			if err != nil {
				log.Println("cronPeerReaper: failed to reap peers on file ID:", f.ID)
			}

			// Increment reap counter
			atomic.AddInt64(count, int64(reaped))

			if reaped > 0 {
				log.Printf("cronPeerReaper: reaped %d peers on file ID: %d", reaped, f.ID)
			}

			// Inform WaitGroup this goroutine is done
			wg.Done()
		}(f, &total, &wg)
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
