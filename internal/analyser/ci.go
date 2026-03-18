package analyser

import (
	"os"
	"path/filepath"
)

func detectCI(repoPath string) string {
	if _, err := os.Stat(filepath.Join(repoPath, ".github", "workflows")); err == nil {
		return "GitHub Actions"
	}
	if _, err := os.Stat(filepath.Join(repoPath, ".circleci", "config.yml")); err == nil {
		return "CircleCI"
	}
	if _, err := os.Stat(filepath.Join(repoPath, ".buildkite")); err == nil {
		return "Buildkite"
	}
	return ""
}
