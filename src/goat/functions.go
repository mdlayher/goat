package goat

import (
	"math/rand"
	"time"
)

// Generate a random announce interval in the specified range
func RandRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}
