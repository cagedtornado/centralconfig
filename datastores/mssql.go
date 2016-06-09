package datastores

import (
	"fmt"
	"time"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

//	Database creation DDL
var dbCreateMSSQL = []byte(`
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[configitem](
	[id] [bigint] IDENTITY(1,1) NOT NULL,
	[application] [nvarchar](100) NOT NULL CONSTRAINT [DF_configitem_application]  DEFAULT (N'*'),
	[name] [nvarchar](100) NOT NULL,
	[value] [nvarchar](max) NOT NULL,
	[machine] [nvarchar](100) NOT NULL CONSTRAINT [DF_configitem_machine]  DEFAULT (N''),
	[updated] [datetime] NOT NULL CONSTRAINT [DF_configitem_updated]  DEFAULT (getdate()),
 CONSTRAINT [PK_configitem] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY],
 CONSTRAINT [unique_app_name_machine] UNIQUE NONCLUSTERED 
(
	[application] ASC,
	[name] ASC,
	[machine] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]

GO
`)

//	The MSSQL database information
type MSSqlDB struct {
	Address  string
	Database string
	User     string
	Password string
}

func (store MSSqlDB) InitStore(overwrite bool) error {
	//	Open the database:
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
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

func (store MSSqlDB) Get(configItem ConfigItem) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	//	Open the database:
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Prepare our query
	stmt, err := db.Prepare("select id, application, name, value, machine, updated from configitem where application=? and name=? and machine=?")
	defer stmt.Close()
	if err != nil {
		return retval, err
	}

	//	Get the application/name/machine combo
	rows, err := stmt.Query(configItem.Application, configItem.Name, configItem.Machine)
	defer rows.Close()
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
		defer rows.Close()
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
		defer rows.Close()
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

func (store MSSqlDB) GetAllForApplication(application string) ([]ConfigItem, error) {
	//	Our return item:
	retval := []ConfigItem{}

	//	Open the database:
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Prepare our query
	stmt, err := db.Prepare("select id, application, name, value, machine, updated from configitem where application=?")
	defer stmt.Close()
	if err != nil {
		return retval, err
	}

	//	Get all global config items
	rows, err := stmt.Query("*")
	defer rows.Close()
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

		//	Append to return values
		retval = append(retval, ConfigItem{
			Id:          id,
			Application: application,
			Name:        name,
			Value:       value,
			Machine:     machine,
			LastUpdated: updated})
	}

	//	Get config items for the given application:
	rows, err = stmt.Query(application)
	defer rows.Close()
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

		//	Append to return values
		retval = append(retval, ConfigItem{
			Id:          id,
			Application: application,
			Name:        name,
			Value:       value,
			Machine:     machine,
			LastUpdated: updated})
	}

	return retval, nil
}

func (store MSSqlDB) GetAll() ([]ConfigItem, error) {
	//	Our return items:
	retval := []ConfigItem{}

	//	Open the database:
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get all config items
	rows, err := db.Query("select id, application, name, value, machine, updated from configitem")
	defer rows.Close()
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

		//	Append to return values
		retval = append(retval, ConfigItem{
			Id:          id,
			Application: application,
			Name:        name,
			Value:       value,
			Machine:     machine,
			LastUpdated: updated})
	}

	return retval, nil
}

func (store MSSqlDB) GetAllApplications() ([]string, error) {
	//	Our return items:
	var retval []string

	//	Open the database:
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get all applications
	rows, err := db.Query("select distinct application from configitem")
	defer rows.Close()
	if err != nil {
		return retval, err
	}

	for rows.Next() {
		var application string

		//	Scan the row into our variables
		err = rows.Scan(&application)

		if err != nil {
			return retval, err
		}

		//	Append to return values
		retval = append(retval, application)
	}

	return retval, nil
}

func (store MSSqlDB) Set(configItem ConfigItem) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	//	Open the database:
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
	defer db.Close()
	if err != nil {
		return retval, err
	}

	if configItem.Id == 0 {
		//	If we have a brand new item, insert it
		stmt, err := db.Prepare("insert into configitem(application, name, value, machine) values(?, ?, ?, ?)")
		defer stmt.Close()
		if err != nil {
			return retval, err
		}

		res, err := stmt.Exec(configItem.Application, configItem.Name, configItem.Value, configItem.Machine)
		if err != nil {
			return retval, err
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return retval, err
		}

		retval = ConfigItem{
			Id:          lastId,
			Application: configItem.Application,
			Name:        configItem.Name,
			Value:       configItem.Value,
			Machine:     configItem.Machine,
			LastUpdated: time.Now()}

	} else {
		//	If we have an existing id, just update the old item
		stmt, err := db.Prepare("update configitem set application=?, name=?, value=?, machine=? where id=?")
		defer stmt.Close()
		if err != nil {
			return retval, err
		}

		_, err = stmt.Exec(configItem.Application, configItem.Name, configItem.Value, configItem.Machine, configItem.Id)
		if err != nil {
			return retval, err
		}

		retval = ConfigItem{
			Id:          configItem.Id,
			Application: configItem.Application,
			Name:        configItem.Name,
			Value:       configItem.Value,
			Machine:     configItem.Machine,
			LastUpdated: time.Now()}
	}

	return retval, nil
}

func (store MSSqlDB) Remove(configItem ConfigItem) error {
	//	Open the database:
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", store.Address, store.Database, store.User, store.Password))
	defer db.Close()
	if err != nil {
		return err
	}

	_, err = db.Exec("delete from configitem where application=? and name=? and machine=?", configItem.Application, configItem.Name, configItem.Machine)
	if err != nil {
		return err
	}

	return nil
}

func GetMSsqlCreateDDL() []byte {
	return dbCreateMSSQL
}
