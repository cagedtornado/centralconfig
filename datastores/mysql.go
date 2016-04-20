package datastores

import (
	"fmt"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//	Requires at least MySQL 5.6 (for the auto updating datetime)
var dbCreate = []byte(`
CREATE TABLE configitem (
  id int(11) NOT NULL AUTO_INCREMENT,
  application varchar(100) NOT NULL DEFAULT '*',
  name varchar(100) NOT NULL,
  value longtext NOT NULL,
  machine varchar(100) NOT NULL DEFAULT '',
  updated datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY id_UNIQUE (id),
  UNIQUE KEY app_name_machine (application,name,machine),
  KEY idx_application (application)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;
`)

//	The MysqlDB database information
type MySqlDB struct {
	Protocol string
	Address  string
	Database string
	User     string
	Password string
}

func (store MySqlDB) InitStore(overwrite bool) error {
	//	Open the database:
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s", store.User, store.Password, store.Protocol, store.Address, store.Database))
	defer db.Close()
	if err != nil {
		return err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (store MySqlDB) Get(configItem *ConfigItem) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	//	Open the database:
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true", store.User, store.Password, store.Protocol, store.Address, store.Database))
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Prepare our query
	stmt, err := db.Prepare("select id, application, name, value, machine, updated from configitem where application=? and name=? and machine=?")
	if err != nil {
		return retval, err
	}
	defer stmt.Close()

	//	Get the application/name/machine combo
	rows, err := stmt.Query(configItem.Application, configItem.Name, configItem.Machine)
	if err != nil {
		return retval, err
	}

	for rows.Next() {
		var id int64
		var application string
		var name string
		var value string
		var machine string
		var updated time.Time

		//	Scan the row into our variables
		err = rows.Scan(&id, &application, &name, &value, &machine, &updated)

		if err != nil {
			return retval, err
		}

		//	Set our return value
		retval = ConfigItem{
			Id:          id,
			Application: application,
			Name:        name,
			Value:       value,
			Machine:     machine,
			LastUpdated: updated}

		break
	}

	//	If we haven't found it, get the application/name combo with a blank machine name
	if retval.Id == 0 {
		rows, err = stmt.Query(configItem.Application, configItem.Name, "")
		if err != nil {
			return retval, err
		}

		for rows.Next() {
			var id int64
			var application string
			var name string
			var value string
			var machine string
			var updated time.Time

			//	Scan the row into our variables
			err = rows.Scan(&id, &application, &name, &value, &machine, &updated)

			if err != nil {
				return retval, err
			}

			//	Set our return value
			retval = ConfigItem{
				Id:          id,
				Application: application,
				Name:        name,
				Value:       value,
				Machine:     machine,
				LastUpdated: updated}

			break
		}
	}

	//	If we still haven't found it, get the default application/name and blank machine name
	if retval.Id == 0 {
		rows, err = stmt.Query("*", configItem.Name, "")
		if err != nil {
			return retval, err
		}

		for rows.Next() {
			var id int64
			var application string
			var name string
			var value string
			var machine string
			var updated time.Time

			//	Scan the row into our variables
			err = rows.Scan(&id, &application, &name, &value, &machine, &updated)

			if err != nil {
				return retval, err
			}

			//	Set our return value
			retval = ConfigItem{
				Id:          id,
				Application: application,
				Name:        name,
				Value:       value,
				Machine:     machine,
				LastUpdated: updated}

			break
		}
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
