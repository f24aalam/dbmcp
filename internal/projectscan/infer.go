package projectscan

import "strings"

// inferDialectFromEnvKeys picks postgres vs mysql when DB_CONNECTION / driver is absent,
// based on common libpq (PG*) or MYSQL_* variable families.
func inferDialectFromEnvKeys(values map[string]string) string {
	for _, k := range []string{"PGHOST", "PGPORT", "PGDATABASE", "PGUSER", "PGPASSWORD"} {
		if strings.TrimSpace(values[k]) != "" {
			return "postgres"
		}
	}

	for _, k := range []string{"MYSQL_HOST", "MYSQL_PORT", "MYSQL_DATABASE", "MYSQL_USER", "MYSQL_PASSWORD"} {
		if strings.TrimSpace(values[k]) != "" {
			return "mysql"
		}
	}

	return ""
}
