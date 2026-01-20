/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/f24aalam/godbmcp/cli"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the mcp server",
	Long: `Start a MCP server`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
