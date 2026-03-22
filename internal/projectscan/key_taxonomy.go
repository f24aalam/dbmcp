package projectscan

import "strings"

// Canonical env/config key groups for framework-agnostic discovery.

var urlKeys = []string{
	"DATABASE_URL",
	"DB_URL",
	"SPRING_DATASOURCE_URL",
	"JDBC_URL",
	"SQLALCHEMY_DATABASE_URI",
}

var driverKeys = []string{
	"DB_CONNECTION",
	"DB_DRIVER",
	"DATABASE_DRIVER",
	"SPRING_DATASOURCE_DRIVER_CLASS_NAME",
}

var hostKeys = []string{
	"DB_HOST",
	"DATABASE_HOST",
	"SPRING_DATASOURCE_HOST",
	"PGHOST",
	"MYSQL_HOST",
}

var portKeys = []string{
	"DB_PORT",
	"DATABASE_PORT",
	"SPRING_DATASOURCE_PORT",
	"PGPORT",
	"MYSQL_PORT",
}

var databaseKeys = []string{
	"DB_DATABASE",
	"DATABASE_NAME",
	"DB_NAME",
	"SPRING_DATASOURCE_NAME",
	"PGDATABASE",
	"MYSQL_DATABASE",
}

var userKeys = []string{
	"DB_USERNAME",
	"DATABASE_USER",
	"DB_USER",
	"SPRING_DATASOURCE_USERNAME",
	"PGUSER",
	"MYSQL_USER",
}

var passwordKeys = []string{
	"DB_PASSWORD",
	"DATABASE_PASSWORD",
	"SPRING_DATASOURCE_PASSWORD",
	"PGPASSWORD",
	"MYSQL_PASSWORD",
}

var sqlitePathKeys = []string{
	"DB_DATABASE",
	"SQLITE_PATH",
	"DATABASE_FILE",
	"DB_DSN",
}

// allKeysForGrep returns every normalized key we look for in fallback line scanning.
func allKeysForGrep() map[string]struct{} {
	out := map[string]struct{}{}
	add := func(keys []string) {
		for _, k := range keys {
			out[k] = struct{}{}
		}
	}
	add(urlKeys)
	add(driverKeys)
	add(hostKeys)
	add(portKeys)
	add(databaseKeys)
	add(userKeys)
	add(passwordKeys)
	add(sqlitePathKeys)
	return out
}

func firstNonEmptyGroup(m map[string]string, keys []string) (value string, matchedKey string) {
	for _, k := range keys {
		if v := strings.TrimSpace(m[k]); v != "" {
			return v, k
		}
	}
	return "", ""
}
