package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/robinojw/roster/internal/analyser"
	"github.com/robinojw/roster/internal/personas"
)

const startDelimiter = "<!-- roster:start -->"
const endDelimiter = "<!-- roster:end -->"
const fileMode = 0644

var managedTemplate = template.Must(template.New("managed").Parse(`<!-- roster:start -->
## Agent Personas (managed by roster)

### Repo Signals
` + "```json" + `
{{ .SignalsJSON }}
` + "```" + `
{{ if .FileTree }}
### File Tree
` + "```" + `
{{ range .FileTree -}}
{{ . }}
{{ end -}}` + "```" + `
{{ end -}}
### Available Personas
The following agent personas are available in ` + "`.roster/personas/`" + `:

| File | Description |
|------|-------------|
{{ range .Personas -}}
| {{ .ID }}.md | {{ .Description }} |
{{ end -}}

### Instructions
1. Review the repo signals above
2. Scan the persona filenames and descriptions in the table
3. Select the personas relevant to this codebase
4. Read each selected persona file from ` + "`.roster/personas/`" + `
5. For each selected persona, create a hydrated version in ` + "`.claude/agents/`" + ` that fills the ` + "`## Codebase Context`" + ` section with repo-specific details (real paths, real framework names, real conventions)
6. Define how the agents should be orchestrated: routing rules (which keywords and tasks route to which agent), spawn strategy (parallel vs sequential vs on-demand), and any agent dependencies
7. Write the orchestration rules into this file below this section
{{ if .ContextHints }}
### Codebase Context Hints
Use these detected signals when hydrating persona ` + "`## Codebase Context`" + ` sections:

{{ .ContextHints }}{{ end }}
### Orchestration
{{ .RoutingRules }}
<!-- roster:end -->`))

type managedData struct {
	SignalsJSON  string
	FileTree     []string
	Personas     []personas.Persona
	ContextHints string
	RoutingRules string
}

func WriteManagedSection(
	repoPath string,
	filename string,
	signals *analyser.RepoSignals,
	selectedPersonas []personas.Persona,
) ([]string, error) {
	section, err := renderManagedSection(signals, selectedPersonas)
	if err != nil {
		return nil, fmt.Errorf("render managed section: %w", err)
	}

	filePath := filepath.Join(repoPath, filename)
	existing, err := os.ReadFile(filePath)
	isUnexpectedErr := err != nil && !os.IsNotExist(err)
	if isUnexpectedErr {
		return nil, fmt.Errorf("read %s: %w", filename, err)
	}

	var output string
	if len(existing) == 0 {
		output = section + "\n"
	} else {
		output = insertOrReplace(string(existing), section)
	}

	if err := os.WriteFile(filePath, []byte(output), fileMode); err != nil {
		return nil, fmt.Errorf("write %s: %w", filename, err)
	}

	return []string{filename}, nil
}

func renderManagedSection(
	signals *analyser.RepoSignals,
	selectedPersonas []personas.Persona,
) (string, error) {
	signalsJSON, err := json.MarshalIndent(signals, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal signals: %w", err)
	}

	data := managedData{
		SignalsJSON:  string(signalsJSON),
		FileTree:     signals.FileTree,
		Personas:     selectedPersonas,
		ContextHints: buildContextHints(signals),
		RoutingRules: buildRoutingRules(selectedPersonas),
	}

	var buf bytes.Buffer
	if err := managedTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func insertOrReplace(existing string, section string) string {
	startIdx := strings.Index(existing, startDelimiter)
	endIdx := strings.Index(existing, endDelimiter)

	delimitersMissing := startIdx == -1 || endIdx == -1
	if delimitersMissing {
		return existing + "\n" + section + "\n"
	}

	before := existing[:startIdx]
	after := existing[endIdx+len(endDelimiter):]
	return before + section + after
}

func buildContextHints(signals *analyser.RepoSignals) string {
	var hints []string

	hints = appendLanguageHints(hints, signals)
	hints = appendToolchainHints(hints, signals)
	hints = appendInfraHints(hints, signals)

	if len(hints) == 0 {
		return ""
	}
	return strings.Join(hints, "\n")
}

func appendLanguageHints(hints []string, signals *analyser.RepoSignals) []string {
	separator := ", "
	if len(signals.Languages) > 0 {
		hints = append(hints, fmt.Sprintf("- **Languages:** %s", strings.Join(signals.Languages, separator)))
	}
	if len(signals.Frameworks) > 0 {
		hints = append(hints, fmt.Sprintf("- **Frameworks:** %s", strings.Join(signals.Frameworks, separator)))
	}
	if signals.HasDesignSystem {
		hints = append(hints, fmt.Sprintf("- **Design system:** %s", signals.DesignSystemType))
	}
	return hints
}

func appendToolchainHints(hints []string, signals *analyser.RepoSignals) []string {
	if signals.TestFramework != "" {
		hints = append(hints, formatTestHint(signals))
	}
	if signals.LintConfig != "" {
		hints = append(hints, fmt.Sprintf("- **Linter:** %s", signals.LintConfig))
	}
	return hints
}

func formatTestHint(signals *analyser.RepoSignals) string {
	hint := fmt.Sprintf("- **Test framework:** %s", signals.TestFramework)
	if signals.HasE2E {
		hint += fmt.Sprintf(" (E2E: %s)", signals.E2EFramework)
	}
	return hint
}

func appendInfraHints(hints []string, signals *analyser.RepoSignals) []string {
	if signals.CIProvider != "" {
		hints = append(hints, fmt.Sprintf("- **CI:** %s", signals.CIProvider))
	}
	if signals.HasDocker {
		hints = append(hints, "- **Docker:** yes")
	}
	if signals.IsMonorepo {
		hints = append(hints, "- **Monorepo:** yes")
	}
	return hints
}

func buildRoutingRules(selectedPersonas []personas.Persona) string {
	var rules []string

	rules = append(rules, "Route tasks to personas based on trigger keywords:\n")
	for _, p := range selectedPersonas {
		triggers := strings.Join(p.Triggers, ", ")
		rules = append(rules, fmt.Sprintf("- **%s** (%s): %s", p.Name, p.Role, triggers))
	}

	rules = append(rules, "")
	rules = append(rules, "Persona roles define the interaction pattern:")
	rules = append(rules, "- **planning** personas set direction — consult them before major changes")
	rules = append(rules, "- **execution** personas implement — delegate concrete tasks to them")
	rules = append(rules, "- **review** personas verify — route completed work through them")

	return strings.Join(rules, "\n")
}
