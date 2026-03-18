package analyser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectJest(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "jest.config.js"), []byte("module.exports = {}"), 0644)

	signals := &RepoSignals{}
	detectTesting(repoPath, signals)

	if signals.TestFramework != "Jest" {
		t.Errorf("expected Jest, got %s", signals.TestFramework)
	}
}

func TestDetectVitest(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "vitest.config.ts"), []byte("export default {}"), 0644)

	signals := &RepoSignals{}
	detectTesting(repoPath, signals)

	if signals.TestFramework != "Vitest" {
		t.Errorf("expected Vitest, got %s", signals.TestFramework)
	}
}

func TestDetectGoTest(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "main_test.go"), []byte("package main"), 0644)

	signals := &RepoSignals{}
	detectTesting(repoPath, signals)

	if signals.TestFramework != "Go test" {
		t.Errorf("expected Go test, got %s", signals.TestFramework)
	}
}

func TestDetectPytest(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "conftest.py"), []byte(""), 0644)

	signals := &RepoSignals{}
	detectTesting(repoPath, signals)

	if signals.TestFramework != "Pytest" {
		t.Errorf("expected Pytest, got %s", signals.TestFramework)
	}
}

func TestDetectPlaywright(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "playwright.config.ts"), []byte("export default {}"), 0644)

	signals := &RepoSignals{}
	detectTesting(repoPath, signals)

	if !signals.HasE2E {
		t.Error("expected HasE2E to be true")
	}
	if signals.E2EFramework != "Playwright" {
		t.Errorf("expected Playwright, got %s", signals.E2EFramework)
	}
}

func TestDetectCypress(t *testing.T) {
	repoPath := t.TempDir()
	os.WriteFile(filepath.Join(repoPath, "cypress.config.js"), []byte("module.exports = {}"), 0644)

	signals := &RepoSignals{}
	detectTesting(repoPath, signals)

	if !signals.HasE2E {
		t.Error("expected HasE2E to be true")
	}
	if signals.E2EFramework != "Cypress" {
		t.Errorf("expected Cypress, got %s", signals.E2EFramework)
	}
}
