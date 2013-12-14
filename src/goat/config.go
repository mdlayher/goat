package goat

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Port      string
	Passkey   bool
	Http      bool
	Udp       bool
	Map       bool
	Sql       bool
	CacheSize int //must be an int between 1-8. Higher values result in
	//faster performance and larger memory requirements
}

// Load configuration
func LoadConfig() Config {
	// Read in JSON file
	var conf Config
	file, err := filepath.Abs("config.json")
	configFile, err := os.Open(file)
	read := json.NewDecoder(configFile)

	// Decode JSON
	err = read.Decode(&conf)
	if err != nil {
		Static.LogChan <- "config.json could not be read, using default configuration"
		conf.Port = "8080"
		conf.Passkey = false
		conf.Http = true
		conf.Udp = false
		conf.Map = true
		conf.Sql = false
		conf.CacheSize = 3
	}

	return conf
}
