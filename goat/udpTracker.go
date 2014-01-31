package goat

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/url"
	"strconv"
)

// udpTracker generates responses in the UDP datagram format
type udpTracker struct {
	TransID uint32
}

// Announce announces using UDP format
func (u udpTracker) Announce(query url.Values, file fileRecord) []byte {
	// Response buffer
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (1 for announce)
	err := binary.Write(res, binary.BigEndian, uint32(1))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	// Transaction ID
	err = binary.Write(res, binary.BigEndian, u.TransID)
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	// Interval
	err = binary.Write(res, binary.BigEndian, uint32(static.Config.Interval))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	// Leechers
	err = binary.Write(res, binary.BigEndian, uint32(file.Leechers()))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	// Seeders
	err = binary.Write(res, binary.BigEndian, uint32(file.Seeders()))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	// Peer list
	numwant, err := strconv.Atoi(query.Get("numwant"))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	err = binary.Write(res, binary.BigEndian, file.PeerList(query.Get("ip"), numwant))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	return res.Bytes()
}

// Error reports a UDP []byte response packed datagram
func (u udpTracker) Error(msg string) []byte {
	// Response buffer
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (3 for error)
	err := binary.Write(res, binary.BigEndian, uint32(3))
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Transaction ID
	err = binary.Write(res, binary.BigEndian, u.TransID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Error message
	err = binary.Write(res, binary.BigEndian, []byte(msg))
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return res.Bytes()
}

// Protocol returns the protocol used by this tracker
func (u udpTracker) Protocol() string {
	return "UDP"
}

// Scrape scrapes using UDP format
func (u udpTracker) Scrape(files []fileRecord) []byte {
	// Response buffer
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (2 for scrape)
	err := binary.Write(res, binary.BigEndian, uint32(2))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP scrape response")
	}

	// Transaction ID
	err = binary.Write(res, binary.BigEndian, u.TransID)
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP scrape response")
	}

	// Iterate all files, writing their statistics into the buffer
	for _, file := range files {
		// Seeders
		err = binary.Write(res, binary.BigEndian, uint32(file.Seeders()))
		if err != nil {
			log.Println(err.Error())
			return u.Error("Could not create UDP scrape response")
		}

		// Completed
		err = binary.Write(res, binary.BigEndian, uint32(file.Completed()))
		if err != nil {
			log.Println(err.Error())
			return u.Error("Could not create UDP scrape response")
		}

		// Leechers
		err = binary.Write(res, binary.BigEndian, uint32(file.Leechers()))
		if err != nil {
			log.Println(err.Error())
			return u.Error("Could not create UDP scrape response")
		}
	}

	return res.Bytes()
}
