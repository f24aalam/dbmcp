package projectscan

import (
	"os"
	"strings"
	"testing"
)

func TestScanProjectFixtures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		root     string
		expectDB string
	}{
		{name: "laravel", root: "testdata/laravel", expectDB: "mysql"},
		{name: "symfony", root: "testdata/symfony", expectDB: "postgres"},
		{name: "nest", root: "testdata/nest", expectDB: "postgres"},
		{name: "spring", root: "testdata/spring", expectDB: "postgres"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := ScanProject(tc.root)
			if err != nil {
				t.Fatalf("ScanProject() error = %v", err)
			}

			if len(res.Candidates) == 0 {
				t.Fatalf("expected at least one candidate")
			}

			if got := res.Candidates[0].DBType; got != tc.expectDB {
				t.Fatalf("db type mismatch: got %q want %q", got, tc.expectDB)
			}
		})
	}
}

func TestScanProjectGenericPGEnv(t *testing.T) {
	t.Parallel()
	res, err := ScanProject("testdata/generic")
	if err != nil {
		t.Fatalf("ScanProject() error = %v", err)
	}

	if len(res.Candidates) == 0 {
		t.Fatalf("expected at least one candidate")
	}

	top := res.Candidates[0]
	if top.DBType != "postgres" {
		t.Fatalf("expected postgres, got %q", top.DBType)
	}

	if top.Host != "db.internal" || top.Database != "analytics" {
		t.Fatalf("unexpected connection fields: %+v", top)
	}

	if top.Parser != "env" {
		t.Fatalf("expected env parser, got %q", top.Parser)
	}
}

func TestPlaceholderPenaltyOnURL(t *testing.T) {
	t.Parallel()
	res, err := ScanProject("testdata/placeholder")
	if err != nil {
		t.Fatalf("ScanProject() error = %v", err)
	}

	if len(res.Candidates) == 0 {
		t.Fatalf("expected at least one candidate")
	}

	top := res.Candidates[0]
	if top.DBType != "postgres" {
		t.Fatalf("expected postgres, got %q", top.DBType)
	}

	if top.Confidence >= scoreDSNRecognized+scoreConfigPathBonus {
		t.Fatalf("expected placeholder penalty to reduce confidence, got %d", top.Confidence)
	}

	found := false
	for _, e := range top.Evidence {
		if strings.Contains(strings.ToLower(e), "placeholder") {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected placeholder mention in evidence: %v", top.Evidence)
	}
}

func TestFallbackScanProject(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	content := "DB_HOST=localhost\nDB_PORT=5432\nDB_DATABASE=mydb\nDB_USERNAME=me\nDB_PASSWORD=topsecret\nDB_CONNECTION=pgsql\n"
	if err := osWriteFile(root+"/config.txt", []byte(content)); err != nil {
		t.Fatal(err)
	}

	candidates, err := fallbackScanProject(root)
	if err != nil {
		t.Fatalf("fallback scan error: %v", err)
	}

	if len(candidates) == 0 {
		t.Fatalf("expected candidate from fallback scan")
	}

	if candidates[0].DBType != "postgres" {
		t.Fatalf("expected postgres, got %q", candidates[0].DBType)
	}
}

func osWriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}
