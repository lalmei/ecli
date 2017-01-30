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

// OpenSession returns a token upon successul authentication.
func OpenSession(url, login, password string) (string, error) {
	serverAddress = url
	p := map[string]interface{}{
		"login":    login,
		"password": password,
	}
	v, err := sendRequest("KeenEye.OpenSession", p)
	if err != nil {
		return "", err
	}
	return v.(map[string]interface{})["token"].(string), err
}

func CloseSession() error {
	p := make(map[string]interface{})
	_, err := sendRequest("KeenEye.CloseSession", p)
	if err != nil {
		return err
	}
	return nil
}
