package walk

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func Files(root string, keep func(string) bool) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		clean := filepath.ToSlash(path)
		if d.IsDir() {
			if d.Name() == ".cache" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if keep(clean) {
			files = append(files, clean)
		}
		return nil
	})
	return files, err
}

func HasPrefix(path string, prefixes ...string) bool {
	path = filepath.ToSlash(path)
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
