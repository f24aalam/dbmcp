package repositories

import (
	"context"
	"database/sql"
)

type postgresRepository struct {
	BaseRepository
}

func NewPostgresRepository() TableRepository {
	return &postgresRepository{}
}

func (r *postgresRepository) GetDatabaseName(ctx context.Context, db *sql.DB) (string, error) {
	var name string
	err := db.QueryRowContext(ctx, "SELECT current_database()").Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *postgresRepository) GetDatabaseVersion(ctx context.Context, db *sql.DB) (string, error) {
	var version string
	err := db.QueryRowContext(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return "", err
	}
	return version, nil
}

func (r *postgresRepository) GetTables(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
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

func (r *postgresRepository) DescribeTable(ctx context.Context, db *sql.DB, tableName string) ([]Column, string, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT column_name, data_type, is_nullable, column_key
		FROM information_schema.columns 
		WHERE table_schema = 'public' AND table_name = $1
	`, tableName)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var columns []Column
	var primaryKey string

	for rows.Next() {
		var name, colType, isNullable, columnKey string
		err := rows.Scan(&name, &colType, &isNullable, &columnKey)
		if err != nil {
			return nil, "", err
		}

		columns = append(columns, Column{
			Name:     name,
			Type:     colType,
			Nullable: isNullable == "YES",
			Key:      columnKey,
		})

		if columnKey == "PRI" {
			primaryKey = name
		}
	}

	return columns, primaryKey, nil
}
