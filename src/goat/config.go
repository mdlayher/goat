package goat

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Load configuration
func LoadConfig() Conf {
	// Read in JSON file
	var conf Conf
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
		conf.Size = 10000000000
	}

	return conf
}
