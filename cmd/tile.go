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
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/keeneyetech/ecli/api"
	"github.com/spf13/cobra"
)

// tileCmd represents the info command
var tileCmd = &cobra.Command{
	Use:   "tile SLIDE_ID",
	Short: "Explicitely run the tiling on a slide",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a slide ID (integer)")
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("Slide ID must be an integer, not %q", args[0])
		}
		if err := api.Tile(id); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Tiling slide %d in the background...\n", id)
	},
}

func init() {
	slideCmd.AddCommand(tileCmd)
}
