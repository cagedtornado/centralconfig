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

	"github.com/spf13/cobra"
)

var buildVersion = "Unknown"
var commitId string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version information",
	Long:  `Shows version information and exits`,
	Run: func(cmd *cobra.Command, args []string) {
		//	Show the version number
		fmt.Printf("Centralconfig version %s", buildVersion)

		//	Show the commitid if available:
		if commitId != "" {
			fmt.Printf(" (%s)", commitId[:7])
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
