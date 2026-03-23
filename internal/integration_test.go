package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robinojw/roster/internal/analyser"
	"github.com/robinojw/roster/internal/writer"
)

func TestIntegrationFullPipeline(t *testing.T) {
	repoPath := t.TempDir()

	packageJSON := `{
		"dependencies": {
			"react": "^18.0.0",
			"next": "^14.0.0"
		},
		"devDependencies": {
			"jest": "^29.0.0"
		}
	}`
	os.WriteFile(filepath.Join(repoPath, "package.json"), []byte(packageJSON), 0644)
	os.WriteFile(filepath.Join(repoPath, "next.config.js"), []byte("module.exports = {}"), 0644)
	os.WriteFile(filepath.Join(repoPath, "tailwind.config.js"), []byte("module.exports = {}"), 0644)
	os.WriteFile(filepath.Join(repoPath, "jest.config.js"), []byte("module.exports = {}"), 0644)

	os.MkdirAll(filepath.Join(repoPath, ".github", "workflows"), 0755)
	os.WriteFile(filepath.Join(repoPath, ".github", "workflows", "ci.yml"), []byte("name: CI"), 0644)

	os.MkdirAll(filepath.Join(repoPath, "src"), 0755)
	os.WriteFile(filepath.Join(repoPath, "src", "App.tsx"), []byte("export default function App() {}"), 0644)
	os.WriteFile(filepath.Join(repoPath, "src", "index.tsx"), []byte("import App from './App'"), 0644)

	os.WriteFile(filepath.Join(repoPath, "go.mod"), []byte("module test"), 0644)
	os.WriteFile(filepath.Join(repoPath, "main_test.go"), []byte("package main"), 0644)

	signals, err := analyser.Analyse(repoPath)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}

	if !contains(signals.Languages, "TypeScript") {
		t.Errorf("expected TypeScript in languages, got %v", signals.Languages)
	}
	if !contains(signals.Languages, "Go") {
		t.Errorf("expected Go in languages, got %v", signals.Languages)
	}
	if !contains(signals.Frameworks, "Next.js") {
		t.Errorf("expected Next.js in frameworks, got %v", signals.Frameworks)
	}
	if !contains(signals.Frameworks, "React") {
		t.Errorf("expected React in frameworks, got %v", signals.Frameworks)
	}
	if !signals.HasDesignSystem {
		t.Error("expected HasDesignSystem to be true (Tailwind)")
	}
	if signals.TestFramework != "Jest" {
		t.Errorf("expected Jest test framework, got %s", signals.TestFramework)
	}
	if signals.CIProvider != "GitHub Actions" {
		t.Errorf("expected GitHub Actions CI, got %s", signals.CIProvider)
	}

	result, err := writer.WriteAll(repoPath, signals)
	if err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	signalsPath := filepath.Join(repoPath, ".roster", "signals.json")
	signalsData, err := os.ReadFile(signalsPath)
	if err != nil {
		t.Fatalf("read signals.json: %v", err)
	}
	var loadedSignals analyser.RepoSignals
	if err := json.Unmarshal(signalsData, &loadedSignals); err != nil {
		t.Fatalf("unmarshal signals: %v", err)
	}
	if loadedSignals.RepoName == "" {
		t.Error("signals.json missing repo name")
	}

	personasDir := filepath.Join(repoPath, ".roster", "personas")
	entries, err := os.ReadDir(personasDir)
	if err != nil {
		t.Fatalf("read personas dir: %v", err)
	}
	if len(entries) != 10 {
		t.Errorf("expected 10 persona files (signal-selected), got %d", len(entries))
	}

	claudeContent, err := os.ReadFile(filepath.Join(repoPath, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	claudeStr := string(claudeContent)
	if !strings.Contains(claudeStr, "<!-- roster:start -->") {
		t.Error("CLAUDE.md missing roster:start")
	}
	if !strings.Contains(claudeStr, "<!-- roster:end -->") {
		t.Error("CLAUDE.md missing roster:end")
	}
	if !strings.Contains(claudeStr, "TypeScript") {
		t.Error("CLAUDE.md managed section missing TypeScript signal")
	}
	if !strings.Contains(claudeStr, "architect.md") {
		t.Error("CLAUDE.md missing persona table entries")
	}

	agentsContent, err := os.ReadFile(filepath.Join(repoPath, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	if !strings.Contains(string(agentsContent), "<!-- roster:start -->") {
		t.Error("AGENTS.md missing roster:start")
	}

	if len(result.FilesWritten) < 13 {
		t.Errorf("expected at least 13 files written, got %d", len(result.FilesWritten))
	}
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
