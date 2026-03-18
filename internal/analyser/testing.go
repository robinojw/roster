package analyser

import (
	"os"
	"path/filepath"
	"strings"
)

var testConfigFiles = map[string]string{
	"jest.config.js":   "Jest",
	"jest.config.ts":   "Jest",
	"jest.config.mjs":  "Jest",
	"vitest.config.ts": "Vitest",
	"vitest.config.js": "Vitest",
	"pytest.ini":       "Pytest",
	"conftest.py":      "Pytest",
}

var e2eConfigFiles = map[string]string{
	"playwright.config.ts": "Playwright",
	"playwright.config.js": "Playwright",
	"cypress.config.js":    "Cypress",
	"cypress.config.ts":    "Cypress",
}

func detectTesting(repoPath string, signals *RepoSignals) {
	for configFile, framework := range testConfigFiles {
		if _, err := os.Stat(filepath.Join(repoPath, configFile)); err == nil {
			signals.TestFramework = framework
			break
		}
	}

	if signals.TestFramework == "" && hasGoTests(repoPath) {
		signals.TestFramework = "Go test"
	}

	for configFile, framework := range e2eConfigFiles {
		if _, err := os.Stat(filepath.Join(repoPath, configFile)); err == nil {
			signals.HasE2E = true
			signals.E2EFramework = framework
			break
		}
	}
}

func hasGoTests(repoPath string) bool {
	found := false
	filepath.WalkDir(repoPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") && path != repoPath {
			return filepath.SkipDir
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "_test.go") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}
