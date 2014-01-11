package goat

import (
	"errors"
	"log"

	"github.com/garyburd/redigo/redis"
)

// redisConnect initiates a connection to Redis server
func redisConnect() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

// redisPing verifies that Redis server is available
func redisPing() bool {
	// Send redis a PING request
	reply, err := redisDo("PING")
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Ensure value is valid
	res, err := redis.String(reply, nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// PONG is valid response
	if res != "PONG" {
		log.Println("redisPing: redis replied to PING with:", res)
		return false
	}

	// Redis OK
	return true
}

// redisDo runs a single Redis command and returns its reply
func redisDo(command string, args ...interface{}) (interface{}, error) {
	// Open Redis connection
	c, err := redisConnect()
	defer c.Close()
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("redisDo: failed to connect to redis")
	}

	// Send Redis command with arguments, receive reply
	reply, err := c.Do(command, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("redisDo: failed to send command to redis: " + command)
	}

	return reply, nil
}
