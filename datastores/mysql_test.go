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

//	MySQL init should ping the database
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

//	MySQL get should return successfully even if the item doesn't exist
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

//	MySQL set should work
func TestMysql_Set_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()

	//	Need to reset database each time -- using TRUNCATE?
	//	https://dev.mysql.com/doc/refman/5.0/en/truncate-table.html

	//	Try storing some config items:
	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	response, err := db.Set(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Set failed: MySQL should have set an item without error: %s", err)
	}

	if ct1.Id == response.Id {
		t.Error("Set failed: MySQL should have set an item with the correct id")
	}
}
