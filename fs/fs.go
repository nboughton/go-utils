// Package fs is for utility functions related to dealing the filesystem
package fs

import (
	"os"
	"path/filepath"
)

// GetBinPath returns the absolute path to the directory of the running binary
func GetBinPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
