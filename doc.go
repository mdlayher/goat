/*
Command goat provides an implementation of a BitTorrent tracker, written in Go.

Installation

goat can be built using Go 1.1+. It can be downloaded, built, and installed,
simply by running 'go get github.com/mdlayher/goat'.

In addition, goat depends on a MySQL server for data storage.  After creating a
database and user for goat, its database schema may be imported from the SQL
files located in 'res/'.  goat will not run unless MySQL is installed, and a
database and user are properly configured for its use.

Configuration

goat is configured using a JSON file, which will be created under
'~/.config/goat/config.json' on UNIX systems.  Here is an example configuration,
describing the settings available to the user.

	{
		// Port: the port number on which goat will listen using both HTTP and UDP interfaces
		"Port": 8080,

		// Passkey: require that a valid passkey is present in announce/scrape URLs (HTTP only)
		// note: this setting is typically used only for private trackers
		// ex: http://localhost:8080/0123456789ABCDEF/announce
		"Passkey": true,

		// Whitelist: require clients to be whitelisted manually for use with the tracker
		// note: this setting is typically used only for private trackers
		"Whitelist": true,

		// Interval: approximately how often, in seconds, clients should announce to the tracker
		// note: this entropy is introduced to stagger time between many client announces
		"Interval": 3600,

		// HTTP: enable listening for client connections via HTTP
		"HTTP": true,

		// API: enable a HTTP RESTful API, used to pull statistics from goat
		// note: only enabled when HTTP is enabled
		"API": true,

		// UDP: enable listening for client connections via UDP
		// note: it is not possible to use a passkey with this listener, so this listener should
		// only be used for public trackers
		"UDP": false,

		// Db: MySQL database configuration
		"DB": {
			// Database: the database goat will use to store its tracker data
			"Database": "goat",

			// Username: the username used to access goat's database
			"Username": "goat",

			// Password: the password used to access goat's database
			"Password": "goat"
		}
	}

*/
package main
