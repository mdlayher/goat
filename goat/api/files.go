package api

import (
	"encoding/json"
	"log"

	"github.com/mdlayher/goat/goat/data"
)

// getFilesJSON returns a JSON representation of one or more data.FileRecords
func getFilesJSON(ID int) []byte {
	// Check for a valid integer ID
	if ID > 0 {
		// Load file
		return new(data.FileRecord).Load(ID, "id").ToJSON()
	}

	// Marshal into JSON
	res, err := json.Marshal(new(data.FileRecordRepository).All())
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return res
}
