package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectNextJS(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "next.config.js"), []byte("module.exports = {}"), 0644)

	frameworks := detectFrameworks(repoPath)
	if !contains(frameworks, "Next.js") {
		t.Errorf("expected Next.js, got %v", frameworks)
	}
}

func TestDetectReactFromPackageJSON(t *testing.T) {
	repoPath := t.TempDir()
	packageJSON := `{"dependencies": {"react": "^18.0.0"}}`
	os.WriteFile(filepath.Join(repoPath, "package.json"), []byte(packageJSON), 0644)

	frameworks := detectFrameworks(repoPath)
	if !contains(frameworks, "React") {
		t.Errorf("expected React, got %v", frameworks)
	}
}

func TestDetectMultipleFrameworks(t *testing.T) {
	repoPath := t.TempDir()
	packageJSON := `{"dependencies": {"react": "^18.0.0", "express": "^4.0.0"}}`
	os.WriteFile(filepath.Join(repoPath, "package.json"), []byte(packageJSON), 0644)

	frameworks := detectFrameworks(repoPath)
	if !contains(frameworks, "React") || !contains(frameworks, "Express") {
		t.Errorf("expected React and Express, got %v", frameworks)
	}
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
