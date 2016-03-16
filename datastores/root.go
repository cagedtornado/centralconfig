package datastores

import (
	"time"

	"github.com/spf13/viper"
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
	Status  int64        `json: "status"`
	Message string       `json: "message"`
	Data    []ConfigItem `json: "data"`
}

//	ConfigService encapsulates account (user) based operations
//	This allows us to create a testable service layer.  See
//	https://github.com/tonyhb/tonyhb.com/blob/master/posts/Building%20a%20testable%20Golang%20database%20layer.md
//	for more information
type ConfigService interface {

	//	Initialize the store (create the DDL if necessary)
	InitStore(overwrite bool) error

	//	Create / update a config item
	Set(c *ConfigItem) error

	//	Get a specific config item
	Get(c *ConfigItem) (ConfigItem, error)

	//	Get all config items for the given application
	GetAll(application string) ([]ConfigItem, error)
}

//	Get the currently configured datastore
func GetConfigDatastore() ConfigService {

	//	Get configuration information and return the appropriate
	//	provider based on what is configured

	//	For now, just always return the BoltDB provider
	return BoltDB{
		Database: viper.GetString("boltdb.database")}

}
