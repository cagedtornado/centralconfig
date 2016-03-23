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

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Prints a default server configuration file",
	Long: `Use this to create a default configuration file for the centralconfig server. 

Example:

centralconfig defaults > centralconfig.yaml`,
	Run: func(cmd *cobra.Command, args []string) {

		if jsonConfig {
			fmt.Printf("%s", jsonDefault)
		} else if yamlConfig {
			fmt.Printf("%s", yamlDefault)
		}
	},
}

func init() {
	RootCmd.AddCommand(defaultsCmd)

	defaultsCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	defaultsCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")

}
