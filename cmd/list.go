/*
Copyright © 2026 Faizan Aalam <f24aalam@gmail.com>
*/
package cmd

import (
	"github.com/f24aalam/godbmcp/app"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the connections",
	Long:  `List all the connections stored with label and connection id`,
	Run: func(cmd *cobra.Command, args []string) {
		app.ListAllConnections()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
