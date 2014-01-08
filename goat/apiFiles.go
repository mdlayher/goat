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
			return
		}

		// Return status
		resChan <- res
		return
	}

	// Marshal into JSON
	res, err := json.Marshal(new(FileRecordRepository).All())
	if err != nil {
		log.Println(err.Error())
		resChan <- nil
		return
	}

	resChan <- res
	return
}
