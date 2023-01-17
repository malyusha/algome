package util

import (
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// WithUserHomeDir returns path where `~` is replaced with user HOME dir.
func WithUserHomeDir(path string) string {
	if runtime.GOOS == "windows" {
		return path
	}
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	return path
}
