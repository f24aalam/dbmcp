package cli

import (
	"context"
	"fmt"

	"github.com/f24aalam/godbmcp/dbmcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func StartServer() {
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

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Println(err)
	}
}
