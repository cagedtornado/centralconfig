package datastores

import (
	"encoding/json"
	"fmt"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//	The MysqlDB database information
type MySqlDB struct {
	Protocol string
	Address  string
	Database string
	User     string
	Password string
}

func (store MySqlDB) InitStore(overwrite bool) error {

	return nil
}

func (store MySqlDB) Get(configItem *ConfigItem) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	//	Open the database:
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s", store.User, store.Password, store.Protocol, store.Address, store.Database))
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Query to get default first
	rows, err := db.Query("select id, application, name, value, machine, updated where application=? and name=? and machine=?", configItem.Application, configItem.Name, "*")
	if err != nil {
		return retval, err
	}

	for rows.Next() {
		var id int
		var application string
		var name string
		var value string
		var machine string

		//	Scan the row into our variables
		err = rows.Scan(&id, &application, &name, &value, &machine)

		if err != nil {
			return retval, err
		}

		//	Set our return value
		retval.Id = id
		retval.Application = application
		retval.Name = name
		retval.Value = value
		retval.Machine = machine
	}

	return retval, nil
}

func (store MySqlDB) GetAllForApplication(application string) ([]ConfigItem, error) {
	//	Our return items:
	retval := []ConfigItem{}

	return retval, nil
}

func (store MySqlDB) GetAll() ([]ConfigItem, error) {
	//	Our return items:
	retval := []ConfigItem{}

	return retval, nil
}

func (store MySqlDB) GetAllApplications() ([]string, error) {
	//	Our return items:
	var retval []string

	return retval, nil
}

func (store MySqlDB) Set(configItem *ConfigItem) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	return retval, nil
}

func (store MySqlDB) Remove(configItem *ConfigItem) error {
	return nil
}
