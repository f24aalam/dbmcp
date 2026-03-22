package projectscan

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const maxStructuredFileSize = 512 * 1024

func scanStructuredFiles(root string) ([]Candidate, error) {
	paths, err := listStructuredConfigFiles(root)
	if err != nil {
		return nil, err
	}

	var candidates []Candidate
	seen := map[string]struct{}{}
	for _, p := range paths {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}

		info, err := os.Stat(p)
		if err != nil || info.IsDir() || info.Size() > maxStructuredFileSize {
			continue
		}

		fileCandidates, err := parseCandidateFile(root, p)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, fileCandidates...)
	}

	return dedupCandidates(candidates), nil
}

func listStructuredConfigFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if shouldSkipScanDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil || info.Size() > maxStructuredFileSize {
			return nil
		}

		base := filepath.Base(path)
		ext := strings.ToLower(filepath.Ext(path))

		switch {
		case strings.HasPrefix(base, ".env"):
			files = append(files, path)
		case strings.HasPrefix(strings.ToLower(base), "application.") &&
			(ext == ".yml" || ext == ".yaml" || ext == ".properties"):
			files = append(files, path)
		case ext == ".json" || ext == ".toml" || ext == ".yaml" || ext == ".yml" || ext == ".ini":
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func shouldSkipScanDir(name string) bool {
	switch name {
	case ".git", "node_modules", "vendor", "dist", "build", ".next", "target":
		return true
	default:
		return false
	}
}

func parserForPath(absPath string) string {
	base := strings.ToLower(filepath.Base(absPath))
	ext := strings.ToLower(filepath.Ext(absPath))
	switch {
	case strings.HasPrefix(base, ".env"):
		return "env"
	case strings.HasSuffix(base, ".properties"):
		return "properties"
	case ext == ".yaml" || ext == ".yml":
		return "yaml"
	case ext == ".json":
		return "json"
	case ext == ".toml":
		return "toml"
	case ext == ".ini":
		return "ini"
	default:
		return "text"
	}
}

func parseCandidateFile(root, absPath string) ([]Candidate, error) {
	rel := relOrAbs(root, absPath)
	parser := parserForPath(absPath)

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	text := string(data)

	merged := parseLooseKeyValues(text)
	if strings.EqualFold(filepath.Ext(absPath), ".json") {
		if jm := flattenJSONStringMap(data); len(jm) > 0 {
			if merged == nil {
				merged = jm
			} else {
				for k, v := range jm {
					if _, ok := merged[k]; !ok {
						merged[k] = v
					}
				}
			}
		}
	}

	if url, key := firstNonEmptyGroup(merged, urlKeys); url != "" {
		if c, ok := candidateFromRawURL(url, rel); ok {
			finalizeRawURLCandidate(&c, rel, parser, []string{"found " + key})
			return []Candidate{c}, nil
		}
	}

	if strings.Contains(text, "://") {
		if c, ok := candidateFromURLInText(text, rel, parser); ok {
			return []Candidate{c}, nil
		}
	}

	c, ok := candidateFromMap(merged, rel, parser)
	if !ok {
		return nil, nil
	}
	return []Candidate{c}, nil
}

func parseLooseKeyValues(text string) map[string]string {
	out := map[string]string{}
	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			k := normalizeKey(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			out[k] = v
			continue
		}

		if strings.Contains(line, ":") && !strings.HasPrefix(line, "-") {
			parts := strings.SplitN(line, ":", 2)
			k := normalizeKey(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			out[k] = v
		}
	}
	return out
}

func candidateFromMap(values map[string]string, relPath, parser string) (Candidate, bool) {
	if url, key := firstNonEmptyGroup(values, urlKeys); url != "" {
		if c, ok := candidateFromRawURL(url, relPath); ok {
			finalizeRawURLCandidate(&c, relPath, parser, []string{"found " + key})
			return c, true
		}
	}

	dbType := strings.ToLower(func() string {
		v, _ := firstNonEmptyGroup(values, driverKeys)
		return v
	}())
	if dbType == "pgsql" {
		dbType = "postgres"
	}
	if dbType == "" {
		dbType = inferDialectFromEnvKeys(values)
	}
	if dbType == "" {
		if dsn, _ := firstNonEmptyGroup(values, []string{"DB_DSN"}); strings.Contains(strings.ToLower(dsn), "sqlite") {
			dbType = "sqlite"
		}
	}

	if dbType == "sqlite" {
		pathVal, key := firstNonEmptyGroup(values, sqlitePathKeys)
		if pathVal == "" {
			return Candidate{}, false
		}
		c := Candidate{
			DBType:     "sqlite",
			SQLitePath: pathVal,
			Database:   pathVal,
			SourceFile: relPath,
			Parser:     parser,
		}
		scoreSQLiteCandidate(&c, relPath, []string{"found " + key}, values)
		return c, true
	}

	if dbType == "" {
		dbType = "mysql"
	}

	host, hk := firstNonEmptyGroup(values, hostKeys)
	portVal, portKey := firstNonEmptyGroup(values, portKeys)
	name, nk := firstNonEmptyGroup(values, databaseKeys)
	user, uk := firstNonEmptyGroup(values, userKeys)
	pass, passK := firstNonEmptyGroup(values, passwordKeys)

	if host == "" && name == "" && user == "" {
		return Candidate{}, false
	}

	port := portVal
	if port == "" {
		if dbType == "postgres" {
			port = "5432"
		} else {
			port = "3306"
		}
	}

	var evidence []string
	for _, x := range []struct {
		envKey string
		label  string
	}{{hk, "host"}, {portKey, "port"}, {nk, "database"}, {uk, "user"}, {passK, "password"}} {
		if x.envKey != "" {
			evidence = append(evidence, "found "+x.label+" ("+x.envKey+")")
		}
	}

	complete := host != "" && name != "" && user != "" && portKey != ""

	c := Candidate{
		DBType:     dbType,
		Host:       host,
		Port:       port,
		Database:   name,
		User:       user,
		Password:   pass,
		SourceFile: relPath,
		Parser:     parser,
	}
	scoreTupleCandidate(&c, relPath, complete, evidence, values)
	return c, true
}

func candidateFromURLInText(text, relPath, parser string) (Candidate, bool) {
	for _, key := range urlKeys {
		prefix := key + "="
		idx := strings.Index(text, prefix)
		if idx < 0 {
			continue
		}
		rest := text[idx+len(prefix):]
		line := rest
		if n := strings.IndexByte(rest, '\n'); n >= 0 {
			line = rest[:n]
		}
		line = strings.Trim(strings.TrimSpace(line), `"'`)
		if c, ok := candidateFromRawURL(line, relPath); ok {
			finalizeRawURLCandidate(&c, relPath, parser, []string{"found " + key + " in file"})
			return c, true
		}
	}

	for _, token := range strings.Fields(text) {
		if c, ok := candidateFromRawURL(strings.Trim(token, `"'`), relPath); ok {
			finalizeRawURLCandidate(&c, relPath, parser, []string{"parsed URL token in file"})
			return c, true
		}
	}

	return Candidate{}, false
}

func candidateFromRawURL(raw, source string) (Candidate, bool) {
	raw = strings.TrimSpace(raw)
	lower := strings.ToLower(raw)
	switch {
	case strings.HasPrefix(lower, "postgres://"), strings.HasPrefix(lower, "postgresql://"):
		return Candidate{DBType: "postgres", URL: raw, SourceFile: source}, true
	case strings.HasPrefix(lower, "mysql://"):
		return Candidate{DBType: "mysql", URL: raw, SourceFile: source}, true
	case strings.HasPrefix(lower, "jdbc:postgresql://"):
		clean := strings.TrimPrefix(raw, "jdbc:")
		return Candidate{DBType: "postgres", URL: clean, SourceFile: source}, true
	case strings.HasPrefix(lower, "jdbc:mysql://"):
		withoutJDBC := strings.TrimPrefix(raw, "jdbc:")
		withoutScheme := strings.TrimPrefix(withoutJDBC, "mysql://")
		return Candidate{DBType: "mysql", URL: "mysql://" + withoutScheme, SourceFile: source}, true
	case strings.HasSuffix(lower, ".sqlite"), strings.HasSuffix(lower, ".sqlite3"), strings.HasSuffix(lower, ".db"):
		return Candidate{DBType: "sqlite", SQLitePath: raw, Database: raw, SourceFile: source}, true
	default:
		return Candidate{}, false
	}
}

func dedupCandidates(in []Candidate) []Candidate {
	type key struct {
		DBType   string
		URL      string
		Host     string
		Port     string
		Database string
		User     string
		SQLite   string
	}
	seen := map[key]Candidate{}
	for _, c := range in {
		k := key{
			DBType:   c.DBType,
			URL:      c.URL,
			Host:     c.Host,
			Port:     c.Port,
			Database: c.Database,
			User:     c.User,
			SQLite:   c.SQLitePath,
		}
		prev, ok := seen[k]
		if !ok || c.Confidence > prev.Confidence {
			seen[k] = c
		}
	}

	out := make([]Candidate, 0, len(seen))
	for _, c := range seen {
		out = append(out, c)
	}
	return out
}

func normalizeKey(k string) string {
	k = strings.TrimSpace(k)
	k = strings.Trim(k, `"`)
	k = strings.ReplaceAll(k, ".", "_")
	k = strings.ReplaceAll(k, "-", "_")
	return strings.ToUpper(k)
}

func relOrAbs(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return rel
}
