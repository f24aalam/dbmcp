package app

import (
	"context"

	"github.com/f24aalam/godbmcp/dbmcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func StartServer() error {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_database_info",
		Description: "Get database information like name, version and connection status",
	}, dbmcp.GetDatabaseInfo)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_tables",
		Description: "Get list of all tables in the database",
	}, dbmcp.GetTables)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "describe_table",
		Description: "Get detailed information about a table including columns, types and primary key",
	}, dbmcp.DescribeTable)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "run_select_query",
		Description: "Run a SELECT query to retrieve data from the database",
	}, dbmcp.RunSelectQuery)

	return server.Run(context.Background(), &mcp.StdioTransport{})
}
