package api

import (
	"encoding/json"

	"github.com/mdlayher/goat/goat/data"
)

// postLogin generates a new API key for this user
func postLogin(session data.UserRecord) ([]byte, error) {
	// Create key for this user's session
	key := new(data.APIKey)
	if err := key.Create(session.ID); err != nil {
		return nil, err
	}

	// Store key in database
	if err := key.Save(); err != nil {
		return nil, err
	}

	// Convert key to JSON form
	jsonKey, err := key.ToJSON()
	if err != nil {
		return nil, err
	}

	// Marshal into JSON
	res, err := json.Marshal(jsonKey)
	if err != nil {
		return nil, err
	}

	return res, nil
}
