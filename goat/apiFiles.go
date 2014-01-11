package goat

import (
	"encoding/json"
	"log"
)

// GetFilesJSON returns a JSON representation of one or more fileRecords
func getFilesJSON(ID int, resChan chan []byte) {
	// Check for a valid integer ID
	if ID > 0 {
		// Load file
		resChan <- new(fileRecord).Load(ID, "id").ToJSON()
		return
	}

	// Marshal into JSON
	res, err := json.Marshal(new(fileRecordRepository).All())
	if err != nil {
		log.Println(err.Error())
		resChan <- nil
		return
	}

	resChan <- res
	return
}
