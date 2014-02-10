package udp

import (
	"bytes"
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/data"
)

// TestAnnounceRequest verifies that AnnounceRequest binary marshal and unmarshal work properly
func TestAnnounceRequest(t *testing.T) {
	log.Println("TestAnnounceRequest()")

	// Generate mock AnnounceRequest
	announce := AnnounceRequest{
		ConnID:     1,
		Action:     1,
		TransID:    1,
		InfoHash:   []byte("abcdef0123456789abcd"),
		PeerID:     []byte("abcdef0123456789abcd"),
		Downloaded: 0,
		Left:       0,
		Uploaded:   0,
		Event:      0,
		IP:         1234,
		Key:        1234,
		Numwant:    50,
		Port:       8080,
	}

	// Marshal to binary representation
	out, err := announce.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal AnnounceRequest to binary: %s", err.Error())
	}

	// Unmarshal announce from binary representation
	announce2 := new(AnnounceRequest)
	if err := announce2.UnmarshalBinary(out); err != nil {
		t.Fatalf("Failed to unmarshal AnnounceRequest from binary: %s", err.Error())
	}

	// Verify announces are identical
	if announce.ConnID != announce2.ConnID || !bytes.Equal(announce.InfoHash, announce2.InfoHash) {
		t.Fatalf("AnnounceRequest results do not match")
	}

	// Convert announce to url.Values map
	query := announce.ToValues()

	// Verify conversion occurred properly
	if query.Get("downloaded") != "0" || query.Get("numwant") != "50" {
		t.Fatalf("AnnounceRequest values map results are not correct")
	}
}

// TestAnnounceResponse verifies that AnnounceResponse binary marshal and unmarshal work properly
func TestAnnounceResponse(t *testing.T) {
	log.Println("TestAnnounceResponse()")

	// Generate mock AnnounceResponse
	announce := AnnounceResponse{
		Action:   1,
		TransID:  1,
		Interval: 3600,
		Leechers: 1,
		Seeders:  1,
		PeerList: []data.Peer{data.Peer{"127.0.0.1", 8080}, data.Peer{"192.168.1.1", 4040}},
	}

	// Marshal to binary representation
	out, err := announce.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal AnnounceResponse to binary: %s", err.Error())
	}

	// Unmarshal announce from binary representation
	announce2 := new(AnnounceResponse)
	if err := announce2.UnmarshalBinary(out); err != nil {
		t.Fatalf("Failed to unmarshal AnnounceResponse from binary: %s", err.Error())
	}

	// Verify announces are identical
	if announce.Action != announce2.Action || announce.PeerList[0].IP != announce2.PeerList[0].IP {
		t.Fatalf("AnnounceResponse results do not match")
	}
}
