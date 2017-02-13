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

package cmd

import (
	"fmt"
	"log"

	"github.com/keeneyetech/ecli/api"
	"github.com/spf13/cobra"
)

// grouplsCmd represents the groupls command
var grouplsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all existing groups",
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME: inefficient, but API v2.10 is currently missing a dedicated `Groups()` endpoint.
		res, err := api.Search(".")
		if err != nil {
			log.Fatal(err)
		}
		for key, reply := range res {
			if key != "items" {
				continue
			}
			if groups, ok := reply.(map[string]interface{})["groups"]; ok {
				for _, s := range groups.([]interface{}) {
					p := s.(map[string]interface{})
					if !cfgQuiet {
						fmt.Printf("%24s %-25q\n", p["id"], p["name"])
					} else {
						fmt.Println(p["name"])
					}
				}
			}
		}
	},
}

func init() {
	groupCmd.AddCommand(grouplsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grouplsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grouplsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
