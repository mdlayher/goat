package goat

import (
	"math/rand"
	"time"
)

// RandRange generates a random announce interval in the specified range
func randRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}
