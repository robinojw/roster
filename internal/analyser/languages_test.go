package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectLanguagesFromExtensions(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(repoPath, "app.ts"), []byte("const x = 1"), 0644)
	os.WriteFile(filepath.Join(repoPath, "style.css"), []byte("body{}"), 0644)

	languages := detectLanguages(repoPath)

	languageSet := make(map[string]bool)
	for _, lang := range languages {
		languageSet[lang] = true
	}
	if !languageSet["Go"] {
		t.Errorf("expected Go, got %v", languages)
	}
	if !languageSet["TypeScript"] {
		t.Errorf("expected TypeScript, got %v", languages)
	}
}

func TestDetectLanguagesFromGoMod(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "go.mod"), []byte("module test"), 0644)

	languages := detectLanguages(repoPath)

	languageSet := make(map[string]bool)
	for _, lang := range languages {
		languageSet[lang] = true
	}
	if !languageSet["Go"] {
		t.Errorf("expected Go from go.mod, got %v", languages)
	}
}

func TestDetectLanguagesFromPackageJSON(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "package.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(repoPath, "index.tsx"), []byte(""), 0644)

	languages := detectLanguages(repoPath)

	languageSet := make(map[string]bool)
	for _, lang := range languages {
		languageSet[lang] = true
	}
	if !languageSet["TypeScript"] {
		t.Errorf("expected TypeScript from .tsx, got %v", languages)
	}
}
