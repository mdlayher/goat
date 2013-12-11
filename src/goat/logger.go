package goat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func LogMng(done chan bool, message chan string) {
	currentTime := time.Now().String()
	logFile, err := os.Create("GoatLog" + currentTime + ".log")
	if err != nil {
		fmt.Println(err)
	}
	writer := bufio.NewWriter(logFile)
	logger := log.New(writer, "", log.Lmicroseconds|log.Lshortfile)

	amIDone := false
	msg := ""

	for !amIDone {
		select {
		case amIDone = <-done:
			continue
		case msg = <-message:
			logger.Println(msg)

		}
	}

}
