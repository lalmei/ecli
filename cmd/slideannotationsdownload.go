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
	"log"
	"os"
	"strconv"

	"github.com/keeneyetech/ecli/api"
	"github.com/spf13/cobra"
)

var cfgOutputFile string

// slideannotationsdownloadCmd represents the slideannotationsdownload command
var slideannotationsdownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download slide annotations and write them to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a slide ID (integer)")
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Slide ID must be an integer, not %q", args[0])
		}
		res, err := api.DownloadAnnotations(id)
		if err != nil {
			log.Fatal(err)
		}
		path := res["path"].(string)
		if path == "" {
			fmt.Println("The annotations file is being built. This is an asynchronous process. Please rerun that command in a few seconds to get the annotations.")
			os.Exit(3)
		}
		var out = os.Stdout
		if cfgOutputFile != "" {
			f, err := os.Create(cfgOutputFile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			out = f
		}
		if err := api.DownloadRessource(path, out); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	slideannotationsCmd.AddCommand(slideannotationsdownloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// slideannotationsdownloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// slideannotationsdownloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	slideannotationsdownloadCmd.Flags().StringVarP(&cfgOutputFile, "output", "o", "", "Write annotations to this file")
}
