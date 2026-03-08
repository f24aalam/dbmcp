/*
Copyright © 2026 Faizan Aalam <f24aalam@gmail.com>
*/
package main

import (
	"github.com/f24aalam/dbmcp/cmd"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func main() {
	cmd.Execute()
}
