package mcpinstall

type Target struct {
	ID         string
	Name       string
	Path       string
	RootKey    string
	Format     string
	ServerType string
}

func SupportedTargets() []Target {
	return []Target{
		{ID: "claude", Name: "Claude Code", Path: ".mcp.json", RootKey: "mcpServers", Format: "json", ServerType: "stdio"},
		{ID: "cursor", Name: "Cursor", Path: ".cursor/mcp.json", RootKey: "mcpServers", Format: "json", ServerType: "stdio"},
		{ID: "copilot-vscode", Name: "Copilot VSCode", Path: ".vscode/mcp.json", RootKey: "servers", Format: "json", ServerType: "stdio"},
		{ID: "copilot-cli", Name: "Copilot CLI", Path: ".copilot/mcp-config.json", RootKey: "mcpServers", Format: "json", ServerType: "local"},
		{ID: "codex", Name: "Codex", Path: ".codex/config.toml", RootKey: "mcp_servers", Format: "toml", ServerType: "stdio"},
		{ID: "opencode", Name: "OpenCode", Path: "opencode.json", RootKey: "mcp", Format: "json", ServerType: "local"},
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
