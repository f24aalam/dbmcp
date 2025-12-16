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
