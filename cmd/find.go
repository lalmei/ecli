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

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:     "find",
	Short:   "Find slides or groups by criteria",
	Example: `ecli find "test slide"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showHelpAndExit(cmd, "Please provide a search term as argument.")
		}
		res, err := api.Search(args[0])
		if err != nil {
			log.Fatal(err)
		}
		for key, reply := range res {
			if key != "items" {
				continue
			}
			groups := reply.(map[string]interface{})["groups"]
			if groups != nil {
				for _, s := range groups.([]interface{}) {
					p := s.(map[string]interface{})
					if !cfgQuiet {
						fmt.Printf("%-5s %24s %-25q\n", "group", p["id"], p["name"])
					} else {
						fmt.Printf("g %s\n", p["name"])
					}
				}
			} else {
				log.Fatal("The response has an unexpected format and cannot be displayed.")
			}
			slides := reply.(map[string]interface{})["slides"]
			if slides != nil {
				for _, s := range slides.([]interface{}) {
					p := s.(map[string]interface{})
					if !cfgQuiet {
						fmt.Printf("%-5s %24s %-.0f %q (%q, %.0f bytes)\n", "slide", p["id"], p["slideId"], p["name"], p["file"], p["size"])
					} else {
						fmt.Printf("s %-.0f %q\n", p["slideId"], p["name"])
					}
				}
			} else {
				log.Fatal("The response has an unexpected format and cannot be displayed.")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
