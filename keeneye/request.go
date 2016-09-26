package keeneye

import (
	"bytes"
	"net/http"

	json "github.com/gorilla/rpc/v2/json2"
)

const (
	serverAddress = "https://dev.keeneyetechnologies.com/api/v2"
)

func sendRequest(method string, args interface{}) (interface{}, error) {
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
		return nil, err
	}
	return result, nil
}
