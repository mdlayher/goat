package goat

import (
	"encoding/json"
	"log"
)

// GetFilesJSON returns a JSON representation of one or more FileRecords
func GetFilesJSON(ID int, resChan chan []byte) {
	// Check for a valid integer ID
	if ID > 0 {
		// Load file
		file := new(FileRecord).Load(ID, "id")

		// Marshal into JSON
		res, err := json.Marshal(file)
		if err != nil {
			log.Println(err.Error())
			resChan <- nil
		}

		// Return status
		resChan <- res
	}

	resChan <- nil
}
