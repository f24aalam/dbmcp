package main

import (
	"context"
	"fmt"
	"os"

	"github.com/f24aalam/godbmcp/cli"
	_ "github.com/go-sql-driver/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Connection struct {
	Database string
	ConectionUrl string
}

type Input struct {
	Name string `json:"name" jsonschema:"the name of a the person to greet"`
}

type Output struct {
	Greetings string `json:"greeting" jsonschema:"the greeting to tell user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{Greetings: "Hi from the godbmcp tool to " + input.Name}, nil
}

func startServer() {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Println(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dbmcp <command>")
		fmt.Println("Commands: add-new, list, mcp")

		return
	}

	switch os.Args[1] {
	case "add-new":
		cli.AddNewConnection(nil)
	case "list":
		cli.ListAllConnections()
	case "mcp":
		cli.StartServer()
	}
}
