package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robinojw/roster/internal/analyser"
	"github.com/robinojw/roster/internal/personas"
)

func TestWriteManagedSectionCreatesFile(t *testing.T) {
	outputDir := t.TempDir()
	signals := &analyser.RepoSignals{RepoName: "test-repo"}
	allPersonas, _ := personas.LoadAll()

	_, err := WriteManagedSection(outputDir, "CLAUDE.md", signals, allPersonas)
	if err != nil {
		t.Fatalf("WriteManagedSection failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	if !strings.Contains(string(content), "<!-- roster:start -->") {
		t.Error("missing roster:start delimiter")
	}
	if !strings.Contains(string(content), "<!-- roster:end -->") {
		t.Error("missing roster:end delimiter")
	}
}

func TestWriteManagedSectionPreservesExisting(t *testing.T) {
	outputDir := t.TempDir()
	existingContent := "# My Project\n\nSome existing content.\n"
	os.WriteFile(filepath.Join(outputDir, "CLAUDE.md"), []byte(existingContent), 0644)

	signals := &analyser.RepoSignals{RepoName: "test-repo"}
	allPersonas, _ := personas.LoadAll()

	_, err := WriteManagedSection(outputDir, "CLAUDE.md", signals, allPersonas)
	if err != nil {
		t.Fatalf("WriteManagedSection failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	if !strings.Contains(string(content), "# My Project") {
		t.Error("existing content was overwritten")
	}
	if !strings.Contains(string(content), "<!-- roster:start -->") {
		t.Error("managed section not appended")
	}
}

func TestWriteManagedSectionReplacesExisting(t *testing.T) {
	outputDir := t.TempDir()
	existingContent := "# My Project\n\n<!-- roster:start -->\nold content\n<!-- roster:end -->\n\nMore stuff.\n"
	os.WriteFile(filepath.Join(outputDir, "CLAUDE.md"), []byte(existingContent), 0644)

	signals := &analyser.RepoSignals{RepoName: "test-repo"}
	allPersonas, _ := personas.LoadAll()

	_, err := WriteManagedSection(outputDir, "CLAUDE.md", signals, allPersonas)
	if err != nil {
		t.Fatalf("WriteManagedSection failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	if strings.Contains(string(content), "old content") {
		t.Error("old managed section was not replaced")
	}
	if !strings.Contains(string(content), "More stuff.") {
		t.Error("content after managed section was lost")
	}
	if !strings.Contains(string(content), "<!-- roster:start -->") {
		t.Error("new managed section not written")
	}
}
