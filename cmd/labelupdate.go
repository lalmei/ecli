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

// labelupdateCmd represents the labelupdate command
var labelupdateCmd = &cobra.Command{
	Use:     "update NAME",
	Aliases: []string{"up"},
	Short:   "Update information on existing label",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing label name.")
		}
		oldName := args[0]
		z, err := api.Label(oldName)
		if err != nil {
			log.Fatal(err)
		}
		// Need to build a label with non-empty values so older values don't
		// get overwritten by empty strings.
		oldColor := z["color"].(string)
		oldDesc := z["description"].(string)
		name := cfgLabelName
		if name == "" {
			name = oldName
		}
		color := cfgLabelColor
		if color == "" {
			color = oldColor
		}
		desc := cfgLabelDesc
		if desc == "" {
			desc = oldDesc
		}

		if err := api.EditLabel(oldName, name, color, desc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Label %q succcessfully updated.\n", name)
	},
}

func init() {
	labelCmd.AddCommand(labelupdateCmd)

	labelupdateCmd.Flags().StringVar(&cfgLabelColor, "color", "", "Background color")
	labelupdateCmd.Flags().StringVar(&cfgLabelDesc, "desc", "", "Short description")
	labelupdateCmd.Flags().StringVar(&cfgLabelName, "name", "", "Name")
}
