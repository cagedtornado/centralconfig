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
	defer db.Close()

	return err
}

func (store BoltDB) Get(configItem *ConfigItem) (ConfigItem, error) {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	checkErr(err)
	defer db.Close()

	//	Our return item:
	retval := ConfigItem{}

	//	TODO: Get global first, then get application, then get application + machine
	//	Get the key based on the application
	err = db.View(func(tx *bolt.Tx) error {
		//	Get the item from the bucket with the app name
		b := tx.Bucket([]byte(configItem.Application))

		if b != nil {
			//	Get the item based on the config name:
			configBytes := b.Get([]byte(configItem.Name))

			//	Need to make sure we got something back here before we try to unmarshal?
			if len(configBytes) > 0 {
				//	Unmarshal data into our config item
				if err := json.Unmarshal(configBytes, &retval); err != nil {
					return err
				}
			}
		}

		return nil
	})

	checkErr(err)

	return retval, err
}

func (store BoltDB) GetAll(application string) ([]ConfigItem, error) {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	checkErr(err)
	defer db.Close()

	//	Our return items:
	retval := []ConfigItem{}

	//	TODO: Get global first, then get application, then get application + machine
	//	Get the key based on the application
	err = db.View(func(tx *bolt.Tx) error {

		//	Get the items from the bucket with the app name
		b := tx.Bucket([]byte(application))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			//	Unmarshal data into our config item
			ci := ConfigItem{}
			if err := json.Unmarshal(v, &ci); err != nil {
				return err
			}

			//	Add the item to our list of config items
			retval = append(retval, ci)
		}

		return nil
	})

	checkErr(err)

	return retval, err
}

func (store BoltDB) Set(configItem *ConfigItem) error {

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
