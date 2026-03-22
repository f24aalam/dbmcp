package projectscan

import (
	"path/filepath"
	"strings"
)

const (
	scoreDSNRecognized      = 90
	scoreCompleteTuple      = 70
	scorePartialTuple       = 55
	scoreConfigPathBonus    = 15
	scorePlaceholderPenalty = -15
	scoreFallbackPenalty    = -10
)

func configPathBonus(relPath string) int {
	base := strings.ToLower(filepath.Base(relPath))
	if strings.HasPrefix(base, ".env") {
		return scoreConfigPathBonus
	}
	if strings.HasPrefix(base, "application.") {
		return scoreConfigPathBonus
	}
	if strings.HasPrefix(base, "database.") {
		return scoreConfigPathBonus
	}
	return 0
}

func valuesContainPlaceholder(m map[string]string) bool {
	for _, v := range m {
		if stringContainsPlaceholder(v) {
			return true
		}
	}
	return false
}

func stringContainsPlaceholder(s string) bool {
	return strings.Contains(s, "${") || strings.Contains(s, "%env(")
}

// scoreURLCandidate sets confidence for a parsed DSN/URL candidate.
func scoreURLCandidate(c *Candidate, relPath string, evidence []string) {
	conf := scoreDSNRecognized
	conf += configPathBonus(relPath)
	ev := append([]string(nil), evidence...)
	if stringContainsPlaceholder(c.URL) {
		conf += scorePlaceholderPenalty
		ev = append(ev, "placeholder in URL (lower confidence)")
	}
	c.Confidence = conf
	c.Evidence = ev
}

// finalizeRawURLCandidate scores candidates produced by candidateFromRawURL (DSNs or sqlite file paths).
func finalizeRawURLCandidate(c *Candidate, relPath, parser string, ev []string) {
	if c.DBType == "sqlite" && c.SQLitePath != "" {
		vals := map[string]string{}
		if stringContainsPlaceholder(c.SQLitePath) {
			vals["SQLITE_PATH"] = c.SQLitePath
		}
		scoreSQLiteCandidate(c, relPath, ev, vals)
	} else {
		scoreURLCandidate(c, relPath, ev)
	}
	c.Parser = parser
	c.SourceFile = relPath
}

// scoreTupleCandidate sets confidence for host/port/db/user style candidates.
func scoreTupleCandidate(c *Candidate, relPath string, complete bool, evidence []string, vals map[string]string) {
	conf := scorePartialTuple
	if complete {
		conf = scoreCompleteTuple
	}
	conf += configPathBonus(relPath)
	ev := append([]string(nil), evidence...)
	if vals != nil && valuesContainPlaceholder(vals) {
		conf += scorePlaceholderPenalty
		ev = append(ev, "placeholder in values (lower confidence)")
	}
	c.Confidence = conf
	c.Evidence = ev
}

func scoreSQLiteCandidate(c *Candidate, relPath string, evidence []string, vals map[string]string) {
	conf := scorePartialTuple + 10 // slightly above generic partial when path is explicit
	conf += configPathBonus(relPath)
	ev := append([]string(nil), evidence...)
	if vals != nil && valuesContainPlaceholder(vals) {
		conf += scorePlaceholderPenalty
		ev = append(ev, "placeholder in values (lower confidence)")
	}
	c.Confidence = conf
	c.Evidence = ev
}

func applyFallbackPenalty(c *Candidate) {
	c.Confidence += scoreFallbackPenalty
	if c.Confidence < 0 {
		c.Confidence = 0
	}
	c.Evidence = append(c.Evidence, "fallback project scan")
}

func maxCandidateConfidence(candidates []Candidate) int {
	m := 0
	for _, c := range candidates {
		if c.Confidence > m {
			m = c.Confidence
		}
	}
	return m
}
