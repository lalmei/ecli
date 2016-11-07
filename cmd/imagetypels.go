// Copyright © 2016 Mathias Monnerville <mathias.monnerville@keeneye.tech>

package cmd

import (
	"fmt"
	"log"

	"ecli/api"

	"github.com/spf13/cobra"
)

// imagetypelsCmd represents the imagetypels command
var imagetypelsCmd = &cobra.Command{
	Use:     "imagetypes",
	Aliases: []string{"it"},
	Short:   "List all supported image types",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := api.ImageTypes()
		if err != nil {
			log.Fatal(err)
		}
		its := res["imageTypes"]
		if its != nil {
			for _, it := range its.([]interface{}) {
				p := it.(map[string]interface{})
				fmt.Printf("%-30s %-30s %-30s", p["id"], p["name"], p["shortDescription"])
				if p["icon"] != "" {
					fmt.Printf(" [icon: yes]")
				}
				fmt.Println()
			}
		} else {
			log.Fatal("The response has an unexpected format and cannot be displayed.")
		}
	},
}

func init() {
	RootCmd.AddCommand(imagetypelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagetypelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagetypelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
