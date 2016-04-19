package datastores_test

import (
	"os"
	"testing"

	"github.com/danesparza/centralconfig/datastores"
)

func getDBConnection() datastores.MySqlDB {

	//	Set this information from environment variables?
	return datastores.MySqlDB{
		Address:  os.Getenv("centralconfig_test_mysql_server"), /* Ex: test-server:3306 If this is blank, it assumes a local database on port 3306 */
		Database: os.Getenv("centralconfig_test_msyql_database"),
		User:     os.Getenv("centralconfig_test_mysql_user"),
		Password: os.Getenv("centralconfig_test_mysql_password")}
}

//	Bolt init should create a new BoltDB file
func TestMysql_Init_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()

	//	Act
	err := db.InitStore(true)

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Can't connect to MySQL database: %s", err)
	}
}

//	Bolt get should return successfully even if the item doesn't exist
func TestMysql_Get_ItemDoesntExist_Successful(t *testing.T) {

	//	Arrange
	db := getDBConnection()

	query := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2"}

	//	Act
	response, err := db.Get(query)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: MySQL should have returned an empty dataset without error: %s", err)
	}

	if query.Value != response.Value && response.Value != "" {
		t.Errorf("Get failed: MySQL shouldn't have returned the value %s", response.Value)
	}
}
