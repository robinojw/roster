package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/robinojw/roster/internal/analyser"
)

func TestWriteSignals(t *testing.T) {
	outputDir := t.TempDir()
	signals := &analyser.RepoSignals{
		RepoName:  "test-repo",
		Languages: []string{"Go", "TypeScript"},
		FileCount: 42,
	}

	written, err := WriteSignals(outputDir, signals)
	if err != nil {
		t.Fatalf("WriteSignals failed: %v", err)
	}

	signalsPath := filepath.Join(outputDir, ".roster", "signals.json")
	data, err := os.ReadFile(signalsPath)
	if err != nil {
		t.Fatalf("read signals.json: %v", err)
	}

	var loaded analyser.RepoSignals
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("unmarshal signals: %v", err)
	}
	if loaded.RepoName != "test-repo" {
		t.Errorf("expected test-repo, got %s", loaded.RepoName)
	}
	if loaded.FileCount != 42 {
		t.Errorf("expected 42 files, got %d", loaded.FileCount)
	}
	if len(written) != 1 {
		t.Errorf("expected 1 file written, got %d", len(written))
	}
}
