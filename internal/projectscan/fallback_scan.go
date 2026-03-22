package projectscan

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func fallbackScanProject(root string) ([]Candidate, error) {
	var candidates []Candidate

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
		if err != nil || info.Size() > 1024*1024 {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(b)
		rel := relOrAbs(root, path)

		if strings.Contains(text, "://") {
			if c, ok := candidateFromURLInText(text, rel, "fallback"); ok {
				applyFallbackPenalty(&c)
				candidates = append(candidates, c)

				return nil
			}
		}

		kv := grepLikelyDBKeys(text)
		if len(kv) == 0 {
			return nil
		}

		if c, ok := candidateFromMap(kv, rel, "fallback"); ok {
			applyFallbackPenalty(&c)
			candidates = append(candidates, c)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dedupCandidates(candidates), nil
}

func grepLikelyDBKeys(text string) map[string]string {
	out := map[string]string{}
	targetKeys := allKeysForGrep()

	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := normalizeKey(parts[0])
			if _, ok := targetKeys[key]; ok {
				out[key] = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			}

			continue
		}

		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := normalizeKey(parts[0])
			if _, ok := targetKeys[key]; ok {
				out[key] = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			}
		}
	}

	return out
}
