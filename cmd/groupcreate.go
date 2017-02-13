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

var (
	cfgGroupName   string
	cfgGroupDesc   string
	cfgGroupLabels []string
)

// groupcreateCmd represents the groupcreate command
var groupcreateCmd = &cobra.Command{
	Use:   "create NAME",
	Short: "Create a group",
	Long: `The group will be attached to the root element by default. A parent group
ID can also be specified by using --group-id. Example:

  group create "My inner group" --group-id 5881df5ae77989696a8b1702

A description can be provided with --desc:

  group create "My inner group" --group-id 5881df5ae77989696a8b1702 --desc "New study"

Also, existing labels can be added a group using --label:

  group create "My inner group" --group-id 5881df5ae77989696a8b1702 --label urgent --label eye`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing group name.")
		}
		name := args[0]
		customLabels, err := checkLabels(cfgGroupLabels)
		if err != nil {
			log.Fatal(err)
		}
		if err := api.NewGroup(name, cfgLabelDesc, customLabels, cfgParentId); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Group %q created.\n", name)
	},
}

func init() {
	groupCmd.AddCommand(groupcreateCmd)

	groupcreateCmd.Flags().StringVar(&cfgGroupDesc, "desc", "", "Short description")
	groupcreateCmd.Flags().StringArrayVarP(&cfgGroupLabels, "label", "l", []string{}, "Label to apply on group")
	groupcreateCmd.Flags().StringVarP(&cfgParentId, "group-id", "g", "root", "Group parent ID")
}
