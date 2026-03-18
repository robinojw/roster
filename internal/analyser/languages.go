package analyser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var extensionToLanguage = map[string]string{
	".go":    "Go",
	".ts":    "TypeScript",
	".tsx":   "TypeScript",
	".js":    "JavaScript",
	".jsx":   "JavaScript",
	".py":    "Python",
	".rb":    "Ruby",
	".rs":    "Rust",
	".java":  "Java",
	".kt":    "Kotlin",
	".swift": "Swift",
	".cs":    "C#",
	".php":   "PHP",
}

var manifestToLanguage = map[string]string{
	"go.mod":           "Go",
	"package.json":     "JavaScript",
	"requirements.txt": "Python",
	"pyproject.toml":   "Python",
	"Gemfile":          "Ruby",
	"Cargo.toml":       "Rust",
	"pom.xml":          "Java",
	"build.gradle":     "Java",
}

func detectLanguages(repoPath string) []string {
	seen := make(map[string]bool)

	for manifest, language := range manifestToLanguage {
		manifestPath := filepath.Join(repoPath, manifest)
		if _, err := os.Stat(manifestPath); err == nil {
			seen[language] = true
		}
	}

	filepath.WalkDir(repoPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		isHiddenDir := entry.IsDir() && strings.HasPrefix(entry.Name(), ".") && path != repoPath
		if isHiddenDir {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		extension := strings.ToLower(filepath.Ext(entry.Name()))
		if language, found := extensionToLanguage[extension]; found {
			seen[language] = true
		}
		return nil
	})

	if seen["TypeScript"] && seen["JavaScript"] {
		delete(seen, "JavaScript")
	}

	languages := make([]string, 0, len(seen))
	for language := range seen {
		languages = append(languages, language)
	}
	sort.Strings(languages)
	return languages
}
