package app

import (
	"strings"
	"testing"
)

func TestSanitizePostgresURLForPQ(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "prisma_schema",
			in:   "postgresql://vision_user:vision_pass@localhost:5432/vision?schema=public",
			out:  "postgresql://vision_user:vision_pass@localhost:5432/vision",
		},
		{
			name: "postgres_scheme",
			in:   "postgres://u:p@host:5432/db?schema=app&sslmode=disable",
			out:  "postgres://u:p@host:5432/db?sslmode=disable",
		},
		{
			name: "unchanged_no_orm_keys",
			in:   "postgres://u:p@localhost:5432/mydb?sslmode=require",
			out:  "postgres://u:p@localhost:5432/mydb?sslmode=require",
		},
		{
			name: "connection_limit_stripped",
			in:   "postgres://u:p@localhost/db?connection_limit=10&schema=public",
			out:  "postgres://u:p@localhost/db",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := sanitizePostgresURLForPQ(tc.in)
			if err != nil {
				t.Fatalf("sanitizePostgresURLForPQ: %v", err)
			}

			if got != tc.out {
				t.Fatalf("got %q, want %q", got, tc.out)
			}
		})
	}
}

func TestCanonicalMySQLDSNFromURL(t *testing.T) {
	t.Parallel()

	raw := "mysql://root:secret@127.0.0.1:3306/appdb"
	got, err := canonicalMySQLDSNFromURL(raw)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(got, "root:secret@tcp(127.0.0.1:3306)/appdb") {
		t.Fatalf("unexpected DSN: %q", got)
	}
}
