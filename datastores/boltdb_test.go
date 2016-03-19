package datastores

import (
	"os"
	"testing"
)

//	Bolt init should create a new BoltDB file
func TestBoltDB_Init(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := BoltDB{
		Database: filename}

	//	Act
	db.InitStore(true)

	//	Assert
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Init failed: BoltDB file %s was not created", filename)
	}
}
