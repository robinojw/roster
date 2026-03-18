package analyser

import (
	"os"
	"path/filepath"
	"strings"
)

var designConfigFiles = map[string]string{
	"tailwind.config.js":  "Tailwind CSS",
	"tailwind.config.ts":  "Tailwind CSS",
	"tailwind.config.mjs": "Tailwind CSS",
}

var designDeps = map[string]string{
	"styled-components": "styled-components",
	"@emotion/react":    "Emotion",
	"@chakra-ui/react":  "Chakra UI",
	"@mui/material":     "Material UI",
}

func detectDesignSystem(repoPath string, signals *RepoSignals) {
	if _, err := os.Stat(filepath.Join(repoPath, ".storybook")); err == nil {
		signals.HasStorybook = true
	}

	for configFile, designType := range designConfigFiles {
		if _, err := os.Stat(filepath.Join(repoPath, configFile)); err == nil {
			signals.HasDesignSystem = true
			signals.DesignSystemType = designType
			return
		}
	}

	deps := readPackageJSONDeps(repoPath)
	for dep, designType := range designDeps {
		if deps[dep] {
			signals.HasDesignSystem = true
			signals.DesignSystemType = designType
			return
		}
	}

	if hasCSSModules(repoPath) {
		signals.HasDesignSystem = true
		signals.DesignSystemType = "CSS Modules"
	}
}

func hasCSSModules(repoPath string) bool {
	found := false
	filepath.WalkDir(repoPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") && path != repoPath {
			return filepath.SkipDir
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".module.css") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}
