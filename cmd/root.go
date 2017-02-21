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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	cfgQuiet bool
)

func errorExit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func usageErrorExit(cmd *cobra.Command, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Printf("\n\n")
	cmd.Help()
	os.Exit(1)
}

func showHelpAndExit(cmd *cobra.Command, msg string) {
	fmt.Printf("%s\n\n", msg)
	cmd.Help()
	os.Exit(1)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ecli",
	Short: "Simple client for Keen Eye API",
	Long: `Ecli is a command line client for the Keen Eye API. Its primary use is to
perform slide upload and get slide information in a convenient way.`,
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
	RootCmd.PersistentFlags().BoolVarP(&cfgQuiet, "quiet", "q", false, "Quiet mode, no verbose output")
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

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println(`Please first create a config file in $HOME/.ecli.json (or anywhere using --config) with the following content adapted to your configuration:
{
  "platform": {
    "login": "something",
    "password: "pwd",
	"url": "https://somehost.keeneyetechnologies.com/api/v2",
  }
}
Then log in with
  $ ecli login platform`)
			os.Exit(1)
		}
	}
}
