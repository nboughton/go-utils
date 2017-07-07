package fs

import (
	"os"
	"path/filepath"
)

// getBinPath returns the absolute path to the directory of the running binary
func getBinPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
