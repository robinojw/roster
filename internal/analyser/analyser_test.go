package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyseEmptyRepo(t *testing.T) {
	repoPath := t.TempDir()

	signals, err := Analyse(repoPath)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if signals.RepoName != filepath.Base(repoPath) {
		t.Errorf("expected repo name %s, got %s", filepath.Base(repoPath), signals.RepoName)
	}
	if signals.FileCount != 0 {
		t.Errorf("expected 0 files, got %d", signals.FileCount)
	}
}

func TestAnalyseCountsFiles(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(repoPath, "README.md"), []byte("# hello"), 0644)

	signals, err := Analyse(repoPath)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if signals.FileCount != 2 {
		t.Errorf("expected 2 files, got %d", signals.FileCount)
	}
}

func TestAnalyseTopLevelDirs(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, "src"), 0755)
	os.MkdirAll(filepath.Join(repoPath, "infra"), 0755)
	os.MkdirAll(filepath.Join(repoPath, ".git"), 0755)

	signals, err := Analyse(repoPath)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if len(signals.TopLevelDirs) != 2 {
		t.Errorf("expected 2 top-level dirs (excluding .git), got %v", signals.TopLevelDirs)
	}
}
