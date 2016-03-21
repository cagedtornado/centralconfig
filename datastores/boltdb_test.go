package datastores_test

import (
	"os"
	"testing"

	"github.com/danesparza/centralconfig/datastores"
)

//	Sanity check: The database shouldn't exist yet
func TestBoltDB_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	filename := "testing.db"

	//	Act

	//	Assert
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("BoltDB database file check failed: BoltDB file %s already exists, and shouldn't", filename)
	}
}

//	Bolt init should create a new BoltDB file
func TestBoltDB_Init_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

	//	Act
	db.InitStore(true)

	//	Assert
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Init failed: BoltDB file %s was not created", filename)
	}
}

//	Bolt get should return successfully even if the item doesn't exist
func TestBoltDB_Get_ItemDoesntExist_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

	query := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem2"}

	//	Act
	response, err := db.Get(query)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: BoltDB should have returned an empty dataset without error: %s", err)
	}

	if query.Value != response.Value && response.Value != "" {
		t.Errorf("Get failed: BoltDB shouldn't have returned the value %s", response.Value)
	}
}

//	Bolt set should work
func TestBoltDB_Set_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

	//	Try storing some config items:
	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	err := db.Set(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Set failed: BoltDB should have set an item without error: %s", err)
	}
}

//	Bolt set then get should work
func TestBoltDB_Set_ThenGet_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

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
		t.Errorf("Get failed: BoltDB should have returned a config item without error: %s", err)
	}

	if response.Value != ct2.Value {
		t.Errorf("Get failed: BoltDB should have returned the value %s but returned %s instead", ct2.Value, response.Value)
	}
}

//	Bolt set then get should work - and default to global settings if
//	app-specific settings aren't found
func TestBoltDB_Set_ThenGet_Global_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

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
		t.Errorf("Get failed: BoltDB should have returned a config item without error: %s", err)
	}

	if response.Value != ct3.Value {
		t.Errorf("Get (global) failed: BoltDB should have returned the value %s but returned %s instead", ct3.Value, response.Value)
	}

	if response2.Value != ct2.Value {
		t.Errorf("Get failed: BoltDB should have returned the value %s but returned %s instead", ct2.Value, response2.Value)
	}
}

//	Bolt set then get should work - and default to global settings if
//	app-specific settings aren't found
func TestBoltDB_Set_ThenGet_WithMachine_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

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
		t.Errorf("Get failed: BoltDB should have returned a config item without error: %s", err)
	}

	if response.Value != ct_nomachine.Value || response.Machine != ct_nomachine.Machine {
		t.Errorf("Get (not machine specific) failed: BoltDB should have returned the value %s (for machine %s) but returned %s (for machine %s) instead", ct_nomachine.Value, ct_nomachine.Machine, response.Value, response.Machine)
	}

	if response2.Value != ct_withmachine.Value || response2.Machine != ct_withmachine.Machine {
		t.Errorf("Get (machine specific) failed: BoltDB should have returned the value %s (for machine %s) but returned %s (for machine %s) instead", ct_withmachine.Value, ct_withmachine.Machine, response2.Value, response2.Machine)
	}
}

//	Bolt getall should work
func TestBoltDB_GetAll_NoMachine_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

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
	response, err := db.GetAll(query.Application)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: BoltDB should have returned all config items without error: %s", err)
	}

	if len(response) != 2 {
		t.Error("Get failed: BoltDB should have returned 2 items")
	}
}

//	Bolt getall should work - even when some items have a specified machine
func TestBoltDB_GetAll_WithMachine_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

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
	response, err := db.GetAll(query.Application)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: BoltDB should have returned all config items without error: %s", err)
	}

	if len(response) != 3 {
		t.Error("Get failed: BoltDB should have returned 3 items")
	}
}

//	Bolt remove should work - even with a non-existant item
func TestBoltDB_Remove_ItemDoesntExist_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	err := db.Remove(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Remove failed: BoltDB should have attempted to remove a non-existant item without error: %s", err)
	}
}

//	Bolt remove should work
func TestBoltDB_Remove_NoMachine_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

	ct1 := &datastores.ConfigItem{
		Application: "MyTestAppName",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	db.Set(ct1)
	err := db.Remove(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Remove failed: BoltDB should have removed an item without error: %s", err)
	}
}

//	Bolt remove should work, even with machine specified
func TestBoltDB_Remove_WithMachine_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

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
		t.Errorf("Remove failed: BoltDB should have removed an item without error: %s", err)
	}
}
