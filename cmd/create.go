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
	Long:  `Use this command to create a default configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("create called")

		if jsonConfig {
			fmt.Println("Creating JSON config")
		} else if yamlConfig {
			//	Write the yamlDefault to the location of the default
			fmt.Println("Writing config to ?")
			fmt.Println("Creating YAML config")
		}
	},
}

func init() {
	configCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	createCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")

}
