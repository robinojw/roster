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
<!-- roster:end -->`))

type managedData struct {
	SignalsJSON string
	FileTree    []string
	Personas    []personas.Persona
}

func WriteManagedSection(
	repoPath string,
	filename string,
	signals *analyser.RepoSignals,
	allPersonas []personas.Persona,
) ([]string, error) {
	section, err := renderManagedSection(signals, allPersonas)
	if err != nil {
		return nil, fmt.Errorf("render managed section: %w", err)
	}

	filePath := filepath.Join(repoPath, filename)
	existing, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read %s: %w", filename, err)
	}

	var output string
	if len(existing) == 0 {
		output = section + "\n"
	} else {
		output = insertOrReplace(string(existing), section)
	}

	if err := os.WriteFile(filePath, []byte(output), 0644); err != nil {
		return nil, fmt.Errorf("write %s: %w", filename, err)
	}

	return []string{filename}, nil
}

func renderManagedSection(
	signals *analyser.RepoSignals,
	allPersonas []personas.Persona,
) (string, error) {
	signalsJSON, err := json.MarshalIndent(signals, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal signals: %w", err)
	}

	data := managedData{
		SignalsJSON: string(signalsJSON),
		FileTree:    signals.FileTree,
		Personas:    allPersonas,
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

	if startIdx == -1 || endIdx == -1 {
		return existing + "\n" + section + "\n"
	}

	before := existing[:startIdx]
	after := existing[endIdx+len(endDelimiter):]
	return before + section + after
}
