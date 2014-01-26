package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChimeraCoder/telescope/goat"
)

var test = flag.Bool("test", false, "Make goat start, and exit shortly after. Used for testing.")
var redisPass = flag.String("redispass", "", "password for the Redis database, if any")
var mysqlDSN = flag.String("mysqldsn", "", "msql data source name")

func main() {
	// Set up command line options
	flag.Parse()

	// Set password for Redis database and DSN for MySQL database
	goat.RedisPass = redisPass
	goat.MysqlDSN = mysqlDSN

	// If test mode, trigger quit shortly after startup
	// Used for CI tests, so that we ensure goat starts up and is able to stop gracefully
	if *test {
		go func() {
			fmt.Println(goat.App, ": launched in test mode")
			time.Sleep(5 * time.Second)

			fmt.Println(goat.App, ": test mode triggering graceful shutdown")
			err := syscall.Kill(os.Getpid(), syscall.SIGTERM)
			if err != nil {
				fmt.Println(goat.App, ": failed to invoke graceful shutdown, halting")
				os.Exit(1)
			}
		}()
	}

	// Launch manager via goroutine
	killChan := make(chan bool)
	exitChan := make(chan int)
	go goat.Manager(killChan, exitChan)

	// Gracefully handle termination via UNIX signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	for sig := range sigChan {
		// Trigger manager shutdown
		fmt.Println(goat.App, ": caught signal:", sig)
		killChan <- true
		break
	}

	// Force terminate if signaled twice
	go func(sigChan chan os.Signal) {
		for sig := range sigChan {
			_ = sig
			fmt.Println(goat.App, ": force halting now!")
			os.Exit(1)
		}
	}(sigChan)

	// Exit with specified code from manager
	code := <-exitChan
	fmt.Println(goat.App, ": graceful shutdown complete")
	os.Exit(code)
}
