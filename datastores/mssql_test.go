package datastores_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/danesparza/centralconfig/datastores"
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

//	MySQL get should return successfully even if the item doesn't exist
func TestMssql_Get_ItemDoesntExist_Successful(t *testing.T) {

	//	Arrange
	db := getMSSQLDBConnection()
	resetMSSQLTestDB(db)

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
