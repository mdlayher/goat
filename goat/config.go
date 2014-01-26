package goat

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/user"
)

// dbConf represents database configuration
type dbConf struct {
	Database string
	Username string
	Password string
}

// sslConf represents SSL configuration
type sslConf struct {
	Port        int
	Certificate string
	Key         string
}

// Conf represents server configuration
type conf struct {
	Port      int
	Passkey   bool
	Whitelist bool
	Interval  int
	HTTP      bool
	HTTPS     bool
	API       bool
	UDP       bool
	Redis     bool
	SSL       sslConf
	DB        dbConf
}

// LoadConfig loads configuration
func loadConfig() conf {
	// Configuration path
	var path string
	config := "config.json"

	// If running on Travis, use alternate configuration
	if os.Getenv("TRAVIS") == "true" {
		config = ".config.travis.json"
	}

	// Store config in standard location
	path = user.HomeDir + "/.config/goat/"

	log.Println("Loading configuration: " + path + config)

	// Check file existence
	_, err = os.Stat(path + config)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Could not find configuration, attempting to create it...")

			if err = os.MkdirAll(path, 0775); err != nil {
				log.Println("Failed to create directory: " + path)
			}

			// Attempt to copy config to home directory
			source, err := os.Open(config)
			if err != nil {
				log.Println("Failed to read source file: " + config)
			}

			// Open destination file
			dest, err := os.Create(path + config)
			if err != nil {
				log.Println("Failed to create destination file: " + path + config)
			}

			// Copy contents
			if _, err = io.Copy(dest, source); err != nil {
				log.Println("Failed to copy to configuration file: " + path + config)
			}

			// Close files
			if err = source.Close(); err != nil {
				log.Println(err.Error())
			}

			if err = dest.Close(); err != nil {
				log.Println(err.Error())
			}
		}
	}

	// Load configuration file
	c := conf{}
	configFile, err := os.Open(path + config)
	if err != nil {
		log.Println("Failed to open config.json")
		return conf{}
	}

	// Decode JSON
	err = json.NewDecoder(configFile).Decode(&c)
	if err != nil {
		log.Println("Could not parse config.json")
		return conf{}
	}

	return c
}
