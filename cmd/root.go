/*
Copyright © 2026 Faizan Aalam (f24aalam@gmail.com)
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godbmcp",
	Short: "Database MCP server and CLI",
	Long:  "godbmcp manages database connections and runs an MCP server for databases.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//
}
