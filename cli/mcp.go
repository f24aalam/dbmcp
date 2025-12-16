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

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Println(err)
	}
}
