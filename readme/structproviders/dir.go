package structproviders

import (
	"fmt"
	"os"
	"path/filepath"
)

// collectFilesOfDirectory returns all files inside given `root` path recursively.
// Returns map of {base_path: [file_name...]}.
func collectFilesOfDirectory(root string) (map[string][]string, error) {
	out := make(map[string][]string, 0)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		// create directory of source for readme
		if err := os.MkdirAll(root, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory of solved problems: %w", err)
		}
		return nil, nil
	}

	// walk through provided root directory and collect all files recursively, ignoring subdirectories
	err := filepath.Walk(root, func(path string, finfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk func error: %w", err)
		}
		dir := filepath.Dir(path)
		if finfo.IsDir() || dir == root {
			return nil
		}

		out[dir] = append(out[dir], finfo.Name())
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("collectFilesOfDirectory error: %w", err)
	}

	return out, nil
}
