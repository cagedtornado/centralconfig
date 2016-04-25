package datastores

import (
	"github.com/spf13/viper"
	"time"
)

//	ConfigItem represents a configuration item
type ConfigItem struct {
	Id          int64     `sql:"id" json:"id"`
	Application string    `sql:"application" json:"application"`
	Machine     string    `sql:"machine" json:"machine"`
	Name        string    `sql:"name" json:"name"`
	Value       string    `sql:"value" json:"value"`
	LastUpdated time.Time `sql:"updated" json:"updated"`
}

type ConfigResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//	ConfigService encapsulates account (user) based operations
//	This allows us to create a testable service layer.  See
//	https://github.com/tonyhb/tonyhb.com/blob/master/posts/Building%20a%20testable%20Golang%20database%20layer.md
//	for more information
type ConfigService interface {

	//	Initialize the store (create the DDL if necessary)
	InitStore(overwrite bool) error

	//	Create / update a config item
	Set(c *ConfigItem) (ConfigItem, error)

	//	Get a specific config item
	Get(c *ConfigItem) (ConfigItem, error)

	//	Get all config items for the given application
	GetAllForApplication(application string) ([]ConfigItem, error)

	//	Get all config items for all applications (including global)
	GetAll() ([]ConfigItem, error)

	//	Get all applications (including global)
	GetAllApplications() ([]string, error)

	//	Remove a config item
	Remove(c *ConfigItem) error
}

//	Get the currently configured datastore
func GetConfigDatastore() ConfigService {

	//	Get configuration information and return the appropriate
	//	provider based on what is configured
	if viper.InConfig("datastore") {
		dsConfig := viper.Sub("datastore")

		//	If we have MSSQL, use that:
		if dsConfig.InConfig("mssql") {
			return MSSqlDB{
				Database: dsConfig.GetString("mssql.database"),
				Address:  dsConfig.GetString("mssql.address"),
				User:     dsConfig.GetString("mssql.user"),
				Password: dsConfig.GetString("mssql.password")}
		}

		//	If we have MySQL, use that:
		if dsConfig.InConfig("mysql") {
			return MySqlDB{
				Database: dsConfig.GetString("mysql.database"),
				Address:  dsConfig.GetString("mysql.address"),
				User:     dsConfig.GetString("mysql.user"),
				Password: dsConfig.GetString("mysql.password")}
		}
	}

	//	Last resort: use BoltDB
	return BoltDB{
		Database: viper.GetString("datastore.boltdb.database")}

}
