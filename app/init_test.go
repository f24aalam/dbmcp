package app

import "testing"

func TestBuildConnectionURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   initInput
		want string
	}{
		{
			name: "postgres",
			in: initInput{
				DBType:   "postgres",
				Host:     "localhost",
				Port:     "5432",
				Database: "appdb",
				User:     "user",
				Password: "pass",
			},
			want: "postgres://user:pass@localhost:5432/appdb?sslmode=disable",
		},
		{
			name: "mysql",
			in: initInput{
				DBType:   "mysql",
				Host:     "127.0.0.1",
				Port:     "3306",
				Database: "appdb",
				User:     "root",
				Password: "pass",
			},
			want: "root:pass@tcp(127.0.0.1:3306)/appdb",
		},
		{
			name: "sqlite",
			in: initInput{
				DBType:     "sqlite",
				SQLitePath: "test.db",
			},
			want: "test.db",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := buildConnectionURL(tc.in)
			if err != nil {
				t.Fatalf("buildConnectionURL() error = %v", err)
			}
			if got != tc.want {
				t.Fatalf("buildConnectionURL() = %q, want %q", got, tc.want)
			}
		})
	}
}
