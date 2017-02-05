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

import "fmt"

func NewLabel(name, color, desc string) error {
	p := map[string]interface{}{
		"name":        name,
		"color":       color,
		"description": desc,
	}
	_, err := sendRequest("KeenEye.NewLabel", p)
	if err != nil {
		return err
	}
	return nil
}

func EditLabel(oldName, name, color, desc string) error {
	p := map[string]interface{}{
		"oldName":     oldName,
		"name":        name,
		"color":       color,
		"description": desc,
	}
	_, err := sendRequest("KeenEye.EditLabel", p)
	if err != nil {
		return err
	}
	return nil
}

func DeleteLabel(name string) error {
	p := map[string]interface{}{
		"name": name,
	}
	_, err := sendRequest("KeenEye.DeleteLabel", p)
	if err != nil {
		return err
	}
	return nil
}

func Label(name string) (map[string]interface{}, error) {
	p := map[string]interface{}{
		"name": name,
	}
	v, err := sendRequest("KeenEye.Label", p)
	if err != nil {
		return nil, err
	}
	return v.(map[string]interface{}), nil
}

func Labels() ([]interface{}, error) {
	p := map[string]interface{}{}
	v, err := sendRequest("KeenEye.Labels", p)
	if err != nil {
		return nil, err
	}
	labels, ok := v.(map[string]interface{})["labels"]
	if !ok {
		return nil, fmt.Errorf("labels: wrong reply format")
	}
	return labels.([]interface{}), nil
}
