package udp

import (
	"log"
	"testing"
)

// TestErrorResponse verifies that ErrorResponse binary marshal and unmarshal work properly
func TestErrorResponse(t *testing.T) {
	log.Println("TestErrorResponse()")

	// Generate mock ErrorResponse
	errRes := ErrorResponse{
		Action: 3,
		TransID: 1,
		Error: "Testing",
	}

	// Marshal to binary representation
	out, err := errRes.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal ErrorResponse to binary: %s", err.Error())
	}

	// Unmarshal errRes from binary representation
	errRes2 := new(ErrorResponse)
	if err := errRes2.UnmarshalBinary(out); err != nil {
		t.Fatalf("Failed to unmarshal ErrorResponse from binary: %s", err.Error())
	}

	// Verify errRess are identical
	if errRes.Action != errRes2.Action || errRes.Error != errRes2.Error {
		t.Fatalf("ErrorResponse results do not match")
	}
}
