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
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns 
		WHERE table_schema = 'public' AND table_name = $1
	`, tableName)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var columns []Column

	for rows.Next() {
		var name, colType, isNullable string
		err := rows.Scan(&name, &colType, &isNullable)
		if err != nil {
			return nil, "", err
		}

		columns = append(columns, Column{
			Name:     name,
			Type:     colType,
			Nullable: isNullable == "YES",
		})
	}

	var primaryKey string
	err = db.QueryRowContext(ctx, `
		SELECT kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		WHERE tc.constraint_type = 'PRIMARY KEY'
			AND tc.table_schema = 'public'
			AND tc.table_name = $1
	`, tableName).Scan(&primaryKey)

	if primaryKey != "" {
		for i := range columns {
			if columns[i].Name == primaryKey {
				columns[i].Key = "PRI"
				break
			}
		}
	}

	return columns, primaryKey, nil
}
