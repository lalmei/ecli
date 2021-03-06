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

	"github.com/spf13/cobra"
)

// worklistCmd represents the worklist command
var worklistCmd = &cobra.Command{
	Aliases: []string{"wl"},
	Use:     "worklist",
	Short:   "Navigate and edit the work list",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.
`,
	Example: `Go to the root of the work list with
  worklist cd /
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

func init() {
	// RootCmd.AddCommand(worklistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// worklistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// worklistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	worklistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
