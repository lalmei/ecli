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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/keeneyetech/ecli/api"
	"github.com/spf13/cobra"
)

var cfgApply bool

// grouprmCmd represents the grouprm command
var grouprmCmd = &cobra.Command{
	Use:   "rm GROUP_ID",
	Short: "Delete a group",
	Long: `Deleting a group recursively deletes all the data it! As a security,
You must pass the -y option to really perform the delete action.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing group ID.")
		}
		if !cfgApply {
			fmt.Println("WARNING: removing a group will recursively remove all data inside it!" +
				" If you are doing it on purpose, rerun the command with the -y flag.")
			os.Exit(2)
		}
		if err := api.DeleteGroup(args[0]); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Group %q (and all its inner data) deleted.\n", args[0])
	},
}

func init() {
	groupCmd.AddCommand(grouprmCmd)

	grouprmCmd.Flags().BoolVarP(&cfgApply, "apply", "y", false, "Really run the command")
}
