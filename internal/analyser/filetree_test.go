package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildFileTreeSmallRepo(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, "src"), 0755)
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(repoPath, "src", "app.go"), []byte("package src"), 0644)
	os.WriteFile(filepath.Join(repoPath, "README.md"), []byte("# hello"), 0644)

	tree := buildFileTree(repoPath, 3)

	if tree == nil {
		t.Fatal("expected file tree for small repo")
	}
	if len(tree) != 4 {
		t.Errorf("expected 4 entries (3 files + 1 dir), got %d: %v", len(tree), tree)
	}
	if !contains(tree, "main.go") {
		t.Errorf("expected main.go in tree, got %v", tree)
	}
	if !contains(tree, "src/") {
		t.Errorf("expected src/ in tree, got %v", tree)
	}
	if !contains(tree, "src/app.go") {
		t.Errorf("expected src/app.go in tree, got %v", tree)
	}
}

func TestBuildFileTreeExcludesHiddenDirs(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, ".git", "objects"), 0755)
	os.WriteFile(filepath.Join(repoPath, ".git", "config"), []byte(""), 0644)
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)

	tree := buildFileTree(repoPath, 1)

	for _, path := range tree {
		if filepath.Base(path) == ".git" || filepath.Base(path) == "config" {
			t.Errorf("hidden dir contents should be excluded, found %s", path)
		}
	}
}

func TestBuildFileTreeReturnsNilForLargeRepo(t *testing.T) {
	repoPath := t.TempDir()
	tree := buildFileTree(repoPath, fileTreeMaxFiles+1)

	if tree != nil {
		t.Error("expected nil file tree for large repo")
	}
}

func TestBuildFileTreeIntegrationWithAnalyse(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(repoPath, "go.mod"), []byte("module test"), 0644)

	signals, err := Analyse(repoPath)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if signals.FileTree == nil {
		t.Error("expected FileTree to be populated for small repo")
	}
	if len(signals.FileTree) != 2 {
		t.Errorf("expected 2 entries, got %d: %v", len(signals.FileTree), signals.FileTree)
	}
}
