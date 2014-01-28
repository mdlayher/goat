package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mdlayher/goat/goat"
)

// config is a flag which allows override of the default configuration file location
var config = flag.String("config", "", "Override config file location with custom path.")

// mySQLDSN is a flag which allows manual override of MySQL DSN configuration
var mySQLDSN = flag.String("mysqldsn", "", "Override config file with the following MySQL DSN configuration.")

// test is a flag which causes goat to start, and exit shortly after
var test = flag.Bool("test", false, "Make goat start, and exit shortly after. Used for testing.")

func main() {
	// Set up command line options
	flag.Parse()

	// Pass command-line flags to override default configurations
	goat.ConfigPath = config
	goat.MySQLDSN = mySQLDSN

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
