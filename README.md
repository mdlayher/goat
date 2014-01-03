goat
====

goat: __Go__ __A__wesome __T__racker.  BitTorrent tracker implementation, written in Go.  MIT Licensed.

Configuration
-------------

goat is configured using a JSON file, which will be created under `~/.config/goat/config.json`.  Here is an example
configuration, describing the settings available to the user.

```
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

	// Http: enable listening for client connections via HTTP
	"Http": true,

	// Udp: enable listening for client connections via UDP
	// note: it is not possible to use a passkey with this listener, so this listener should
	// only be used for public trackers
	"Udp": false
}
```
