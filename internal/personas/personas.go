package personas

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Persona struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Triggers    []string `yaml:"triggers"`
	Content     string   `yaml:"-"`
}

func LoadAll() ([]Persona, error) {
	entries, err := personaFS.ReadDir("lib")
	if err != nil {
		return nil, fmt.Errorf("read persona lib: %w", err)
	}

	var personas []Persona
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		persona, err := loadPersona("lib/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("load persona %s: %w", entry.Name(), err)
		}
		personas = append(personas, persona)
	}
	return personas, nil
}

func loadPersona(path string) (Persona, error) {
	data, err := personaFS.ReadFile(path)
	if err != nil {
		return Persona{}, fmt.Errorf("read file: %w", err)
	}

	content := string(data)
	frontmatter, body, err := splitFrontmatter(content)
	if err != nil {
		return Persona{}, fmt.Errorf("parse frontmatter: %w", err)
	}

	var persona Persona
	if err := yaml.Unmarshal([]byte(frontmatter), &persona); err != nil {
		return Persona{}, fmt.Errorf("unmarshal frontmatter: %w", err)
	}
	persona.Content = strings.TrimSpace(body)
	return persona, nil
}

func splitFrontmatter(content string) (string, string, error) {
	const delimiter = "---"
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, delimiter) {
		return "", "", fmt.Errorf("missing opening frontmatter delimiter")
	}
	rest := trimmed[len(delimiter):]
	endIndex := strings.Index(rest, "\n"+delimiter)
	if endIndex == -1 {
		return "", "", fmt.Errorf("missing closing frontmatter delimiter")
	}
	frontmatter := rest[:endIndex]
	body := rest[endIndex+len("\n"+delimiter):]
	return frontmatter, body, nil
}
