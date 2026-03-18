package analyser

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var lintConfigFiles = map[string]string{
	".eslintrc.json": "ESLint",
	".eslintrc.js":   "ESLint",
	".eslintrc.yml":  "ESLint",
	".eslintrc.yaml": "ESLint",
	"eslint.config.js":  "ESLint",
	"eslint.config.mjs": "ESLint",
	".golangci.yml":     "golangci-lint",
	".golangci.yaml":    "golangci-lint",
	"ruff.toml":         "Ruff",
	".prettierrc":       "Prettier",
	".prettierrc.json":  "Prettier",
	".prettierrc.js":    "Prettier",
}

func detectConventions(repoPath string, signals *RepoSignals) {
	for configFile, linter := range lintConfigFiles {
		if _, err := os.Stat(filepath.Join(repoPath, configFile)); err == nil {
			signals.LintConfig = linter
			break
		}
	}

	if _, err := os.Stat(filepath.Join(repoPath, "Dockerfile")); err == nil {
		signals.HasDocker = true
	}

	signals.IsMonorepo = detectMonorepo(repoPath)
}

func detectMonorepo(repoPath string) bool {
	monoDirs := []string{"packages", "apps"}
	for _, dir := range monoDirs {
		if _, err := os.Stat(filepath.Join(repoPath, dir)); err == nil {
			return true
		}
	}

	pkgPath := filepath.Join(repoPath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false
	}

	var pkg struct {
		Workspaces json.RawMessage `json:"workspaces"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}
	return pkg.Workspaces != nil
}
