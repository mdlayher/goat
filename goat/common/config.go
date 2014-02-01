package common

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/user"
	ospath "path"
)

// ConfigPath is set via command-line, and can be used to override config file path location
var ConfigPath *string

// dbConf represents database configuration
type dbConf struct {
	Host     string
	Database string
	Username string
	Password string
}

// sslConf represents SSL configuration
type sslConf struct {
	Enabled     bool
	Port        int
	Certificate string
	Key         string
}

// redisConf represents Redis configuration
type redisConf struct {
	Enabled  bool
	Host     string
	Password string
}

// Conf represents server configuration
type Conf struct {
	Port      int
	Passkey   bool
	Whitelist bool
	Interval  int
	HTTP      bool
	API       bool
	UDP       bool
	SSL       sslConf
	DB        dbConf
	Redis     redisConf
}

// LoadConfig loads configuration
func LoadConfig() Conf {
	// Configuration path
	var path string
	config := "config.json"

	// If running on Travis, use alternate configuration
	if os.Getenv("TRAVIS") == "true" {
		config = ".config.travis.json"
	}

	// Load current user from OS, to get home directory
	user, err := user.Current()
	if err != nil {
		log.Println(err.Error())
		path = "./"
	} else {
		// Store config in standard location
		path = user.HomeDir + "/.config/goat/"
	}

	// Allow manual override of config path, if flag is set
	if ConfigPath != nil && *ConfigPath != "" {
		// Split config path into path and filename
		path = ospath.Dir(*ConfigPath) + "/"
		config = ospath.Base(*ConfigPath)
	}

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
	c := Conf{}
	configFile, err := os.Open(path + config)
	if err != nil {
		log.Println("Failed to open config.json")
		return Conf{}
	}

	// Decode JSON
	err = json.NewDecoder(configFile).Decode(&c)
	if err != nil {
		log.Println("Could not parse config.json")
		return Conf{}
	}

	return c
}
