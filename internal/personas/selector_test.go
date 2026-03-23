package personas

import (
	"sort"
	"testing"

	"github.com/robinojw/roster/internal/analyser"
)

const (
	testLangGo         = "Go"
	testLangTS         = "TypeScript"
	testLangPython     = "Python"
	testFwReact        = "React"
	testFwNextJS       = "Next.js"
	testFwDjango       = "Django"
	testFwExpress      = "Express"
	testDesignTailwind = "Tailwind CSS"
	testCIGitHub       = "GitHub Actions"
	testIDDesign       = "design"
	testIDAccess       = "accessibility"
	testIDDevops       = "devops"
	testIDAPI          = "api"
	testIDData         = "data"
	testIDPerf         = "performance"

	totalPersonaCount = 11
	corePersonaCount  = 5
)

func loadTestPersonas(t *testing.T) []Persona {
	t.Helper()
	all, err := LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	return all
}

func selectedIDs(selected []Persona) []string {
	ids := make([]string, len(selected))
	for i, p := range selected {
		ids[i] = p.ID
	}
	sort.Strings(ids)
	return ids
}

func hasID(selected []Persona, id string) bool {
	for _, p := range selected {
		if p.ID == id {
			return true
		}
	}
	return false
}

func TestSelectCorePersonasAlwaysIncluded(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages: []string{testLangGo},
	}

	selected := Select(all, signals)

	for _, id := range []string{"architect", "reviewer", "test", "docs", "security"} {
		if !hasID(selected, id) {
			t.Errorf("core persona %q should always be included", id)
		}
	}
}

func TestSelectGoOnlyRepo(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages: []string{testLangGo},
	}

	selected := Select(all, signals)

	for _, id := range []string{testIDDesign, testIDAccess} {
		if hasID(selected, id) {
			t.Errorf("persona %q should not be included for Go-only repo", id)
		}
	}
}

func TestSelectGoWithCIIncludesDevops(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages:  []string{testLangGo},
		CIProvider: testCIGitHub,
	}

	selected := Select(all, signals)

	if !hasID(selected, testIDDevops) {
		t.Error("devops should be included when CI provider is detected")
	}
}

func TestSelectGoWithDockerIncludesDevops(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages: []string{testLangGo},
		HasDocker: true,
	}

	selected := Select(all, signals)

	if !hasID(selected, testIDDevops) {
		t.Error("devops should be included when Docker is detected")
	}
}

func TestSelectReactAppIncludesFrontend(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages:  []string{testLangTS},
		Frameworks: []string{testFwReact, testFwNextJS},
	}

	selected := Select(all, signals)

	for _, id := range []string{testIDDesign, testIDAccess, testIDPerf, testIDAPI} {
		if !hasID(selected, id) {
			t.Errorf("persona %q should be included for React/Next.js app", id)
		}
	}
}

func TestSelectDesignSystemIncludesFrontend(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages:        []string{testLangTS},
		HasDesignSystem:  true,
		DesignSystemType: testDesignTailwind,
	}

	selected := Select(all, signals)

	for _, id := range []string{testIDDesign, testIDAccess} {
		if !hasID(selected, id) {
			t.Errorf("persona %q should be included when design system detected", id)
		}
	}
}

func TestSelectDjangoIncludesDataAndAPI(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages:  []string{testLangPython},
		Frameworks: []string{testFwDjango},
	}

	selected := Select(all, signals)

	for _, id := range []string{testIDData, testIDAPI} {
		if !hasID(selected, id) {
			t.Errorf("persona %q should be included for Django app", id)
		}
	}
}

func TestSelectPythonOnlyIncludesData(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages: []string{testLangPython},
	}

	selected := Select(all, signals)

	if !hasID(selected, testIDData) {
		t.Error("data should be included for Python projects")
	}
}

func TestSelectFullStackApp(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages:        []string{testLangTS, testLangGo, testLangPython},
		Frameworks:       []string{testFwReact, testFwNextJS, testFwExpress, testFwDjango},
		HasDesignSystem:  true,
		DesignSystemType: testDesignTailwind,
		HasStorybook:     true,
		TestFramework:    "Jest",
		HasE2E:           true,
		E2EFramework:     "Playwright",
		CIProvider:       testCIGitHub,
		HasDocker:        true,
	}

	selected := Select(all, signals)

	if len(selected) != totalPersonaCount {
		ids := selectedIDs(selected)
		t.Errorf("full-stack app should include all personas, got %d: %v", len(selected), ids)
	}
}

func TestSelectMinimalRepo(t *testing.T) {
	all := loadTestPersonas(t)
	signals := &analyser.RepoSignals{
		Languages: []string{testLangGo},
	}

	selected := Select(all, signals)

	if len(selected) != corePersonaCount {
		ids := selectedIDs(selected)
		t.Errorf("minimal Go repo should include only core personas, got %d: %v", len(selected), ids)
	}
}
