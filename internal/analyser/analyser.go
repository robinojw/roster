package analyser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Analyse(repoPath string) (*RepoSignals, error) {
	signals := &RepoSignals{
		RepoName: filepath.Base(repoPath),
	}

	signals.TopLevelDirs = collectTopLevelDirs(repoPath)
	fileCount, err := countFiles(repoPath)
	if err != nil {
		return nil, fmt.Errorf("count files: %w", err)
	}
	signals.FileCount = fileCount
	signals.Languages = detectLanguages(repoPath)
	signals.Frameworks = detectFrameworks(repoPath)
	detectDesignSystem(repoPath, signals)
	detectTesting(repoPath, signals)
	signals.CIProvider = detectCI(repoPath)
	detectConventions(repoPath, signals)
	signals.FileTree = buildFileTree(repoPath, signals.FileCount)

	return signals, nil
}

func collectTopLevelDirs(repoPath string) []string {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return nil
	}

	var dirs []string
	for _, entry := range entries {
		isHidden := strings.HasPrefix(entry.Name(), ".")
		if entry.IsDir() && !isHidden {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs
}

func countFiles(repoPath string) (int, error) {
	count := 0
	err := filepath.WalkDir(repoPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		isHiddenDir := entry.IsDir() && strings.HasPrefix(entry.Name(), ".") && path != repoPath
		if isHiddenDir {
			return filepath.SkipDir
		}
		if !entry.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}
