package main

import (
	"encoding/json"
	"fmt"
	"goat"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const APP = "goat"

func main() {
	// Read configuration
	conf := load_config()

	// Launch manager via goroutine
	killChan := make(chan bool)
	doneChan := make(chan int)
	log.Println("hello")
	go goat.Manager(killChan, doneChan, conf.Port)

	// Gracefully handle termination via UNIX signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	signal.Notify(sigc, syscall.SIGTERM)
	for sig := range sigc {
		fmt.Println(APP, ": caught signal:", sig)
		killChan <- true
		hold := <-doneChan
		os.Exit(hold)
	}
}

// Load configuration
func load_config() goat.Config {
	// Read in JSON file
	var conf goat.Config
	configFile, err := os.Open("config.json")
	read := json.NewDecoder(configFile)

	// Decode JSON
	err = read.Decode(&conf)
	if err != nil {
		fmt.Println(APP, ": config.json could not be read, using default configuration")
		conf.Port = "8080"
		conf.Http = true
		conf.Udp = false

	}

	return conf
}
