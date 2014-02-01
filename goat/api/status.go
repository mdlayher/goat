package api

import (
	"encoding/json"
	"log"

	"github.com/mdlayher/goat/goat/common"
)

// getStatusJSON returns a JSON representation of server status
func getStatusJSON() []byte {
	// Marshal into JSON from request
	res, err := json.Marshal(common.GetServerStatus())
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Return status
	return res
}
