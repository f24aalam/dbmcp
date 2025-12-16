package dbmcp

import (
	"context"

	"github.com/f24aalam/godbmcp/database"
	"github.com/f24aalam/godbmcp/storage"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConnectionInput struct {
	ConnectionID string `json:"connection_id" jsonschema:"the connection id to connect with the database"`
}

type TableInput struct {
	ConnectionID string `json:"connection_id" jsonschema:"the connection id to connect with the database"`
	TableName    string `json:"table_name" jsonschema:"the table name"`
}

type GetDatabaseInfoOutput struct {
	DatabaseName     string `json:"database_name"`
	DatabaseVendor   string `json:"database_vendor"`
	DatabaseVersion  string `json:"database_version"`
	ConnectionStatus string `json:"connection_status"`
}

func GetDatabaseInfo(ctx context.Context, req *mcp.CallToolRequest, input ConnectionInput) (
	*mcp.CallToolResult,
	GetDatabaseInfoOutput,
	error,
) {
	dbType, dbUrl, err := storage.GetCredentialById(input.ConnectionID)
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}

	conn := &database.Connection{
		Database:      dbType,
		ConnectionUrl: dbUrl,
	}

	err = conn.Open()
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}
	defer conn.Close()

	var dbName string
	err = conn.DB.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}

	var dbVersion string
	err = conn.DB.QueryRow("SELECT VERSION()").Scan(&dbVersion)
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}

	return nil, GetDatabaseInfoOutput{
		DatabaseName:     dbName,
		DatabaseVendor:   dbType,
		DatabaseVersion:  dbVersion,
		ConnectionStatus: "connected",
	}, nil
}

type GetTablesOutput struct {
	Tables     []string `json:"tables"`
	TableCount int      `json:"table_count"`
}

func GetTables(ctx context.Context, req *mcp.CallToolRequest, input ConnectionInput) (
	*mcp.CallToolResult,
	GetTablesOutput,
	error,
) {
	dbType, dbUrl, err := storage.GetCredentialById(input.ConnectionID)
	if err != nil {
		return nil, GetTablesOutput{}, err
	}

	conn := &database.Connection{
		Database:      dbType,
		ConnectionUrl: dbUrl,
	}

	err = conn.Open()
	if err != nil {
		return nil, GetTablesOutput{}, err
	}
	defer conn.Close()

	rows, err := conn.DB.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE()")
	if err != nil {
		return nil, GetTablesOutput{}, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, GetTablesOutput{}, err
		}

		tables = append(tables, tableName)
	}

	return nil, GetTablesOutput{
		Tables:     tables,
		TableCount: len(tables),
	}, nil
}

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Key      string `json:"key"`
}

type DescribeTableOutput struct {
	Columns    []Column `json:"columns"`
	PrimaryKey string   `json:"primary_key"`
}

func DescribeTable(ctx context.Context, req *mcp.CallToolRequest, input TableInput) (
	*mcp.CallToolResult,
	DescribeTableOutput,
	error,
) {
	dbType, dbUrl, err := storage.GetCredentialById(input.ConnectionID)
	if err != nil {
		return nil, DescribeTableOutput{}, err
	}

	conn := database.Connection{
		Database:      dbType,
		ConnectionUrl: dbUrl,
	}

	err = conn.Open()
	if err != nil {
		return nil, DescribeTableOutput{}, err
	}
	defer conn.Close()

	rows, err := conn.DB.Query("SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", input.TableName)
	if err != nil {
		return nil, DescribeTableOutput{}, err
	}
	defer rows.Close()

	var columns []Column
	var primaryKey string

	for rows.Next() {
		var name, colType, isNullable, columnKey string
		err := rows.Scan(&name, &colType, &isNullable, &columnKey)
		if err != nil {
			return nil, DescribeTableOutput{}, nil
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

	return nil, DescribeTableOutput{
		Columns:    columns,
		PrimaryKey: primaryKey,
	}, nil
}
