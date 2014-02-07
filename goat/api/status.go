package api

import (
	"encoding/json"

	"github.com/mdlayher/goat/goat/common"
)

// getStatusJSON returns a JSON representation of server status
func getStatusJSON() ([]byte, error) {
	// Retrieve status
	status, err := common.GetServerStatus()
	if err != nil {
		return nil, err
	}

	// Marshal into JSON from request
	res, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	// Return status
	return res, nil
}
