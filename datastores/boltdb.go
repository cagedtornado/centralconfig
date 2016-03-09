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

	//	Read the 'configitems' bucket
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("configitems"))

		//	Get the item with the expected key:
		v := b.Get([]byte(configItem.Id))
		fmt.Printf("%s", v)
		return nil
	})
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
