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

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	tokenFile = ".token"
)

func absTokenFile() string {
	return path.Join(os.Getenv("HOME"), tokenFile)
}

func StoreSession(url, token string) error {
	f, err := os.Create(absTokenFile())
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	b := &struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	}{url, token}
	return enc.Encode(b)
}

// Loads token and url endpoint information.
func LoadSession() (string, string, error) {
	f, err := os.Open(absTokenFile())
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("You are not logged in. Please run \"login\" command first.")
		}
		return "", "", err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	b := &struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	}{}
	if err := dec.Decode(b); err != nil {
		return "", "", err
	}
	return b.Url, b.Token, nil
}

func DeleteTokenFile() error {
	return os.Remove(absTokenFile())
}
