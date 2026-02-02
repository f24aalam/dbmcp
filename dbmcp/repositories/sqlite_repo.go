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

func (r *sqliteRepository) GetTables(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}

		tables = append(tables, tableName)
	}

	return tables, nil
}
