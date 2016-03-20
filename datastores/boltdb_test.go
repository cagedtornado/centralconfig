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
		Application: "Formbuilder",
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
		Application: "Formbuilder",
		Name:        "TestItem1",
		Value:       "Value1"}

	//	Act
	err := db.Set(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: BoltDB should have set an item without error: %s", err)
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
		Application: "Formbuilder",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "Formbuilder",
		Name:        "TestItem2",
		Value:       "Value2"}

	query := &datastores.ConfigItem{
		Application: "Formbuilder",
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

//	Bolt getall should work
func TestBoltDB_GetAll_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := datastores.BoltDB{
		Database: filename}

	ct1 := &datastores.ConfigItem{
		Application: "Formbuilder",
		Name:        "TestItem1",
		Value:       "Value1"}

	ct2 := &datastores.ConfigItem{
		Application: "Formbuilder",
		Name:        "TestItem2",
		Value:       "Value2"}

	query := &datastores.ConfigItem{
		Application: "Formbuilder"}

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
