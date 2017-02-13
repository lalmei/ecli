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
	cfgLabelName  string
	cfgLabelColor string
	cfgLabelDesc  string
)

const dfltBackgroundColor = "#FF0000"

// labelcreateCmd represents the labelcreate command
var labelcreateCmd = &cobra.Command{
	Use:   "create NAME",
	Short: "Create a label",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing label name.")
		}
		name := args[0]
		color := cfgLabelColor
		if color == "" {
			color = dfltBackgroundColor
		}
		if err := api.NewLabel(name, color, cfgLabelDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Label %q created.\n", name)
	},
}

func init() {
	labelCmd.AddCommand(labelcreateCmd)

	labelcreateCmd.Flags().StringVar(&cfgLabelColor, "color", "", "Background color")
	labelcreateCmd.Flags().StringVar(&cfgLabelDesc, "desc", "", "Short description")
}
