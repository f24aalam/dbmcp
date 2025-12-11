package database

import (
	"database/sql"
	"fmt"
)

type Connection struct {
	Database string
	ConnectionUrl string
	DB *sql.DB
}

func (c *Connection) Open() error {
	db, err := sql.Open(c.Database, c.ConnectionUrl)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	c.DB = db
	return nil;
}

func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close();
	}

	return nil
}
