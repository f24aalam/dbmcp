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
