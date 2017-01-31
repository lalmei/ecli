// Copyright © 2017 The Keen Eye Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/keeneyetech/ecli/config"

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
