package mcpinstall

import (
	"os"
	"path/filepath"
)

type Target struct {
	ID         string
	Name       string
	Path       string // project-relative path
	GlobalPath string // absolute path for global install
	RootKey    string
	Format     string
	ServerType string
}

func SupportedTargets() []Target {
	return []Target{
		{
			ID: "claude", Name: "Claude Code",
			Path: ".mcp.json", GlobalPath: getClaudeGlobalPath(),
			RootKey: "mcpServers", Format: "json", ServerType: "stdio",
		},
		{
			ID: "cursor", Name: "Cursor",
			Path: ".cursor/mcp.json", GlobalPath: getCursorGlobalPath(),
			RootKey: "mcpServers", Format: "json", ServerType: "stdio",
		},
		{
			ID: "copilot-vscode", Name: "Copilot VSCode",
			Path: ".vscode/mcp.json", GlobalPath: "", // Platform-dependent, leaving empty for now
			RootKey: "servers", Format: "json", ServerType: "stdio",
		},
		{
			ID: "copilot-cli", Name: "Copilot CLI",
			Path: ".copilot/mcp-config.json", GlobalPath: getCopilotGlobalPath(),
			RootKey: "mcpServers", Format: "json", ServerType: "local",
		},
		{
			ID: "codex", Name: "Codex",
			Path: ".codex/config.toml", GlobalPath: getCodexGlobalPath(),
			RootKey: "mcp_servers", Format: "toml", ServerType: "stdio",
		},
		{
			ID: "opencode", Name: "OpenCode",
			Path: "opencode.json", GlobalPath: getOpenCodeGlobalPath(),
			RootKey: "mcp", Format: "json", ServerType: "local",
		},
	}
}

func GetTarget(id string) (Target, bool) {
	for _, t := range SupportedTargets() {
		if t.ID == id {
			return t, true
		}
	}

	return Target{}, false
}

func getClaudeGlobalPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Global Path: ~/.claude.json
	return filepath.Join(home, ".claude.json")
}

func getCursorGlobalPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Global Path: ~/.cursor/mcp.json
	return filepath.Join(home, ".cursor", "mcp.json")
}

func getCopilotGlobalPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Global Path: ~/.copilot/mcp-config.json
	return filepath.Join(home, ".copilot", "mcp-config.json")
}

func getCodexGlobalPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Global Path: ~/.codex/config.toml
	return filepath.Join(home, ".codex", "config.toml")
}

func getOpenCodeGlobalPath() string {
	config, err := os.UserConfigDir()
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, ".config", "opencode", "opencode.json")
	}
	// Global Path: ~/.config/opencode/opencode.json
	return filepath.Join(config, "opencode", "opencode.json")
}
