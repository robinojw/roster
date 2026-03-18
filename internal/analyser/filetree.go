package analyser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const fileTreeMaxFiles = 500

func buildFileTree(repoPath string, fileCount int) []string {
	if fileCount > fileTreeMaxFiles {
		return nil
	}

	var paths []string
	filepath.WalkDir(repoPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		isHiddenDir := entry.IsDir() && strings.HasPrefix(entry.Name(), ".") && path != repoPath
		if isHiddenDir {
			return filepath.SkipDir
		}
		if path == repoPath {
			return nil
		}

		relativePath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return nil
		}

		if entry.IsDir() {
			relativePath += "/"
		}
		paths = append(paths, relativePath)
		return nil
	})

	sort.Strings(paths)
	return paths
}
