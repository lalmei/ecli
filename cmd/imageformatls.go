// Copyright © 2016 Mathias Monnerville <mathias.monnerville@keeneye.tech>

package cmd

import (
	"fmt"
	"log"

	"ecli/api"

	"github.com/spf13/cobra"
)

// imageformatlsCmd represents the imageformatls command
var imageformatlsCmd = &cobra.Command{
	Use:     "imageformats",
	Aliases: []string{"if"},
	Short:   "List all supported image formats",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := api.ImageFormats()
		if err != nil {
			log.Fatal(err)
		}
		its := res["imageFormats"]
		if its != nil {
			for _, it := range its.([]interface{}) {
				p := it.(map[string]interface{})
				if !cfgQuiet {
					fmt.Printf("%-30s %-30s %-30s", p["id"], p["name"], p["shortDescription"])
					if p["icon"] != "" {
						fmt.Printf(" [icon: yes]")
					}
				} else {
					fmt.Print(p["id"])
				}
				fmt.Println()
			}
		} else {
			log.Fatal("The response has an unexpected format and cannot be displayed.")
		}
	},
}

func init() {
	RootCmd.AddCommand(imageformatlsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageformatlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageformatlsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
