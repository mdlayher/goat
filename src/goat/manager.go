package goat

import (
	"fmt"
)

func Manager(killChan chan bool, doneChan chan int, port string) {
	go new(HttpListener).Listen(port)

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
