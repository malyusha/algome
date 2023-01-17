//go:build windows

package util

// WithUserHomeDir returns path where `~` is replaced with user HOME dir.
func WithUserHomeDir(path string) string {
	return path
}
