package repositories

import (
	"context"
	"database/sql"
)

type mysqlRepository struct{}

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
