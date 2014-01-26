package goat

import (
	"errors"
	"log"

	"github.com/garyburd/redigo/redis"
)

var RedisPass *string

// redisConnect initiates a connection to Redis server
func redisConnect() (c redis.Conn, err error) {
	c, err = redis.Dial("tcp", "crestfish.redistogo.com:11107")
	if err != nil {
		return
	}

	// Authenticate with Redis database if necessary
	if RedisPass != nil && *RedisPass != "" {
		_, err = c.Do("AUTH", RedisPass)
	}
	return
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

	if err := c.Close(); err != nil {
		log.Println(err.Error())
	}

	return reply, nil
}
