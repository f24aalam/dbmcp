package repositories

import (
	"context"
	"database/sql"
)

type Column struct {
	Name     string
	Type     string
	Nullable bool
	Key      string
}

type TableRepository interface {
	GetDatabaseName(ctx context.Context, db *sql.DB) (string, error)
	GetDatabaseVersion(ctx context.Context, db *sql.DB) (string, error)

	GetTables(ctx context.Context, db *sql.DB) ([]string, error)
}
