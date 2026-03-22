package projectscan

import (
	"os"
	"path/filepath"
	"sort"
)

const lowConfidenceThreshold = 50

// ScanProject discovers database connection candidates using framework-agnostic
// heuristics (env keys, URLs, JSON/YAML config files), then optionally falls
// back to a broader project scan when nothing strong is found.
func ScanProject(root string) (Result, error) {
	structured, err := scanStructuredFiles(root)
	if err != nil {
		return Result{}, err
	}

	var candidates []Candidate
	candidates = append(candidates, structured...)

	if len(structured) == 0 || maxCandidateConfidence(structured) < lowConfidenceThreshold {
		fb, err := fallbackScanProject(root)
		if err != nil {
			return Result{}, err
		}

		candidates = append(candidates, fb...)
	}

	candidates = dedupCandidates(candidates)

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Confidence > candidates[j].Confidence
	})

	return Result{Candidates: candidates}, nil
}

// IsLikelyProjectDir returns true if the directory looks like a software project root.
func IsLikelyProjectDir(root string) bool {
	markers := []string{
		".git",
		"go.mod",
		"package.json",
		"composer.json",
		"pom.xml",
		"build.gradle",
		"settings.gradle",
		"requirements.txt",
		"pyproject.toml",
	}

	for _, marker := range markers {
		if _, err := os.Stat(filepath.Join(root, marker)); err == nil {
			return true
		}
	}

	return false
}
