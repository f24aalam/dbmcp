package app

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
)

// ormPostgresQueryKeys are URI query keys used by Prisma and similar tools that
// lib/pq forwards to PostgreSQL as startup options where they are invalid.
var ormPostgresQueryKeys = []string{
	"schema",
	"connection_limit",
	"pool_timeout",
}

// sanitizePostgresURLForPQ removes ORM-only query parameters from postgres /
// postgresql URLs so github.com/lib/pq can connect (e.g. Prisma ?schema=public).
func sanitizePostgresURLForPQ(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return raw, nil
	}

	q := u.Query()
	changed := false
	for _, k := range ormPostgresQueryKeys {
		if len(q[k]) > 0 {
			q.Del(k)
			changed = true
		}
	}

	if !changed {
		return raw, nil
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// canonicalMySQLDSNFromURL converts a mysql:// URI (common in env files) into a
// go-sql-driver/mysql DSN string.
func canonicalMySQLDSNFromURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	if !strings.EqualFold(u.Scheme, "mysql") {
		return "", fmt.Errorf("expected mysql URL scheme, got %q", u.Scheme)
	}

	cfg := mysqldriver.NewConfig()
	if u.User != nil {
		cfg.User = u.User.Username()
		cfg.Passwd, _ = u.User.Password()
	}

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "3306"
	}

	cfg.Net = "tcp"
	if host == "" {
		cfg.Addr = "127.0.0.1:" + port
	} else {
		cfg.Addr = net.JoinHostPort(host, port)
	}

	db := strings.TrimPrefix(u.Path, "/")
	if db == "" {
		return "", fmt.Errorf("mysql URL missing database name in path")
	}

	cfg.DBName = db

	if len(u.Query()) > 0 {
		cfg.Params = make(map[string]string)
		for k, vals := range u.Query() {
			if len(vals) == 0 {
				continue
			}

			cfg.Params[k] = vals[0]
		}
	}

	return cfg.FormatDSN(), nil
}
