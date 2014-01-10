package goat

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// Connect to Redis server
func redisConnect() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

// Verify that Redis server is available
func redisPing() bool {
	c, err := redisConnect()
	defer c.Close()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Send Redis a PING to verify it is working
	err = c.Send("PING")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	c.Flush()

	// Receive a PONG in return
	val, err := c.Receive()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Ensure value is valid
	res, err := redis.String(val, nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// PONG is valid response
	if res != "PONG" {
		log.Println("error: redis replied to PING with:", res)
		return false
	}

	// Redis OK
	return true
}
