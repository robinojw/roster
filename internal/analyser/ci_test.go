package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGitHubActions(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, ".github", "workflows"), 0755)
	os.WriteFile(filepath.Join(repoPath, ".github", "workflows", "ci.yml"), []byte("name: CI"), 0644)

	ci := detectCI(repoPath)
	if ci != "GitHub Actions" {
		t.Errorf("expected GitHub Actions, got %s", ci)
	}
}

func TestDetectCircleCI(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, ".circleci"), 0755)
	os.WriteFile(filepath.Join(repoPath, ".circleci", "config.yml"), []byte("version: 2"), 0644)

	ci := detectCI(repoPath)
	if ci != "CircleCI" {
		t.Errorf("expected CircleCI, got %s", ci)
	}
}

func TestDetectBuildkite(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, ".buildkite"), 0755)

	ci := detectCI(repoPath)
	if ci != "Buildkite" {
		t.Errorf("expected Buildkite, got %s", ci)
	}
}

func TestDetectNoCI(t *testing.T) {
	repoPath := t.TempDir()

	ci := detectCI(repoPath)
	if ci != "" {
		t.Errorf("expected empty, got %s", ci)
	}
}
