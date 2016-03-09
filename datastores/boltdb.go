package datastores

import (
	"fmt"

	"github.com/boltdb/bolt"
)

//	The BoltDB database information
type BoltDB struct {
	Database string
	User     string
	Password string
}

func (store BoltDB) Get(configItem ConfigItem) (ConfigItem, error) {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	checkErr(err)
	defer db.Close()

	//	For now, just return a blank config item.
	ci := ConfigItem{
		Id:    77,
		Name:  "bogus",
		Value: "fubar"}

	//	Need to decide on appropriate nesting structure for BoltDB...
	//	Perhaps {appname} or {} -> Configitems?
	//	                      ^ this means 'global'
	db.View(func(tx *bolt.Tx) error {
		// b := tx.Bucket([]byte("configitems"))

		//	Get the item with the expected key:
		// v := b.Get([]byte(configItem.Name))
		// fmt.Printf("%s", v)

		return nil
	})

	return ci, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
