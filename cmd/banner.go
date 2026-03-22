/*
Copyright © 2026 Faizan Aalam (f24aalam@gmail.com)
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

// ASCII wordmark; backticks in art use double-quoted segment.
const dbmcpBanner = "       ____                       \n" +
	"  ____/ / /_  ____ ___  _________ \n" +
	" / __  / __ \\/ __ `__ \\/ ___/ __ \\\n" +
	"/ /_/ / /_/ / / / / / / /__/ /_/ /\n" +
	"\\__,_/_.___/_/ /_/ /_/\\___/ .___/ \n" +
	"                         /_/      "

const dbmcpTagline = "Database MCP server and CLI"

func printDbmcpBanner(out io.Writer) {
	bannerStyle := lipgloss.NewStyle().
		Foreground(ThemeGreen).
		Bold(true)
	tagStyle := lipgloss.NewStyle().
		Foreground(ThemeGreenMuted).
		MarginLeft(4)

	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, bannerStyle.Render(dbmcpBanner))
	_, _ = fmt.Fprintln(out, tagStyle.Render(dbmcpTagline))
	_, _ = fmt.Fprintln(out)
}
