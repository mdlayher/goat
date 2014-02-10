package udp

import (
	"log"
	"testing"
)

// TestPacket verifies that Packet binary marshal and unmarshal work properly
func TestPacket(t *testing.T) {
	log.Println("TestPacket()")

	// Generate mock Packet
	packet := Packet{
		ConnID: 1,
		Action: 0,
		TransID: 1,
	}

	// Marshal to binary representation
	out, err := packet.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal Packet to binary: %s", err.Error())
	}

	// Unmarshal packet from binary representation
	packet2 := new(Packet)
	if err := packet2.UnmarshalBinary(out); err != nil {
		t.Fatalf("Failed to unmarshal Packet from binary: %s", err.Error())
	}

	// Verify packets are identical
	if packet.ConnID != packet2.ConnID || packet.Action != packet2.Action {
		t.Fatalf("Packet results do not match")
	}
}
