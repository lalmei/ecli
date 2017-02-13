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

func NewGroup(name, desc string, labels []*Label, parentId string) error {
	p := map[string]interface{}{
		"name":        name,
		"description": desc,
		"labels":      labels,
	}
	if parentId != "root" {
		p["parentId"] = parentId
	}
	_, err := sendRequest("KeenEye.NewGroup", p)
	if err != nil {
		return err
	}
	return nil
}

// FIXME: will take a labels parameter once the API has a `Groups()` endpoint
// to fetch the current labels applied on a group. So for now, no label editing
// when updating a group.
func EditGroup(id, name, desc string) error {
	p := map[string]interface{}{
		"id":          id,
		"name":        name,
		"description": desc,
	}
	_, err := sendRequest("KeenEye.EditGroup", p)
	if err != nil {
		return err
	}
	return nil
}

func DeleteGroup(id string) error {
	p := map[string]interface{}{
		"id": id,
	}
	_, err := sendRequest("KeenEye.DeleteGroup", p)
	if err != nil {
		return err
	}
	return nil
}

func Groups() ([]interface{}, error) {
	// FIXME: add implementation when API has a `Groups()` endpoint.
	return nil, fmt.Errorf("Not implemented")
}
