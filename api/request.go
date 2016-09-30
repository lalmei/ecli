package api

import (
	"bytes"
	"fmt"
	"net/http"

	"ecli/config"

	json "github.com/gorilla/rpc/v2/json2"
)

var serverAddress string

func sendRequest(method string, args interface{}) (interface{}, error) {
	serverAddress, tok, err := config.LoadSession()
	if err != nil {
		return nil, err
	}
	if serverAddress == "" {
		return nil, fmt.Errorf("empty server url, cannot access endpoint")
	}
	if args != nil {
		args.(map[string]interface{})["token"] = tok
	}

	message, err := json.EncodeClientRequest(method, args)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(message))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result interface{}
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		if err.Error() == "TokenExpired" {
			if err := config.DeleteTokenFile(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("Your session expired. Please run \"login\" command again.")
		}
		return nil, err
	}
	return result, nil
}
