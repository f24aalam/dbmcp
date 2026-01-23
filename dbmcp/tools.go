package dbmcp

import (
	"context"
	"fmt"
	"strings"
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

	var dbName string
	err = conn.DB.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		return nil, GetDatabaseInfoOutput{}, err
	}

	var dbVersion string
	err = conn.DB.QueryRowContext(ctx, "SELECT VERSION()").Scan(&dbVersion)
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

	rows, err := conn.DB.QueryContext(ctx, "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE()")
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
	conn, err := GetDB()
	if err != nil {
		return nil, DescribeTableOutput{}, err
	}

	rows, err := conn.DB.QueryContext(ctx, "SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", input.TableName)
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
			return nil, DescribeTableOutput{}, err
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

	rows, err := conn.DB.QueryContext(ctx, strings.TrimSpace(input.Query))
	if err != nil {
		return nil, SelectQueryOutput{}, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, SelectQueryOutput{}, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuesPtr := make([]interface{}, len(columns))

		for i := range columns {
			valuesPtr[i] = &values[i]
		}

		err = rows.Scan(valuesPtr...)
		if err != nil {
			return nil, SelectQueryOutput{}, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}

			val := values[i]
			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}

			entry[col] = v
		}

		result = append(result, entry)
	}

	return nil, SelectQueryOutput{
		Rows:     result,
		RowCount: len(result),
	}, nil
}
