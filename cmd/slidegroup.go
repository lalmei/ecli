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
	"log"

	"github.com/keeneyetech/ecli/core"

	"github.com/spf13/cobra"
)

// slidegroupCmd represents the slidegroup command
var slidegroupCmd = &cobra.Command{
	Use:     "worklist",
	Aliases: []string{"w"},
	Short:   "Open default browser to view the slide item in the worklist",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := slideInfo(args)
		if err != nil {
			log.Fatal(err)
		}
		u, ok := res["url"]
		if !ok {
			log.Fatal("slide has no URL info")
		}
		v, ok := u.(map[string]interface{})["worklist"]
		if !ok {
			log.Fatal("slide has no worklist URL info")
		}
		vurl := v.(string)
		if err := core.Open(vurl); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	slideCmd.AddCommand(slidegroupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// slidegroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// slidegroupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
