package goat

import (
	"fmt"
)

const APP = "goat"

func Manager(killChan chan bool, doneChan chan int, port string) {
	// Launch listeners
	go new(HttpListener).Listen(port)
	go new(UdpListener).Listen(port)

	fmt.Println(APP, ": HTTP and UDP listeners launched on port " + port)

	for {
		select {
		case <-killChan:
			//change this to kill workers gracefully and exit
			fmt.Println("done")
			doneChan <- 0
			// case freeWorker := <-ioReturn:
		}
	}
}
