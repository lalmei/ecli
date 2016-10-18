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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.3.2"

var cfgFile string

func showHelpAndExit(cmd *cobra.Command, msg string) {
	fmt.Printf("%s\n\n", msg)
	fmt.Println(cmd.UsageString())
	os.Exit(1)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ecli",
	Short: "Simple client for Keen Eye Engine",
	Long: `Ecli is a CLI client for the Keen Eye Engine. It can be used to upload,
find a slide or get slide information.

Also it is possible no browse the work list like browsing a UNIX filesystem.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetOutput(os.Stderr)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.ecli.json)")
	RootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode, no verbose output")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ecli") // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(`Please first create a config file in $HOME/.ecli.json with the following content:
{
  "profile1": {
    "login": "something",
    "password: "else",
	"url": "https://somehost.keeneyetechnologies.com/api/v2",
  },
  "profile2": {
    "login": "something",
    "password: "else",
	"url": "https://somehost.keeneyetechnologies.com/api/v2",
  }
}`)
		os.Exit(1)
	}
}
