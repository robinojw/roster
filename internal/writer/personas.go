package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/robinojw/roster/internal/personas"
)

func WritePersonas(repoPath string, allPersonas []personas.Persona) ([]string, error) {
	personasDir := filepath.Join(repoPath, ".roster", "personas")
	if err := os.MkdirAll(personasDir, 0755); err != nil {
		return nil, fmt.Errorf("create personas dir: %w", err)
	}

	var written []string
	for _, persona := range allPersonas {
		filename := persona.ID + ".md"
		filePath := filepath.Join(personasDir, filename)
		if err := os.WriteFile(filePath, []byte(persona.Content), 0644); err != nil {
			return nil, fmt.Errorf("write persona %s: %w", persona.ID, err)
		}
		written = append(written, filepath.Join(".roster", "personas", filename))
	}
	return written, nil
}
