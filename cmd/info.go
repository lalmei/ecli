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
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/keeneyetech/ecli/api"
	"github.com/spf13/cobra"
)

var cfgInfoLast bool

func slideInfo(args []string) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("Please provide a slide ID (integer)")
	}
	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Slide ID must be an integer, not %q", args[0])
	}
	res, err := api.SlideInfo(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "info SLIDE_ID",
	Aliases: []string{"i"},
	Short:   "Get slide information",
	Run: func(cmd *cobra.Command, args []string) {
		if cfgInfoLast {
			res, err := api.Search(".") // any
			if err != nil {
				log.Fatal(err)
			}
			for key, reply := range res {
				if key != "items" {
					continue
				}
				slides, ok := reply.(map[string]interface{})["slides"]
				if ok {
					lslides := slides.([]interface{})
					if len(lslides) > 0 {
						s := lslides[len(lslides)-1]
						slideId := s.(map[string]interface{})["slideId"].(float64)
						// Rewrite args for later call to slideInfo()
						args = []string{fmt.Sprintf("%.0f", slideId)}
						break
					} else {
						log.Fatal("No slide uploaded so far, so no info to display.")
					}
				}
			}
		}

		res, err := slideInfo(args)
		if err != nil {
			log.Fatal(err)
		}
		d, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(string(d))
	},
}

func init() {
	slideCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	infoCmd.Flags().BoolVarP(&cfgInfoLast, "last", "l", false, "Show info about last uploaded image")
}
