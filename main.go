package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/f24aalam/godbmcp/database"
	"github.com/f24aalam/godbmcp/storage"
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
	var dbName string
	var dbType string
	var dbConnUrl string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Connection Name").
				Value(&dbName),
			huh.NewSelect[string]().
				Title("Select Database").
				Options(
					huh.NewOption("MySQL", "mysql"),
				).
				Value(&dbType),
			huh.NewInput().
				Title("Enter connection string").
				Value(&dbConnUrl),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println(err)
	}

	conn := database.Connection{
		Database: dbType,
		ConnectionUrl: dbConnUrl,
	}

	err = conn.Open()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	id, err := storage.SaveCredential(dbName, dbType, dbConnUrl)
	if err != nil {
		fmt.Println("Error in saving connection: ", err)
	}

	fmt.Println("Database connection success, saved with id: ", id)

	defer conn.Close()
}
