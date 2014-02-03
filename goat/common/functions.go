package common

import (
	crand "crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"strconv"
	"time"
)

// RandRange generates a random announce interval in the specified range
func RandRange(min int, max int) int {
	mrand.Seed(time.Now().Unix())
	return min + mrand.Intn(max-min)
}

// RandString generates a random hex string, used for generating hashes
// Thanks: http://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
func RandString() string {
	u := make([]byte, 16)
	_, err := crand.Read(u)
	if err != nil {
		// On failure, get a random number
		return strconv.Itoa(RandRange(0, 1000000))
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u)
}
