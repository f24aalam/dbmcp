package dbmcp

import (
	"sync"

	"github.com/f24aalam/godbmcp/database"
)

var (
	dbOnce sync.Once
	dbConn *database.Connection
	dbErr  error

	dbType string
	dbUrl  string
)

func InitDBConfig(t, url string) {
	dbType = t
	dbUrl = url
}

func GetDBType() string {
	return dbType
}

func GetDB() (*database.Connection, error) {
	dbOnce.Do(func() {
		dbConn = &database.Connection{
			Database:      dbType,
			ConnectionURL: dbUrl,
		}

		dbErr = dbConn.Open()
	})

	return dbConn, dbErr
}

func CloseDB() error  {
	if dbConn != nil {
		return dbConn.Close()
	}

	return nil
}
