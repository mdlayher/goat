package goat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func LogMng(doneChan chan bool, logChan chan string) {
	// create log file and pull current time to add to logfile name
	currentTime := time.Now().String()
	logFile, err := os.Create("GoatLog" + currentTime + ".log")
	if err != nil {
		fmt.Println(err)
	}
	writer := bufio.NewWriter(logFile)
	// create a logger that will use the writer created above
	logger := log.New(writer, "", log.Lmicroseconds|log.Lshortfile)
	amIDone := false
	msg := ""
	// wait for errer to be passed on the logChan channel or the done chan
	for !amIDone {
		select {
		case amIDone = <-doneChan:
			logFile.Close()
		case msg = <-logChan:
			logger.Println(msg)
		}
	}
}
