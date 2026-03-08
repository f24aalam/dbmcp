/*
Copyright © 2026 Faizan Aalam <f24aalam@gmail.com>
*/
package cmd

import (
	"github.com/f24aalam/dbmcp/app"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new connection",
	Long:  `Add a new connection to the dbmcp, this will generate a unique connection id which can be used with mcp server`,
	Run: func(cmd *cobra.Command, args []string) {
		app.AddNewConnection(nil)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
