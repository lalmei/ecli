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
	"strings"

	"ecli/api"
	"ecli/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open a session",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showHelpAndExit(cmd, fmt.Sprintf("Please provide a profile name among %q",
				strings.Join(viper.AllKeys(), ", ")))
		}
		creds := viper.GetStringMapString(args[0])
		if len(creds) == 0 {
			log.Fatalf("no profile named %q. Please check your config file.", args[0])
		}
		// Hack to first store a session file, giving an empty token.
		_, tok, err := config.LoadSession()
		ztok := ""
		if err == nil {
			ztok = tok
		}
		if err := config.StoreSession(creds["url"], ztok); err != nil {
			log.Fatalf("can't save session info: %q", err.Error())
		}
		token, err := api.OpenSession(creds["url"], creds["login"], creds["password"])
		if err != nil {
			log.Fatal(err)
		}
		if err := config.StoreSession(creds["url"], token); err != nil {
			log.Fatalf("can't save session info: %q", err.Error())
		}
		fmt.Printf("%s: you have been logged in successfully.\n", args[0])
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
