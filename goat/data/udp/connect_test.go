package udp

import (
	"log"
	"testing"
)

// TestConnectResponse verifies that ConnectResponse binary marshal and unmarshal work properly
func TestConnectResponse(t *testing.T) {
	log.Println("TestConnectResponse()")

	// Generate mock ConnectResponse
	connect := ConnectResponse{
		ConnID: 1,
		Action: 0,
		TransID: 1,
	}

	// Marshal to binary representation
	out, err := connect.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal ConnectResponse to binary: %s", err.Error())
	}

	// Unmarshal connect from binary representation
	connect2 := new(ConnectResponse)
	if err := connect2.UnmarshalBinary(out); err != nil {
		t.Fatalf("Failed to unmarshal ConnectResponse from binary: %s", err.Error())
	}

	// Verify connects are identical
	if connect.ConnID != connect2.ConnID || connect.TransID != connect2.TransID {
		t.Fatalf("ConnectResponse results do not match")
	}
}
