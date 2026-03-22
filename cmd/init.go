/*
Copyright © 2026 Faizan Aalam <f24aalam@gmail.com>
*/
package cmd

import (
	"github.com/f24aalam/dbmcp/app"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dbmcp in current project",
	Long:  "Detect database configuration in this project, save a connection, and install project-local MCP config.",
	Run: func(cmd *cobra.Command, args []string) {
		// Banner on stderr so it shows before the stepflow TUI (stderr).
		printDbmcpBanner(cmd.ErrOrStderr())
		app.InitProject()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
