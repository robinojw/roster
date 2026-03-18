package analyser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

var configFileToFramework = map[string]string{
	"next.config.js":     "Next.js",
	"next.config.mjs":    "Next.js",
	"next.config.ts":     "Next.js",
	"nuxt.config.ts":     "Nuxt",
	"nuxt.config.js":     "Nuxt",
	"svelte.config.js":   "SvelteKit",
	"angular.json":       "Angular",
	"tailwind.config.js":  "Tailwind CSS",
	"tailwind.config.ts":  "Tailwind CSS",
	"tailwind.config.mjs": "Tailwind CSS",
}

var depToFramework = map[string]string{
	"react":   "React",
	"vue":     "Vue",
	"svelte":  "Svelte",
	"express": "Express",
	"fastify": "Fastify",
	"django":  "Django",
	"flask":   "Flask",
	"rails":   "Rails",
	"gin":     "Gin",
	"fiber":   "Fiber",
	"echo":    "Echo",
}

func detectFrameworks(repoPath string) []string {
	seen := make(map[string]bool)

	for configFile, framework := range configFileToFramework {
		path := filepath.Join(repoPath, configFile)
		if _, err := os.Stat(path); err == nil {
			seen[framework] = true
		}
	}

	deps := readPackageJSONDeps(repoPath)
	for dep, framework := range depToFramework {
		if _, found := deps[dep]; found {
			seen[framework] = true
		}
	}

	frameworks := make([]string, 0, len(seen))
	for framework := range seen {
		frameworks = append(frameworks, framework)
	}
	sort.Strings(frameworks)
	return frameworks
}

func readPackageJSONDeps(repoPath string) map[string]bool {
	deps := make(map[string]bool)
	pkgPath := filepath.Join(repoPath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return deps
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return deps
	}

	for dep := range pkg.Dependencies {
		deps[dep] = true
	}
	for dep := range pkg.DevDependencies {
		deps[dep] = true
	}
	return deps
}
