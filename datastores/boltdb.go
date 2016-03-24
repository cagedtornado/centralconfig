package datastores

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

//	The BoltDB database information
type BoltDB struct {
	Database string
	User     string
	Password string
}

//	If we need to list applications, we can do so by listing buckets:
//	https://github.com/boltdb/bolt/issues/295

func (store BoltDB) InitStore(overwrite bool) error {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()

	return err
}

func (store BoltDB) Get(configItem *ConfigItem) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get the key based on the 'global' application
	err = db.View(func(tx *bolt.Tx) error {
		//	Get the item from the bucket with the app name
		b := tx.Bucket([]byte("*"))

		if b != nil {
			//	TODO: If we have a machine name, encode it in the key:
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

	//	Get the key based on the application name
	err = db.View(func(tx *bolt.Tx) error {
		//	Get the item from the bucket with the app name
		b := tx.Bucket([]byte(configItem.Application))

		if b != nil {
			//	If we have a machine name, append it in the key:
			keyName := configItem.Name
			if configItem.Machine != "" {
				keyName = keyName + "|" + configItem.Machine
			}

			configBytes := b.Get([]byte(keyName))

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

	return retval, err
}

func (store BoltDB) GetAllForApplication(application string) ([]ConfigItem, error) {
	//	Our return items:
	retval := []ConfigItem{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get the global app data first:
	err = db.View(func(tx *bolt.Tx) error {

		//	Get the items from the bucket with the app name
		b := tx.Bucket([]byte("*"))

		if b != nil {

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
		}

		return nil
	})

	//	Get the key based on the application
	err = db.View(func(tx *bolt.Tx) error {

		//	Get the items from the bucket with the app name
		b := tx.Bucket([]byte(application))

		if b != nil {

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
		}

		return nil
	})

	//	TODO: get application + machine

	return retval, err
}

func (store BoltDB) GetAll() ([]ConfigItem, error) {

	//	Our return items:
	retval := []ConfigItem{}
	var bucketList []string

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get a list of all buckets
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			bucketList = append(bucketList, string(name))
			return nil
		})
	})

	//	For each bucket...
	for _, bucketName := range bucketList {

		//	... get a list of configitems for that bucket
		err = db.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(bucketName))

			if b != nil {

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
			}

			return nil
		})
	}

	return retval, err
}

func (store BoltDB) GetAllApplications() ([]string, error) {

	//	Our return items:
	var bucketList []string

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return bucketList, err
	}

	//	Get a list of all buckets
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			bucketList = append(bucketList, string(name))
			return nil
		})
	})

	return bucketList, err
}

func (store BoltDB) Set(configItem *ConfigItem) (ConfigItem, error) {

	//	Our return item:
	retval := *configItem

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {
		//	Put the item in the bucket with the app name
		b, err := tx.CreateBucketIfNotExists([]byte(configItem.Application))
		if err != nil {
			return err
		}

		// If we don't have an id, generate an id for the configitem.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		if configItem.Id == 0 {
			id, _ := b.NextSequence()
			configItem.Id = int64(id)
		}

		//	Set the current datetime:
		configItem.LastUpdated = time.Now()

		//	Serialize to JSON format
		encoded, err := json.Marshal(configItem)
		if err != nil {
			return err
		}

		//	If we have a machine name, append it in the key:
		keyName := configItem.Name
		if configItem.Machine != "" {
			keyName = keyName + "|" + configItem.Machine
		}

		//	Store it, with the 'name' as the key:
		return b.Put([]byte(keyName), encoded)
	})

	return retval, err
}

func (store BoltDB) Remove(configItem *ConfigItem) error {

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {

		//	Get the item from the bucket with the app name
		b := tx.Bucket([]byte(configItem.Application))

		if b != nil {
			//	Delete it, with the 'name' as the key:
			return b.Delete([]byte(configItem.Name))
		}

		return nil
	})

	return err
}
