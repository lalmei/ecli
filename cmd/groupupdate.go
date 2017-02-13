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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing group ID.")
		}
		if err := api.EditGroup(args[0], cfgGroupName, cfgGroupDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Group succcessfully updated.\n")
	},
}

func init() {
	groupCmd.AddCommand(groupupdateCmd)

	groupupdateCmd.Flags().StringVar(&cfgGroupName, "name", "", "Name")
	groupupdateCmd.Flags().StringVar(&cfgGroupDesc, "desc", "", "Short description")
}
