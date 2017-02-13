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

	log "github.com/Sirupsen/logrus"
	"github.com/keeneyetech/ecli/api"
	"github.com/spf13/cobra"
)

// groupupdateCmd represents the groupupdate command
var groupupdateCmd = &cobra.Command{
	Use:     "update GROUP_ID",
	Aliases: []string{"up"},
	Short:   "Update information of existing group",
	Long: `A new name, description or set of labels can be provided to edit a group. 

Change the description of a group:

  group update 58a1a5b4e7798928257123a0 --desc "A better description"

You can pass several --label flags to add several existing labels:

  group update 58a1a5b4e7798928257123a0 --name "New name" --label caution --label eye

Passing no --label at all will make the labels unchanged for a group, if any. Passing
--label "" will remove all labels apply on a group:

  group update 58a1a5b4e7798928257123a0 --label ""`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing group ID.")
		}
		var labels []*api.Label
		if len(cfgGroupLabels) > 0 {
			// -l "" will remove all group labels
			if len(cfgGroupLabels) >= 1 && cfgGroupLabels[0] != "" {
				var err error
				labels, err = checkLabels(cfgGroupLabels)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			// No label specified on the CLI so we must pass all existing group's labels
			// otherwise they will get deleted on update.
			// FIXME: this can be a slow process; will be faster when the API gets a new
			// `Groups()` endpoint.
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
						if p["id"].(string) != args[0] {
							continue
						}
						for _, lab := range p["labels"].([]interface{}) {
							c := lab.(map[string]interface{})
							labels = append(labels, &api.Label{
								c["name"].(string),
								c["color"].(string),
								c["description"].(string),
							})
						}
					}
				}
			}
		}
		if err := api.EditGroup(args[0], cfgGroupName, cfgGroupDesc, labels); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Group succcessfully updated.\n")
	},
}

func init() {
	groupCmd.AddCommand(groupupdateCmd)

	groupupdateCmd.Flags().StringVar(&cfgGroupName, "name", "", "Name")
	groupupdateCmd.Flags().StringVar(&cfgGroupDesc, "desc", "", "Short description")
	groupupdateCmd.Flags().StringArrayVarP(&cfgGroupLabels, "label", "l", []string{}, "Label to apply on group")
}
