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
boltdb:
  database: config.db
`)

var jsonDefault = []byte(`{
"boltdb": {
        "database": "config.db"
    }
}`)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a default configuration file",
	Long: `Use this command to print a default configuration file.  
	You can redirect the output to a file in order to create a default config file`,
	Run: func(cmd *cobra.Command, args []string) {

		if jsonConfig {
			fmt.Printf("%s", jsonDefault)
		} else if yamlConfig {
			fmt.Printf("%s", yamlDefault)
		}
	},
}

func init() {
	configCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	createCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")

}
