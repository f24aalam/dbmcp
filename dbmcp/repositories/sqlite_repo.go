package repositories

import (
	"context"
	"database/sql"
)

type sqliteRepository struct{}

func NewSQLiteRepository() TableRepository {
	return &sqliteRepository{}
}

func (r *sqliteRepository) GetDatabaseName(ctx context.Context, db *sql.DB) (string, error) {
	var (
		seq  int
		name string
		file string
	)

	err := db.QueryRowContext(ctx, "PRAGMA database_list").Scan(&seq, &name, &file)

	if err != nil {
		return "", err
	}

	if file == "" {
		return ":memory:", nil
	}

	return "File: " + file + " | Name: " + name, nil
}

func (r *sqliteRepository) GetDatabaseVersion(ctx context.Context, db *sql.DB) (string, error) {
	var version string

	err := db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)

	if err != nil {
		return "", err
	}

	return version, nil
}
