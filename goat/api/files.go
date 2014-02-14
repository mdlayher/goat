package api

import (
	"encoding/json"

	"github.com/mdlayher/goat/goat/data"
)

// getFilesJSON returns a JSON representation of one or more data.FileRecords
func getFilesJSON(ID int) ([]byte, error) {
	// Check for a valid integer ID
	if ID > 0 {
		// Load file
		file, err := new(data.FileRecord).Load(ID, "id")
		if err != nil {
			return nil, err
		}

		// Create JSON represenation
		jsonFile, err := file.ToJSON()
		if err != nil {
			return nil, err
		}

		// Marshal into JSON
		res, err := json.Marshal(jsonFile)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	// Load all files
	files, err := new(data.FileRecordRepository).All()
	if err != nil {
		return nil, err
	}

	// Marshal into JSON
	res, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}

	return res, err
}
