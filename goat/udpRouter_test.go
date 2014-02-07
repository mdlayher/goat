package goat

import (
	"log"
	"net"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
	"github.com/mdlayher/goat/goat/data/udp"
)

// TestUDPRouter verifies that the main UDP router is working properly
func TestUDPRouter(t *testing.T) {
	log.Println("TestUDPRouter()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock data.FileRecord
	file := data.FileRecord{
		InfoHash: "6465616462656566303030303030303030303030",
		Verified: true,
	}

	// Save mock file
	if err := file.Save(); err != nil {
		t.Fatalf("Failed to save mock file: %s", err.Error())
	}

	// Fake UDP address
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create fake UDP address")
	}

	// Connect packet with handshake
	connect := udp.Packet{udpInitID, 0, 1234}
	connectBuf, err := connect.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to create UDP connect packet")
	}

	// Perform connection handshake
	res, err := parseUDP(connectBuf, addr)
	if err != nil {
		errRes := new(udp.ErrorResponse)
		err2 := errRes.UnmarshalBinary(res)
		if err2 != nil {
			t.Fatalf(err.Error())
		}

		log.Println("ERROR:", errRes.Error)
		t.Fatalf(err.Error())
	}

	// Retrieve response, get new connection ID, which will be expected by router
	connRes := new(udp.ConnectResponse)
	err = connRes.UnmarshalBinary(res)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Create announce request
	announce := udp.AnnounceRequest{
		ConnID:     connRes.ConnID,
		Action:     1,
		TransID:    connRes.TransID,
		InfoHash:   []byte("deadbeef000000000000"),
		PeerID:     []byte("00001111222233334444"),
		Downloaded: 0,
		Left:       0,
		Uploaded:   0,
		IP:         0,
		Key:        1234,
		Port:       5000,
	}

	// Get announce bytes
	announceBuf, err := announce.MarshalBinary()
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Send announce to UDP router
	res, err = parseUDP(announceBuf, addr)
	if err != nil {
		errRes := new(udp.ErrorResponse)
		err2 := errRes.UnmarshalBinary(res)
		if err2 != nil {
			t.Fatalf(err.Error())
		}

		log.Println("ERROR:", errRes.Error)
		t.Fatalf(err.Error())
	}

	// Get UDP announce response
	announceRes := new(udp.AnnounceResponse)
	err = announceRes.UnmarshalBinary(res)
	if err != nil {
		errRes := new(udp.ErrorResponse)
		err2 := errRes.UnmarshalBinary(res)
		if err2 != nil {
			t.Fatalf(err.Error())
		}

		log.Println("ERROR:", errRes.Error)
		t.Fatalf(err.Error())
	}
	log.Println(announceRes)

	// Create scrape request
	scrape := udp.ScrapeRequest{
		ConnID:     connRes.ConnID,
		Action:     2,
		TransID:    connRes.TransID,
		InfoHashes: [][]byte{[]byte("deadbeef000000000000")},
	}

	// Get scrape bytes
	scrapeBuf, err := scrape.MarshalBinary()
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Send scrape to UDP router
	res, err = parseUDP(scrapeBuf, addr)
	if err != nil {
		errRes := new(udp.ErrorResponse)
		err2 := errRes.UnmarshalBinary(res)
		if err2 != nil {
			t.Fatalf(err.Error())
		}

		log.Println("ERROR:", errRes.Error)
		t.Fatalf(err.Error())
	}

	// Get UDP scrape response
	scrapeRes := new(udp.ScrapeResponse)
	err = scrapeRes.UnmarshalBinary(res)
	if err != nil {
		errRes := new(udp.ErrorResponse)
		err2 := errRes.UnmarshalBinary(res)
		if err2 != nil {
			t.Fatalf(err.Error())
		}

		log.Println("ERROR:", errRes.Error)
		t.Fatalf(err.Error())
	}
	log.Println(scrapeRes)

	// Delete mock file
	if err := file.Delete(); err != nil {
		t.Fatalf("Failed to delete mock file: %s", err.Error())
	}
}
