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

//	ConfigResponse represents an API response
type ConfigResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//	WebSocketResponse represents a WebSocket event response
type WebSocketResponse struct {
	Type string     `json:"type"`
	Data ConfigItem `json:"data"`
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
	dsType := viper.GetString("datastore.type")

	if dsType != "" {

		switch dsType {
		case "mssql":
			//	If we have MSSQL, use that:
			return MSSqlDB{
				Database: viper.GetString("datastore.database"),
				Address:  viper.GetString("datastore.address"),
				User:     viper.GetString("datastore.user"),
				Password: viper.GetString("datastore.password")}

		case "mysql":
			//	If we have MySQL, use that:
			return MySqlDB{
				Database: viper.GetString("datastore.database"),
				Address:  viper.GetString("datastore.address"),
				User:     viper.GetString("datastore.user"),
				Password: viper.GetString("datastore.password")}

		case "boltdb":
			return BoltDB{
				Database: viper.GetString("datastore.database")}
		}
	}

	return UnknownDB{}

}
