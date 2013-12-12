package goat

import (
	"fmt"
)

const APP = "goat"

func Manager(killChan chan bool, exitChan chan int, port string) {
	// Launch listeners
	logChan := make(chan string)
	doneChan := make(chan bool)
	go new(HttpListener).Listen(port, logChan)
	go new(UdpListener).Listen(port, logChan)
	go LogMng(doneChan, logChan)

	fmt.Println(APP, ": HTTP and UDP listeners launched on port "+port)

	for {
		select {
		case <-killChan:
			//change this to kill workers gracefully and exit
			logChan <- "done"
			doneChan <- true
			exitChan <- 0
		}
	}
}
