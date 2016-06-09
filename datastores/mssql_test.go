package datastores_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/cagedtornado/centralconfig/datastores"
)

//	Gets the database connection information from the environment
func getMSSQLDBConnection() datastores.MSSqlDB {

	//	Set this information from environment variables?
	return datastores.MSSqlDB{
		Address:  os.Getenv("centralconfig_test_mssql_server"), /* Ex: test-server:3306 If this is blank, it assumes a local database on port 3306 */
		Database: os.Getenv("centralconfig_test_mssql_database"),
		User:     os.Getenv("centralconfig_test_mssql_user"),
		Password: os.Getenv("centralconfig_test_mssql_password")}
}

//	Reset the test database
func resetMSSQLTestDB(store datastores.MSSqlDB) {
	//	Open the database:
	db, _ := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
	defer db.Close()

	db.Exec("TRUNCATE TABLE configitem")
}

//	MSSQL init should ping the database
func TestMssql_Init_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	//	Act
	err := db.InitStore(true)

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Can't connect to database: %s", err)
	}
}

//	MSSQL get should return successfully even if the item doesn't exist
func TestMssql_Get_ItemDoesntExist_Successful(t *testing.T) {

	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	query := datastores.ConfigItem{
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

//	MSSQL set should work
func TestMssql_Set_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	//	Try storing some config items:
	ct1 := datastores.ConfigItem{
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

//	MSSQL set then get should work
func TestMssql_Set_ThenGet_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	query := datastores.ConfigItem{
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

//	MSSQL set then get should work - and default to global settings if
//	app-specific settings aren't found
func TestMssql_Set_ThenGet_Global_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := datastores.ConfigItem{
		Application: "*",
		Name:        "TestItem3",
		Value:       "Value2"}

	query := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem3"} // This item wasn't set for MyTestAppName - the global default should be used

	query2 := datastores.ConfigItem{
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

//	MSSQL set then get should work - and default to global settings if
//	app-specific settings aren't found
func TestMssql_Set_ThenGet_WithMachine_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct_nomachine := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2_NO_MACHINE"}

	ct_withmachine := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Machine:     "APPBOX1", // This config item is machine specific
		Value:       "Value2_SET_WITH_MACHINE"}

	query := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2"}

	query2 := datastores.ConfigItem{
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

//	MSSQL GetAllForApplication should work
func TestMssql_GetAllForApplication_NoMachine_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	query := datastores.ConfigItem{
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

//	Mssql GetAllForApplication should work - even when some items have a specified machine
func TestMssql_GetAllForApplication_WithMachine_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Machine:     "APPBOX1",
		Value:       "Value2"}

	query := datastores.ConfigItem{
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

//	MSSQL getall should work
func TestMssql_GetAll_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := datastores.ConfigItem{
		Application: "OtherTestApp",
		Name:        "TestItem3",
		Value:       "Value2"}

	ct4 := datastores.ConfigItem{
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

//	MSSQL getall should work - even with no initial data
func TestMssql_GetAll_NoInitialData_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	response, err := db.GetAll()

	//	Assert
	if err != nil {
		t.Errorf("GetAll (no initial data) failed: Should have returned all config items without error: %s", err)
	}

	if len(response) != 0 {
		t.Error("GetAll (no initial data) failed: Should have returned 0 items")
	}
}

//	MSSQL GetAllApplications should work
func TestMssql_GetAllApplications_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2",
		Value:       "Value2"}

	ct3 := datastores.ConfigItem{
		Application: "OtherTestApp",
		Name:        "TestItem3",
		Value:       "Value2"}

	ct4 := datastores.ConfigItem{
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

//	MSSQL GetAllApplications should work, even with no data
func TestMssql_GetAllApplications_NoData_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

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

//	MSSQL remove should work - even with a non-existant item
func TestMssql_Remove_ItemDoesntExist_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
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

//	MSSQL remove should work
func TestMssql_Remove_NoMachine_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
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

//	MSSQL remove should work, even with machine specified
func TestMssql_Remove_WithMachine_Successful(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL tests: Not on Windows")
	}

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

	ct1 := datastores.ConfigItem{
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
