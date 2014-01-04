package goat

import (
	"encoding/json"
	"io"
	"os"
	"os/user"
)

// Conf represents server configuration
type Conf struct {
	Port      int
	Passkey   bool
	Whitelist bool
	Interval  int
	HTTP      bool
	UDP       bool
}

// LoadConfig loads configuration
func LoadConfig() Conf {
	// Load current user from OS, to get home directory
	user, err := user.Current()
	if err != nil {
		Static.LogChan <- err.Error()
	}

	// Configuration path
	path := user.HomeDir + "/.config/goat/"
	config := "config.json"

	Static.LogChan <- "Loading configuration: " + path + config

	// Check file existence
	_, err = os.Stat(path + config)
	if err != nil {
		if os.IsNotExist(err) {
			Static.LogChan <- "Could not find configuration, attempting to create it..."

			err = os.MkdirAll(path, 0775)
			if err != nil {
				Static.LogChan <- "Failed to create directory: " + path
			}

			// Attempt to copy config to home directory
			source, err := os.Open(config)
			if err != nil {
				Static.LogChan <- "Failed to read source file: " + config
			}

			// Open destination file
			dest, err := os.Create(path + config)
			if err != nil {
				Static.LogChan <- "Failed to create destination file: " + path + config
			}

			// Copy contents
			_, err = io.Copy(dest, source)
			if err != nil {
				Static.LogChan <- "Failed to copy to configuration file: " + path + config
			}

			// Close files
			source.Close()
			dest.Close()
		}
	}

	// Load configuration file
	configFile, err := os.Open(path + config)

	// Decode JSON
	var conf Conf
	err = json.NewDecoder(configFile).Decode(&conf)
	if err != nil {
		Static.LogChan <- "Could not read config.json, using defaults..."

		// Sane configuration defaults
		conf.Port = 8080
		conf.Passkey = true
		conf.Whitelist = true
		conf.Interval = 3600
		conf.HTTP = true
		conf.UDP = false
	}

	return conf
}
