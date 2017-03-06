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

func SlideInfo(id uint64) (map[string]interface{}, error) {
	p := map[string]interface{}{
		"slideId": id,
	}
	v, err := sendRequest("KeenEye.SlideInfo", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}

func DeleteSlide(id uint64) error {
	p := map[string]interface{}{
		"slideIds": []uint64{id},
	}
	_, err := sendRequest("KeenEye.DeleteSlides", p)
	if err != nil {
		return err
	}
	return nil
}

func EditSlide(id, name, desc string, pixelSizeValue float64, pixelSizeUnit, imageFormat string) error {
	p := map[string]interface{}{
		"id":          id,
		"name":        name,
		"description": desc,
		"imageFormat": imageFormat,
		"pixelSize": map[string]interface{}{
			"value": pixelSizeValue,
			"unit":  pixelSizeUnit,
		},
	}
	_, err := sendRequest("KeenEye.EditSlide", p)
	if err != nil {
		return err
	}
	return nil
}
