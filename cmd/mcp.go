/*
Copyright © 2026 Faizan Aalam <f24aalam@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/f24aalam/godbmcp/app"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the mcp server",
	Long: `Start a MCP server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.StartServer(connectionID)
	},
}

var connectionID string

func init() {
	rootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().StringVar(
		&connectionID,
		"connection-id",
		"",
		"Database connection ID to use for MCP server",
	)
}
