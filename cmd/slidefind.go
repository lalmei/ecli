// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	"ecli/keeneye"

	"github.com/spf13/cobra"
)

// slidefindCmd represents the slidefind command
var slidefindCmd = &cobra.Command{
	Use:     "find",
	Short:   "Find a slide by criteria",
	Example: "Doo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a search term as argument.")
		}
		res, err := keeneye.Search(args[0])
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range res {
			if _, ok := e.(map[string]interface{}); ok {
				gg := e.(map[string]interface{})
				if len(gg) == 0 {
					continue
				}
				fmt.Printf("%-24s %-4s %-25s %s %s\n", "IID", "ID", "Name", "Filename", "Size (B)")
				fmt.Printf("----------------------------------------------------------------------------------\n")
				for _, val := range gg {
					gh := val.(map[string]interface{})
					fmt.Printf("%24s %-4.0f %-25q %q %.0f\n", gh["id"], gh["slideId"], gh["name"], gh["file"], gh["size"])
				}
			}
		}
	},
}

func init() {
	slideCmd.AddCommand(slidefindCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// slidefindCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// slidefindCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
