package datastores_test

import (
	"fmt"
	"os"
	"testing"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/danesparza/centralconfig/datastores"
)

//	Gets the database connection information from the environment
func getDBConnection() datastores.MySqlDB {

	//	Set this information from environment variables?
	return datastores.MySqlDB{
		Address:  os.Getenv("centralconfig_test_mysql_server"), /* Ex: test-server:3306 If this is blank, it assumes a local database on port 3306 */
		Database: os.Getenv("centralconfig_test_mysql_database"),
		User:     os.Getenv("centralconfig_test_mysql_user"),
		Password: os.Getenv("centralconfig_test_mysql_password")}
}

//	Reset the test database
func resetTestDB(store datastores.MySqlDB) {
	//	Open the database:
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s", store.User, store.Password, store.Protocol, store.Address, store.Database))
	defer db.Close()

	db.Exec("TRUNCATE TABLE configitem")
}

//	MySQL init should ping the database
func TestMysql_Init_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	//	Act
	err := db.InitStore(true)

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Can't connect to database: %s", err)
	}
}

//	MySQL get should return successfully even if the item doesn't exist
func TestMysql_Get_ItemDoesntExist_Successful(t *testing.T) {

	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	query := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2"}

	//	Act
	response, err := db.Get(query)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned an empty dataset without error: %s", err)
	}

	if query.Value != response.Value && response.Value != "" {
		t.Errorf("Get failed: Shouldn't have returned the value %s", response.Value)
	}
}

//	MySQL set should work
func TestMysql_Set_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	//	Try storing some config items:
	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	response, err := db.Set(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Set failed: Should have set an item without error: %s", err)
	}

	if ct1.Id == response.Id {
		t.Error("Set failed: Should have set an item with the correct id")
	}
}

//	MySQL set then get should work
func TestMysql_Set_ThenGet_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	query := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2"}

	//	Act
	db.Set(ct1)
	db.Set(ct2)
	response, err := db.Get(query)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned a config item without error: %s", err)
	}

	if response.Value != ct2.Value {
		t.Errorf("Get failed: Should have returned the value %s but returned %s instead", ct2.Value, response.Value)
	}
}

//	MySQL set then get should work - and default to global settings if
//	app-specific settings aren't found
func TestMysql_Set_ThenGet_Global_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := &datastores.ConfigItem{
		Application: "*",
		Name:        "TestItem3",
		Value:       "Value2"}

	query := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem3"} // This item wasn't set for MyTestAppName - the global default should be used

	query2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2"} // This item WAS set for MyTestAppName

	//	Act
	db.Set(ct1)
	db.Set(ct2)
	db.Set(ct3)
	response, err := db.Get(query)
	response2, err2 := db.Get(query2)

	//	Assert
	if err != nil || err2 != nil {
		t.Errorf("Get failed: Should have returned a config item without error: %s", err)
	}

	if response.Value != ct3.Value {
		t.Errorf("Get (global) failed: Should have returned the value %s but returned %s instead", ct3.Value, response.Value)
	}

	if response2.Value != ct2.Value {
		t.Errorf("Get failed: Should have returned the value %s but returned %s instead", ct2.Value, response2.Value)
	}
}

//	MySQL set then get should work - and default to global settings if
//	app-specific settings aren't found
func TestMysql_Set_ThenGet_WithMachine_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct_nomachine := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2_NO_MACHINE"}

	ct_withmachine := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Machine:     "APPBOX1", // This config item is machine specific
		Value:       "Value2_SET_WITH_MACHINE"}

	query := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2"}

	query2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Machine:     "APPBOX1", // Notice that this has a machine in the query
		Name:        "TestItem2"}

	//	Act
	db.Set(ct1)
	db.Set(ct_nomachine)
	db.Set(ct_withmachine)
	response, err := db.Get(query)
	response2, err2 := db.Get(query2)

	//	Assert
	if err != nil || err2 != nil {
		t.Errorf("Get failed: Should have returned a config item without error: %s", err)
	}

	if response.Value != ct_nomachine.Value || response.Machine != ct_nomachine.Machine {
		t.Errorf("Get (not machine specific) failed: Should have returned the value %s (for machine %s) but returned %s (for machine %s) instead", ct_nomachine.Value, ct_nomachine.Machine, response.Value, response.Machine)
	}

	if response2.Value != ct_withmachine.Value || response2.Machine != ct_withmachine.Machine {
		t.Errorf("Get (machine specific) failed: Should have returned the value %s (for machine %s) but returned %s (for machine %s) instead", ct_withmachine.Value, ct_withmachine.Machine, response2.Value, response2.Machine)
	}
}

//	MySQL GetAllForApplication should work
func TestMysql_GetAllForApplication_NoMachine_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	query := &datastores.ConfigItem{
		Application: "MyTestAppName"}

	//	Act
	db.Set(ct1)
	db.Set(ct2)
	response, err := db.GetAllForApplication(query.Application)

	//	Assert
	if err != nil {
		t.Errorf("GetAllForApplication failed: Should have returned all config items without error: %s", err)
	}

	if len(response) != 2 {
		t.Error("GetAllForApplication failed: Should have returned 2 items")
	}
}

//	Mysql GetAllForApplication should work - even when some items have a specified machine
func TestMysql_GetAllForApplication_WithMachine_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Machine:     "APPBOX1",
		Value:       "Value2"}

	query := &datastores.ConfigItem{
		Application: "MyTestAppName"}

	//	Act
	db.Set(ct1)
	db.Set(ct2)
	db.Set(ct3)
	response, err := db.GetAllForApplication(query.Application)

	//	Assert
	if err != nil {
		t.Errorf("GetAllForApplication failed: Should have returned all config items without error: %s", err)
	}

	if len(response) != 3 {
		t.Error("GetAllForApplication failed: Should have returned 3 items")
	}
}

//	MySQL getall should work
func TestMysql_GetAll_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := &datastores.ConfigItem{
		Application: "OtherTestApp",
		Name:        "TestItem3",
		Value:       "Value2"}

	ct4 := &datastores.ConfigItem{
		Application: "*",
		Name:        "TestItem4",
		Value:       "Value2"}

	//	Act
	db.Set(ct1)
	db.Set(ct2)
	db.Set(ct3)
	db.Set(ct4)
	response, err := db.GetAll()

	//	Assert
	if err != nil {
		t.Errorf("GetAll failed: Should have returned all config items without error: %s", err)
	}

	if len(response) != 4 {
		t.Error("GetAll failed: Should have returned 4 items")
	}
}

//	MySQL getall should work - even with no initial data
func TestMysql_GetAll_NoInitialData_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	response, err := db.GetAll()

	//	Assert
	if err != nil {
		t.Errorf("GetAll (no initial data) failed: Should have returned all config items without error: %s", err)
	}

	if len(response) != 0 {
		t.Error("GetAll (no initial data) failed: Should have returned 0 items")
	}
}

//	MySQL GetAllApplications should work
func TestMysql_GetAllApplications_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := &datastores.ConfigItem{
		Application: "OtherTestApp",
		Name:        "TestItem3",
		Value:       "Value2"}

	ct4 := &datastores.ConfigItem{
		Application: "*",
		Name:        "TestItem4",
		Value:       "Value2"}

	//	Act
	db.Set(ct1)
	db.Set(ct2)
	db.Set(ct3)
	db.Set(ct4)
	response, err := db.GetAllApplications()

	//	Assert
	if err != nil {
		t.Errorf("GetAllApplications failed: Should have returned all applications without error: %s", err)
	}

	if len(response) != 3 {
		t.Error("GetAllApplications failed: Should have returned the correct number of applications")
	}
}

//	MySQL GetAllApplications should work, even with no data
func TestMysql_GetAllApplications_NoData_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	//	Act
	response, err := db.GetAllApplications()

	//	Assert
	if err != nil {
		t.Errorf("GetAllApplications failed: Should have returned all applications without error: %s", err)
	}

	if len(response) != 0 {
		t.Error("GetAllApplications failed: Should have returned the correct number of applications")
	}
}

//	MySQL remove should work - even with a non-existant item
func TestMysql_Remove_ItemDoesntExist_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	err := db.Remove(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Remove failed: Should have attempted to remove a non-existant item without error: %s", err)
	}
}

//	MySQL remove should work
func TestMysql_Remove_NoMachine_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	db.Set(ct1)
	err := db.Remove(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Remove failed: Should have removed an item without error: %s", err)
	}
}

//	MySQL remove should work, even with machine specified
func TestMysql_Remove_WithMachine_Successful(t *testing.T) {
	//	Arrange
	db := getDBConnection()
	resetTestDB(db)

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Machine:     "APPBOX1",
		Value:       "Value1"}

	//	Act
	db.Set(ct1)
	err := db.Remove(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Remove failed: Should have removed an item without error: %s", err)
	}
}
