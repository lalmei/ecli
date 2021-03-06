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
	"runtime"

	"github.com/keeneyetech/ecli/api"
	"github.com/keeneyetech/ecli/core"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show tool version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:      %s\n", core.Version)
		fmt.Printf("Git revision: %s\n", core.GitRevision)
		fmt.Printf("Git branch:   %s\n", core.GitBranch)
		fmt.Printf("Go version:   %s\n", runtime.Version())
		fmt.Printf("Built:        %s\n", core.Built)
		fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
		v, err := api.Version()
		if err == nil {
			fmt.Printf("---\nAPI Version:  %s\n", v)
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
