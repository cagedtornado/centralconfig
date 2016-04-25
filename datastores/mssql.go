package datastores

import (
	"fmt"

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

func GetMSsqlCreateDDL() []byte {
	return dbCreateMSSQL
}
