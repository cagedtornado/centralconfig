package datastores

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//	The SQLite database information
type SQLiteDB struct {
	Database string
	User     string
	Password string
}

func (store SQLiteDB) Get(configItem ConfigItem) (ConfigItem, error) {

	//	Open the database
	db, err := sql.Open("sqlite3", store.Database)
	checkErr(err)

	//	Insert or update the item:
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	//	Close the database
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
