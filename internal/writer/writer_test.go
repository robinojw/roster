package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/robinojw/roster/internal/analyser"
)

func TestWriteAll(t *testing.T) {
	outputDir := t.TempDir()
	signals := &analyser.RepoSignals{
		RepoName:  "test-repo",
		Languages: []string{"Go"},
	}

	result, err := WriteAll(outputDir, signals)
	if err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	expectedFiles := []string{
		".roster/signals.json",
		".roster/personas/architect.md",
		"CLAUDE.md",
		"AGENTS.md",
	}
	for _, file := range expectedFiles {
		fullPath := filepath.Join(outputDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", file)
		}
	}
	if len(result.FilesWritten) == 0 {
		t.Error("expected FilesWritten to be populated")
	}
}
