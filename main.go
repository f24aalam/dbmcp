package main

import (
	"fmt"
	"os"

	"github.com/f24aalam/godbmcp/cli"
	_ "github.com/go-sql-driver/mysql"
)

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
