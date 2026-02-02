package repositories

import (
	"context"
	"database/sql"
)

type mysqlRepository struct{
	BaseRepository
}

func NewMySQLRepository() TableRepository {
	return &mysqlRepository{}
}

func (r *mysqlRepository) GetDatabaseName(ctx context.Context, db *sql.DB) (string, error) {
	var name string

	err := db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&name)

	if err != nil {
		return "", err
	}

	return name, nil
}

func (r *mysqlRepository) GetDatabaseVersion(ctx context.Context, db *sql.DB) (string, error) {
	var version string

	err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version)

	if err != nil {
		return "", err
	}

	return version, nil
}

func (r *mysqlRepository) GetTables(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE()")
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

func (r *mysqlRepository) DescribeTable(ctx context.Context, db *sql.DB, tableName string) ([]Column, string, error) {
	rows, err := db.QueryContext(ctx, "SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", tableName)
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
