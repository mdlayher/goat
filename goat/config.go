package goat

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/user"
)

// DBConf represents database configuration
type DBConf struct {
	Database string
	Username string
	Password string
}

// Conf represents server configuration
type Conf struct {
	Port      int
	Passkey   bool
	Whitelist bool
	Interval  int
	HTTP      bool
	UDP       bool
	DB        DBConf
}

// LoadConfig loads configuration
func LoadConfig() Conf {
	// Configuration path
	var path string
	config := "config.json"

	// Load current user from OS, to get home directory
	user, err := user.Current()
	if err != nil {
		log.Println(err.Error())
		path = "."
	} else {
		// Store config in standard location
		path = user.HomeDir + "/.config/goat/"
	}

	log.Println("Loading configuration: " + path + config)

	// Check file existence
	_, err = os.Stat(path + config)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Could not find configuration, attempting to create it...")

			err = os.MkdirAll(path, 0775)
			if err != nil {
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
			_, err = io.Copy(dest, source)
			if err != nil {
				log.Println("Failed to copy to configuration file: " + path + config)
			}

			// Close files
			source.Close()
			dest.Close()
		}
	}

	// Load configuration file
	conf := Conf{}
	configFile, err := os.Open(path + config)
	if err != nil {
		log.Println("Failed to open config.json")
		return Conf{}
	}

	// Decode JSON
	err = json.NewDecoder(configFile).Decode(&conf)
	if err != nil {
		log.Println("Could not parse config.json")
		return Conf{}
	}

	return conf
}
