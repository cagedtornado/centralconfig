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

var (
	jsonConfig bool
	yamlConfig bool
)

var yamlDefault = []byte(`
http:
  port: 3000
boltdb:
  database: config.db
`)

var jsonDefault = []byte(`{
"http" : {
		"port": "3000"
	},
"boltdb": {
        "database": "config.db"
    }
}`)

// createconfigCmd represents the createconfig command
var createconfigCmd = &cobra.Command{
	Use:   "createconfig",
	Short: "Creates a default server configuration file",
	Long: `Use this to create a default configuration file for the centralconfig server. 

Example:

centralconfig createconfig > centralconfig.yaml`,
	Run: func(cmd *cobra.Command, args []string) {

		if jsonConfig {
			fmt.Printf("%s", jsonDefault)
		} else if yamlConfig {
			fmt.Printf("%s", yamlDefault)
		}
	},
}

func init() {
	RootCmd.AddCommand(createconfigCmd)

	createconfigCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	createconfigCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")

}
