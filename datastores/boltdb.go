package datastores

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

//	The BoltDB database information
type BoltDB struct {
	Database string
	User     string
	Password string
}

func (store BoltDB) InitStore(overwrite bool) error {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	checkErr(err)
	defer db.Close()

	return err
}

func (store BoltDB) Get(configItem ConfigItem) (ConfigItem, error) {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	checkErr(err)
	defer db.Close()

	//	Our return item:
	retval := ConfigItem{}

	//	Need to decide on appropriate nesting structure for BoltDB...
	//	Perhaps {appname} or {} -> Configitems?
	//	                      ^ this means 'global'
	err = db.View(func(tx *bolt.Tx) error {
		//	Get the item from the bucket with the app name
		b := tx.Bucket([]byte(configItem.Application))

		//	Get the item based on the config name:
		configBytes := b.Get([]byte(configItem.Name))

		//	Unmarshal data into our config item
		if err := json.Unmarshal(configBytes, &retval); err != nil {
			return err
		}

		return nil
	})

	checkErr(err)

	return retval, err
}

func (store BoltDB) Set(configItem ConfigItem) error {

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	checkErr(err)
	defer db.Close()

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {
		//	Put the item in the bucket with the app name
		b, err := tx.CreateBucketIfNotExists([]byte(configItem.Application))
		checkErr(err)

		//	Serialize to JSON format
		encoded, err := json.Marshal(configItem)
		checkErr(err)

		//	Store it, with the 'name' as the key:
		return b.Put([]byte(configItem.Name), encoded)
	})

	return err
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
