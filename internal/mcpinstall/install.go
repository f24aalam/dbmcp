package mcpinstall

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type InstallResult struct {
	TargetID string
	Path     string
	Written  bool
	Err      error
}

func InstallForAgents(root string, targetIDs []string, connectionID string, isGlobal bool) []InstallResult {
	out := make([]InstallResult, 0, len(targetIDs))
	for _, id := range targetIDs {
		t, ok := GetTarget(id)
		if !ok {
			out = append(out, InstallResult{TargetID: id, Err: os.ErrInvalid})
			continue
		}

		path := ""
		displayPath := ""
		if isGlobal {
			if t.GlobalPath == "" {
				out = append(out, InstallResult{TargetID: id, Err: fmt.Errorf("global install not supported for this agent")})
				continue
			}
			path = t.GlobalPath
			displayPath = t.GlobalPath
		} else {
			path = filepath.Join(root, t.Path)
			displayPath = t.Path
		}

		err := installSingle(path, t, connectionID)
		out = append(out, InstallResult{
			TargetID: id,
			Path:     displayPath,
			Written:  err == nil,
			Err:      err,
		})
	}

	return out
}

func installSingle(path string, target Target, connectionID string) error {
	if target.Format == "toml" {
		return installSingleTOML(path, connectionID)
	}

	root := map[string]any{}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &root)
	}

	servers, _ := root[target.RootKey].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}

	servers["dbmcp"] = buildJSONServerEntry(target, connectionID)
	root[target.RootKey] = servers

	if target.RootKey == "mcp" {
		if _, ok := root["$schema"]; !ok {
			root["$schema"] = "https://opencode.ai/config.json"
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func buildJSONServerEntry(target Target, connectionID string) map[string]any {
	entry := map[string]any{
		"type":    target.ServerType,
		"command": "dbmcp",
		"args":    []string{"mcp", "--connection-id", connectionID},
	}

	switch target.ID {
	case "copilot-cli":
		entry["tools"] = []string{"*"}
		entry["env"] = map[string]string{}
	case "opencode":
		entry["enabled"] = true
	}

	return entry
}

func installSingleTOML(path string, connectionID string) error {
	block := strings.TrimSpace(fmt.Sprintf(`
[mcp_servers.dbmcp]
command = "dbmcp"
args = ["mcp", "--connection-id", "%s"]
`, connectionID)) + "\n"

	existing := ""
	if b, err := os.ReadFile(path); err == nil {
		existing = string(b)
	} else if !os.IsNotExist(err) {
		return err
	}

	re := regexp.MustCompile(`(?ms)^\[mcp_servers\.dbmcp\]\n(?:[^\[]*\n?)*`)
	if re.MatchString(existing) {
		existing = re.ReplaceAllString(existing, block+"\n")
	} else {
		if strings.TrimSpace(existing) != "" && !strings.HasSuffix(existing, "\n") {
			existing += "\n"
		}

		existing += "\n" + block
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(existing), 0o644)
}
