package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/robinojw/roster/internal/personas"
)

func TestWritePersonas(t *testing.T) {
	outputDir := t.TempDir()
	allPersonas, err := personas.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	written, err := WritePersonas(outputDir, allPersonas)
	if err != nil {
		t.Fatalf("WritePersonas failed: %v", err)
	}

	if len(written) != 11 {
		t.Errorf("expected 11 files written, got %d", len(written))
	}

	for _, persona := range allPersonas {
		personaPath := filepath.Join(outputDir, ".roster", "personas", persona.ID+".md")
		data, err := os.ReadFile(personaPath)
		if err != nil {
			t.Errorf("expected persona file %s to exist: %v", persona.ID, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("persona file %s is empty", persona.ID)
		}
	}
}
