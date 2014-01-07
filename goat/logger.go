package goat

import (
	"log"
	"time"
)

// LogManager is responsible for creating the main loggers
func LogManager() {
	// Set up logging flags
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)

	// Start the system status logger
	go StatusLogger()
}

// StatusLogger logs and displays system status at regular intervals
func StatusLogger() {
	ticker := time.NewTicker(5 * time.Minute)

	// Loop infinitely, trigger events via ticker
	for {
		select {
		case <-ticker.C:
			// Fetch status, log it
			go PrintCurrentStatus()
		}
	}
}
