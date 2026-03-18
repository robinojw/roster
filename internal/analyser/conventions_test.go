package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectESLint(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, ".eslintrc.json"), []byte("{}"), 0644)

	signals := &RepoSignals{}
	detectConventions(repoPath, signals)

	if signals.LintConfig != "ESLint" {
		t.Errorf("expected ESLint, got %s", signals.LintConfig)
	}
}

func TestDetectGolangciLint(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, ".golangci.yml"), []byte("version: 2"), 0644)

	signals := &RepoSignals{}
	detectConventions(repoPath, signals)

	if signals.LintConfig != "golangci-lint" {
		t.Errorf("expected golangci-lint, got %s", signals.LintConfig)
	}
}

func TestDetectRuff(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "ruff.toml"), []byte(""), 0644)

	signals := &RepoSignals{}
	detectConventions(repoPath, signals)

	if signals.LintConfig != "Ruff" {
		t.Errorf("expected Ruff, got %s", signals.LintConfig)
	}
}

func TestDetectDocker(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "Dockerfile"), []byte("FROM alpine"), 0644)

	signals := &RepoSignals{}
	detectConventions(repoPath, signals)

	if !signals.HasDocker {
		t.Error("expected HasDocker to be true")
	}
}

func TestDetectMonorepoFromPackagesDir(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, "packages", "app"), 0755)

	signals := &RepoSignals{}
	detectConventions(repoPath, signals)

	if !signals.IsMonorepo {
		t.Error("expected IsMonorepo to be true")
	}
}

func TestDetectMonorepoFromAppsDir(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, "apps", "web"), 0755)

	signals := &RepoSignals{}
	detectConventions(repoPath, signals)

	if !signals.IsMonorepo {
		t.Error("expected IsMonorepo to be true")
	}
}

func TestDetectMonorepoFromWorkspaces(t *testing.T) {
	repoPath := t.TempDir()
	packageJSON := `{"workspaces": ["packages/*"]}`
	os.WriteFile(filepath.Join(repoPath, "package.json"), []byte(packageJSON), 0644)

	signals := &RepoSignals{}
	detectConventions(repoPath, signals)

	if !signals.IsMonorepo {
		t.Error("expected IsMonorepo to be true from workspaces")
	}
}
