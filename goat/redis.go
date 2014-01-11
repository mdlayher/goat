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
	if err != nil {
		log.Println(err.Error())
		return false
	}

	defer c.Close()
	return true
}
