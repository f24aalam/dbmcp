package dbmcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/f24aalam/godbmcp/dbmcp/repositories"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConnectionInput struct {
	ConnectionID string `json:"connection_id,omitempty" jsonschema:"optional; database connection id. If omitted, the server default connection is used"`
}

type TableInput struct {
	ConnectionID string `json:"connection_id,omitempty" jsonschema:"optional; database connection id. If omitted, the server default connection is used"`
	TableName    string `json:"table_name" jsonschema:"the table name"`
}

type QueryInput struct {
	ConnectionID string `json:"connection_id,omitempty" jsonschema:"optional; database connection id. If omitted, the server default connection is used"`
	Query        string `json:"query" jsonschema:"the query to run"`
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
	conn, err := GetDB()
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}

	repo := repositories.GetRepository(GetDBType())

	dbName, err := repo.GetDatabaseName(ctx, conn.DB)
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}

	dbVersion, err := repo.GetDatabaseVersion(ctx, conn.DB)
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}

	return nil, GetDatabaseInfoOutput{
		DatabaseName:     dbName,
		DatabaseVendor:   GetDBType(),
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
	conn, err := GetDB()
	if err != nil {
		return nil, GetTablesOutput{}, err
	}

	repo := repositories.GetRepository(GetDBType())

	tables, err := repo.GetTables(ctx, conn.DB)

	return nil, GetTablesOutput{
		Tables:     tables,
		TableCount: len(tables),
	}, nil
}

type DescribeTableOutput struct {
	Columns    []repositories.Column `json:"columns"`
	PrimaryKey string                `json:"primary_key"`
}

func DescribeTable(ctx context.Context, req *mcp.CallToolRequest, input TableInput) (
	*mcp.CallToolResult,
	DescribeTableOutput,
	error,
) {
	conn, err := GetDB()
	if err != nil {
		return nil, DescribeTableOutput{}, err
	}

	repo := repositories.GetRepository(GetDBType())

	columns, primaryKey, err := repo.DescribeTable(ctx, conn.DB, input.TableName)

	if err != nil {
		return nil, DescribeTableOutput{}, err
	}

	return nil, DescribeTableOutput{
		Columns:    columns,
		PrimaryKey: primaryKey,
	}, nil
}

type SelectQueryOutput struct {
	Rows     []map[string]interface{} `json:"rows"`
	RowCount int                      `json:"row_count"`
}

func RunSelectQuery(ctx context.Context, req *mcp.CallToolRequest, input QueryInput) (
	*mcp.CallToolResult,
	SelectQueryOutput,
	error,
) {
	query := strings.TrimSpace(strings.ToUpper(input.Query))
	if !strings.HasPrefix(query, "SELECT") {
		return nil, SelectQueryOutput{}, fmt.Errorf("only SELECT queries are allowed")
	}

	conn, err := GetDB()
	if err != nil {
		return nil, SelectQueryOutput{}, err
	}

	repo := repositories.GetRepository(GetDBType())
	rows, rowCount, err := repo.RunSelectQuery(ctx, conn.DB, query)

	return nil, SelectQueryOutput{
		Rows:     rows,
		RowCount: rowCount,
	}, nil
}
