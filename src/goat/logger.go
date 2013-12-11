package goat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func LogMng(doneChan chan bool, logChan chan string) {
	// Create log directory and file, and pull current date to add to logfile name
	now := time.Now()
	os.Mkdir("logs", os.ModeDir|os.ModePerm)
	logFile, err := os.Create(fmt.Sprintf("logs/goat-%d-%d-%d.log", now.Year(), now.Month(), now.Day()))
	if err != nil {
		fmt.Println(err)
	}

	// create a logger that will use the writer created above
	logger := log.New(bufio.NewWriter(logFile), "", log.Lmicroseconds|log.Lshortfile)
	amIDone := false
	msg := ""

	// Wait for error to be passed on the logChan channel, or termination signal on doneChan
	for !amIDone {
		select {
		case amIDone = <-doneChan:
			logFile.Close()
		case msg = <-logChan:
			logger.Println(APP, ":", msg)
			fmt.Println(APP, ":", msg)
		}
	}
}
