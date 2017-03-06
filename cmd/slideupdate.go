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

// slideupdateCmd represents the slideupdate command
var slideupdateCmd = &cobra.Command{
	Use:   "update SLIDE_ID",
	Short: "Update information of existing slide",
	Long: `A new name, description, image format or pixel size can be given to a slide.

Change the description of slide ID 58a5a841e779890c2486ca90:

  slide update 58a5a841e779890c2486ca90 --name "DICOM GW9" --desc "A better description"

Change the pixel size with

  slide update 58a5a841e779890c2486ca90 --pixel-size 2 --pixel-size-unit um

Change previous image format to "ndpi"

  slide update 58a5a841e779890c2486ca90 --image-format ndpi`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing slide ID (string).")
		}
		if err := api.EditSlide(args[0],
			cfgSlideName,
			cfgSlideDescription,
			cfgPixelSizeValue,
			cfgPixelSizeUnit,
			cfgImageFormat); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Slide succcessfully updated.\n")
	},
}

func init() {
	slideCmd.AddCommand(slideupdateCmd)

	slideupdateCmd.Flags().StringVar(&cfgSlideName, "name", "", "Name")
	slideupdateCmd.Flags().StringVar(&cfgSlideDescription, "desc", "", "Short description")
	slideupdateCmd.Flags().StringVarP(&cfgImageFormat, "image-format", "f", "tiff", "Image Format")
	slideupdateCmd.Flags().Float64VarP(&cfgPixelSizeValue, "pixel-size", "p", 0, "Pixel size value")
	slideupdateCmd.Flags().StringVarP(&cfgPixelSizeUnit, "pixel-size-unit", "u", "um", "Pixel size unit")
}
