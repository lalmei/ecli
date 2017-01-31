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
	"strconv"

	"github.com/keeneyetech/ecli/api"

	"github.com/spf13/cobra"
)

// rmslideCmd represents the rmslide command
var rmslideCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a slide",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a slide ID, either as int or hex string.")
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Slide ID must be an integer, not %q", args[0])
		}
		if err := api.DeleteSlide(id); err != nil {
			log.Fatalf("can't delete slide: %s", err.Error())
		}
	},
}

func init() {
	slideCmd.AddCommand(rmslideCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmslideCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmslideCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
