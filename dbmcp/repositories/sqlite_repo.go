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

func (r *sqliteRepository) DescribeTable(ctx context.Context, db *sql.DB, tableName string) ([]Column, string, error) {
	rows, err := db.QueryContext(
		ctx,
		"PRAGMA table_info(" + tableName + ")",
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var columns []Column
	var primaryKey string

	for rows.Next() {
		var (
			cid        int
			name       string
			colType    string
			notNull    int
			defaultVal sql.NullString
			pk         int
		)

		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultVal, &pk); err != nil {
			return nil, "", err
		}

		col := Column{
			Name:     name,
			Type:     colType,
			Nullable: notNull == 0,
		}

		if pk == 1 {
			col.Key = "PRI"
			primaryKey = name
		}

		columns = append(columns, col)
	}

	return columns, primaryKey, nil
}
