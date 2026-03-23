package writer

import (
	"fmt"

	"github.com/robinojw/roster/internal/analyser"
	"github.com/robinojw/roster/internal/personas"
)

type Result struct {
	FilesWritten []string
}

func WriteAll(repoPath string, signals *analyser.RepoSignals) (*Result, error) {
	allPersonas, err := personas.LoadAll()
	if err != nil {
		return nil, fmt.Errorf("load personas: %w", err)
	}

	selectedPersonas := personas.Select(allPersonas, signals)
	result := &Result{}

	signalsFiles, err := WriteSignals(repoPath, signals)
	if err != nil {
		return nil, fmt.Errorf("write signals: %w", err)
	}
	result.FilesWritten = append(result.FilesWritten, signalsFiles...)

	personaFiles, err := WritePersonas(repoPath, selectedPersonas)
	if err != nil {
		return nil, fmt.Errorf("write personas: %w", err)
	}
	result.FilesWritten = append(result.FilesWritten, personaFiles...)

	claudeFiles, err := WriteManagedSection(repoPath, "CLAUDE.md", signals, selectedPersonas)
	if err != nil {
		return nil, fmt.Errorf("write CLAUDE.md: %w", err)
	}
	result.FilesWritten = append(result.FilesWritten, claudeFiles...)

	agentsFiles, err := WriteManagedSection(repoPath, "AGENTS.md", signals, selectedPersonas)
	if err != nil {
		return nil, fmt.Errorf("write AGENTS.md: %w", err)
	}
	result.FilesWritten = append(result.FilesWritten, agentsFiles...)

	return result, nil
}
