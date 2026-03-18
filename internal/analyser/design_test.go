package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectTailwind(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "tailwind.config.js"), []byte("module.exports = {}"), 0644)

	signals := &RepoSignals{}
	detectDesignSystem(repoPath, signals)

	if !signals.HasDesignSystem {
		t.Error("expected HasDesignSystem to be true")
	}
	if signals.DesignSystemType != "Tailwind CSS" {
		t.Errorf("expected Tailwind CSS, got %s", signals.DesignSystemType)
	}
}

func TestDetectStorybook(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, ".storybook"), 0755)

	signals := &RepoSignals{}
	detectDesignSystem(repoPath, signals)

	if !signals.HasStorybook {
		t.Error("expected HasStorybook to be true")
	}
}

func TestDetectStyledComponents(t *testing.T) {
	repoPath := t.TempDir()
	packageJSON := `{"dependencies": {"styled-components": "^5.0.0"}}`
	os.WriteFile(filepath.Join(repoPath, "package.json"), []byte(packageJSON), 0644)

	signals := &RepoSignals{}
	detectDesignSystem(repoPath, signals)

	if !signals.HasDesignSystem {
		t.Error("expected HasDesignSystem to be true")
	}
	if signals.DesignSystemType != "styled-components" {
		t.Errorf("expected styled-components, got %s", signals.DesignSystemType)
	}
}

func TestDetectCSSModules(t *testing.T) {
	repoPath := t.TempDir()
	os.MkdirAll(filepath.Join(repoPath, "src"), 0755)
	os.WriteFile(filepath.Join(repoPath, "src", "app.module.css"), []byte("body{}"), 0644)

	signals := &RepoSignals{}
	detectDesignSystem(repoPath, signals)

	if !signals.HasDesignSystem {
		t.Error("expected HasDesignSystem to be true")
	}
	if signals.DesignSystemType != "CSS Modules" {
		t.Errorf("expected CSS Modules, got %s", signals.DesignSystemType)
	}
}

func TestDetectNoDesignSystem(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("package main"), 0644)

	signals := &RepoSignals{}
	detectDesignSystem(repoPath, signals)

	if signals.HasDesignSystem {
		t.Error("expected HasDesignSystem to be false")
	}
}
