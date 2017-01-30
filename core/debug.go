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

package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func DebugResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("An error occurred. Here is some debug info:\n\n")
		fmt.Println(">> RESPONSE HEADER")
		for k, v := range resp.Header {
			fmt.Printf("%s: %q\n", k, v[0])
		}
		fmt.Printf("\n>> RESPONSE BODY\n")
		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n\n>> HTTP ERROR\n", string(d))
		return fmt.Errorf(resp.Status)
	}
	return nil
}
