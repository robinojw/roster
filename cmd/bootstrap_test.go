package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBootstrapCommand(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(repoPath, "go.mod"), []byte("module test"), 0644)

	rootCmd.SetArgs([]string{"bootstrap", "--path", repoPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("bootstrap command failed: %v", err)
	}

	expectedFiles := []string{
		".roster/signals.json",
		".roster/personas/architect.md",
		"CLAUDE.md",
		"AGENTS.md",
	}
	for _, file := range expectedFiles {
		fullPath := filepath.Join(repoPath, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", file)
		}
	}
}

func TestBootstrapDryRun(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)

	rootCmd.SetArgs([]string{"bootstrap", "--path", repoPath, "--dry-run"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("bootstrap --dry-run failed: %v", err)
	}

	signalsPath := filepath.Join(repoPath, ".roster", "signals.json")
	if _, err := os.Stat(signalsPath); !os.IsNotExist(err) {
		t.Error("dry-run should not create files")
	}
}
