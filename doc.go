/*
Command goat provides an implementation of a BitTorrent tracker, written in Go.

Installation

goat can be built using Go 1.1+. It can be downloaded, built, and installed,
simply by running:

	$ go get github.com/mdlayher/goat

In addition, goat depends on a MySQL server for data storage.  After creating a
database and user for goat, its database schema may be imported from the SQL
files located in 'res/'.  goat will not run unless MySQL is installed, and a
database and user are properly configured for its use.

Optionally, goat can be built to use ql (https://github.com/cznic/ql) as its storage
backend. This is done by supplying the 'ql' tag in the go get command:

	$ go get -tags='ql' github.com/mdlayher/goat

A blank ql database file is located under 'res/ql/goat.db', and will be copied to
'~/.config/goat/goat.db' on UNIX systems.  goat is now able to use ql as its
storage backend, for those who do not wish to use an external, MySQL backend.

Listeners

goat is capable of listening for torrent traffic in three modes: HTTP, HTTPS,
and UDP.  HTTP/HTTPS are the recommended methods, and are required in order for
goat to serve its API, and to allow use of private tracker passkeys.

HTTP is considered the standard mode of operation for goat.  HTTP allows gathering
a great number of metrics, use of passkeys, use of a client whitelist, and access
to goat's RESTful API, when configured.  For most trackers, this will be the only
listener which is necessary in order for goat to function properly.

The HTTPS listener provides a method to encrypt traffic to the tracker, but must
be used with caution.  Unless the SSL certificate in use is signed by a proper
certificate authority, it will distress most clients, and they may outright refuse
to announce to it.  If you are in possession of a certificate signed by a certificate
authority, this mode may be more ideal, as it provides added security for your
clients.

The UDP listener is the most unusual method of the three, and should only be used
for public trackers.  The BitTorrent UDP tracker protocol specifies a very specific
packet format, meaning that additional information or parameters cannot be packed
into a UDP datagram in a standard way.  The UDP tracker may be the fastest and least
bandwidth-intensive, but as stated, should only be used for public trackers.

API

A new feature goat added to goat in order to allow better interoperability with many
languages is a RESTful API, which is served using the HTTP or HTTPS listeners.  This
API enables easy retrieval of tracker statistics, while allowing goat to run as a
completely independent process.

It should be noted that the API is only enabled when configured, and when a HTTP or
HTTPS listener is enabled.  Without a transport mechanism, the API will be inaccessible.

The API features several modes of authentication, including HTTP Basic for login and
HMAC-SHA1 other calls.  Upon logging into the API using HTTP Basic with a username
and password pair, an API public key and secret will be generated.  The public key
is used as the username for HTTP Basic authentication, and the secret key is used
to calculate a HMAC-SHA1 signature for the password.

As part of API signature generation, a random nonce value must be generated and added
to the request.  It is added to the password portion of the HTTP Basic request, and
also to the string which is used to create the signature.  Nonce values must be changed
on every request, or the request will fail.

The current pseudocode format of the HMAC-SHA1 signature is as follows:

	signString = UserID-Nonce-HTTPMethod-HTTPResource
	ex: 1-0123abc-GET-/api/status

	signature = hmac_sha1(signString, apiSecret)

The proper format for a HTTP Basic request is as follows:

	Authorization: Basic base64(pubkey:nonce/signature)
	ex: Authorization: Basic base64(abcdef0123456789:0123abc/0123abcd4567ef89)

When the public key, nonce, and API signature are sent via HTTP Basic, the server will
verify the signature.  Successful authentication will allow access to the API.

API Calls

This list contains all API calls currently recognized by goat.  Each call must be
authenticated using the aforementioned methods.

	POST /api/login

	$ curl --user username:password http://localhost:8080/api/login
	{
		"userId": 1,
		"pubkey": "abcdef0123456789",
		"secret": "0123456789abcdef",
		"expire": 1389737644
	}

Request an API public key and secret key for this user.  The public key, user ID,
and secret key are used to authenticate further API calls.  The expire time indicates
when this key is set to expire.  Further API calls will extend the expiration time.

	GET /api/files

	$ curl --user pubkey:nonce/signature http://localhost:8080/api/files
	[
		{
			"id": 1,
			"infoHash": "abcdef0123456789",
			"verified": true,
			"createTime": 1389737644,
			"updateTime": 1389737644
		}
	]

Retrieve a list of all files tracked by goat.  Some extended attributes are not added
to reduce strain on database, and to provide a more general overview.

	GET /api/files/:id

	$ curl --user pubkey:nonce/signature http://localhost:8080/api/files/1
	{
		"id": 1,
		"infoHash": "abcdef0123456789",
		"verified": true,
		"createTime": 1389737644,
		"updateTime": 1389737644,
		"completed": 0,
		"seeders": 0,
		"leechers": 0,
		"fileUsers": [
			{
				"fileId": 1,
				"userId": 1,
				"ip": "8.8.8.8",
				"active": true,
				"completed": false,
				"announced": 1,
				"uploaded": 0,
				"downloaded": 0,
				"left": 0,
				"time": 1389983002
			}
		]
	}

Retrieve extended attributes about a specific file with matching ID.  This provides
counts for number of completions, seeders, leechers, and a list of fileUser relationships
associated with a given file.

	GET /api/status

	$ curl --user pubkey:nonce/signature http://localhost:8080/api/status
	{
		"pid": 27796,
		"hostname": "goat",
		"platform": "linux",
		"architecture": "amd64",
		"numCpu": 4,
		"numGoroutine": 14,
		"memoryMb": 1.03678,
		"maintenance": false,
		"status": "",
		"api": {
			"minute": 1
			"halfHour": 2,
			"hour": 3,
			"total": 4
		},
		"http": {
			"minute": 1
			"halfHour": 2,
			"hour": 3,
			"total": 4
		},
		"udp": {
			"minute": 1
			"halfHour": 2,
			"hour": 3,
			"total": 4
		}
	}

Retrieve a variety of metrics about the current status of goat, including its PID,
hostname, memory usage, number of HTTP/UDP hits, etc.

Configuration

goat is configured using a JSON file, which will be created under
'~/.config/goat/config.json' on UNIX systems.  Here is an example configuration,
describing the settings available to the user.

	{
		// Port: the port number on which goat will listen using both HTTP and UDP
		"Port": 8080,

		// Passkey: require that a valid passkey is present in HTTP tracker requests
		// note: this setting is typically used only for private trackers
		// ex: http://localhost:8080/0123456789ABCDEF/announce
		"Passkey": true,

		// Whitelist: require clients to be whitelisted for use with the tracker
		// note: this setting is typically used only for private trackers
		"Whitelist": true,

		// Interval: number of seconds clients should wait between announces
		"Interval": 3600,

		// HTTP: enable listening for client connections via HTTP
		"HTTP": true,

		// API: enable a HTTP RESTful API, used to pull statistics from goat
		// note: only enabled when HTTP/HTTPS is enabled
		"API": true,

		// UDP: enable listening for client connections via UDP
		// note: it is not possible to use a passkey with this listener, so this
		// listener should only be used for public trackers
		"UDP": false,

		// SSL: HTTPS configuration
		"SSL": {
			// Enabled: enable listening for client connections via HTTPS
			"Enabled": false,

			// Port: the port number on which goat will listen using HTTPS
			"Port": 8443,

			// Certificate: the path to the certificate file used for HTTPS
			"Certificate": "goat.crt",

			// Key: the path to the key file used for HTTPS
			"Key": "goat.key"
		},

		// DB: MySQL database configuration
		"DB": {
			// Host: the host and port of the MySQL database server
			"Host": "localhost:3306",

			// Database: the database goat will use to store its tracker data
			"Database": "goat",

			// Username: the username used to access goat's database
			"Username": "goat",

			// Password: the password used to access goat's database
			"Password": "goat"
		},

		// Redis: Redis cache configuration
		"Redis": {
			// Enabled: enable caching and rate limiting features using Redis
			"Enabled": false,

			// Host: the host and port of the Redis cache server
			"Host": "localhost:6379",

			// Password: optional, the password used to connect to Redis
			// note: if blank, goat will not attempt to authenticate to Redis
			"Password": ""
		}
	}

*/
package main
