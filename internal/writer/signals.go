package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/robinojw/roster/internal/analyser"
)

func WriteSignals(repoPath string, signals *analyser.RepoSignals) ([]string, error) {
	rosterDir := filepath.Join(repoPath, ".roster")
	if err := os.MkdirAll(rosterDir, 0755); err != nil {
		return nil, fmt.Errorf("create .roster dir: %w", err)
	}

	data, err := json.MarshalIndent(signals, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal signals: %w", err)
	}

	signalsPath := filepath.Join(rosterDir, "signals.json")
	if err := os.WriteFile(signalsPath, data, 0644); err != nil {
		return nil, fmt.Errorf("write signals.json: %w", err)
	}

	return []string{".roster/signals.json"}, nil
}
