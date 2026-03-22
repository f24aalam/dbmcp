package mcpinstall

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallForAgentsMergesExistingConfig(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, ".cursor", "mcp.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	existing := `{"mcpServers":{"other":{"command":"echo","args":["ok"]}}}`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	results := InstallForAgents(root, []string{"cursor"}, "conn-123")
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("unexpected results: %+v", results)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	s := string(b)
	if !strings.Contains(s, `"other"`) {
		t.Fatalf("expected existing server to be preserved: %s", s)
	}

	if !strings.Contains(s, `"dbmcp"`) || !strings.Contains(s, `"conn-123"`) {
		t.Fatalf("expected dbmcp server with connection id: %s", s)
	}
}

func TestInstallForAgentsWritesAgentSpecificSchemas(t *testing.T) {
	t.Parallel()
	root := t.TempDir()

	results := InstallForAgents(root, []string{"claude", "copilot-vscode", "copilot-cli", "opencode"}, "conn-xyz")
	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("target %s failed: %v", r.TargetID, r.Err)
		}
	}

	assertContains := func(path, want string) {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(b), want) {
			t.Fatalf("%s missing %q:\n%s", path, want, string(b))
		}
	}

	assertContains(".mcp.json", `"mcpServers"`)
	assertContains(".mcp.json", `"type": "stdio"`)
	assertContains(".vscode/mcp.json", `"servers"`)
	assertContains(".copilot/mcp-config.json", `"tools"`)
	assertContains("opencode.json", `"mcp"`)
	assertContains("opencode.json", `"enabled": true`)
}

func TestInstallForAgentsMergesCodexTOML(t *testing.T) {
	t.Parallel()
	root := t.TempDir()

	path := filepath.Join(root, ".codex", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	initial := "model = \"gpt-5\"\n[mcp_servers.other]\ncommand = \"echo\"\n"
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	results := InstallForAgents(root, []string{"codex"}, "mysql-1")
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("unexpected codex install result: %+v", results)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	s := string(b)
	if !strings.Contains(s, "[mcp_servers.other]") {
		t.Fatalf("expected existing block preserved:\n%s", s)
	}

	if !strings.Contains(s, "[mcp_servers.dbmcp]") || !strings.Contains(s, "mysql-1") {
		t.Fatalf("expected dbmcp toml block:\n%s", s)
	}
}
