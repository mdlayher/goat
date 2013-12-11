package main

import (
	"encoding/json"
	"fmt"
	"goat"
	"os"
	"os/signal"
	"syscall"
)

const APP = "goat"

func main() {
	//load configureation file
	var conf goat.Config
	configFile, err := os.Open("config.json")
	read := json.NewDecoder(configFile)
	err = read.Decode(&conf)
	if err != nil {
		fmt.Println("config.json was unable to be found \nusing default port of 8080")
		conf.Port = "8080"
		conf.Protocals = ["tcp","udp"]
	}
	// Create goroutine to handle termination via UNIX signal
	killChan := make(chan bool)
	doneChan := make(chan int)
	go goat.Manager(killChan, doneChan, conf.Port)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	signal.Notify(sigc, syscall.SIGTERM)
	func() {
		for sig := range sigc {
			fmt.Println(APP, ": caught signal:", sig)
			killChan <- true
			hold := <-doneChan
			os.Exit(hold)
		}
	}()

}
