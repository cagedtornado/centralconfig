package cmd

import (
	"fmt"

	"github.com/danesparza/centralconfig/datastores"
	"github.com/spf13/cobra"
)

var (
	jsonConfig bool
	yamlConfig bool
	mysqlDDL   bool
	mssqlDDL   bool
)

var yamlDefault = []byte(`
server:
  port: 3000
datastore:
  boltdb:
    database: config.db
`)

var jsonDefault = []byte(`{
  "server" : {
    "port": "3000"
  },
  "datastore" : {
    "boltdb": {
      "database": "config.db"
    }
  }
}`)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Prints default server configuration files or DDL",
	Long: `Use this to create a default configuration file for the centralconfig server. 

Example:

centralconfig defaults > centralconfig.yaml

You can also use this to print SQL database scripts to create your datastore.  Example: 

centralconfig defaults --mysql > centralconfigdb.sql

`,
	Run: func(cmd *cobra.Command, args []string) {

		if jsonConfig {
			fmt.Printf("%s", jsonDefault)
		} else if mysqlDDL {
			fmt.Printf("%s", datastores.GetMysqlCreateDDL())
		} else if mssqlDDL {
			fmt.Printf("%s", datastores.GetMSsqlCreateDDL())
		} else if yamlConfig {
			fmt.Printf("%s", yamlDefault)
		}
	},
}

func init() {
	RootCmd.AddCommand(defaultsCmd)

	defaultsCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	defaultsCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")
	defaultsCmd.Flags().BoolVarP(&mysqlDDL, "mysql", "m", false, "Create a MySQL database script")
	defaultsCmd.Flags().BoolVarP(&mssqlDDL, "mssql", "s", false, "Create a MSSQL database script")

}
