// Package fs is for utility functions related to dealing the filesystem
package fs

import (
	"os"
	"path/filepath"
	"time"
)

// GetBinPath returns the absolute path to the directory of the running binary
func GetBinPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// YearsOld returns the number of years since f modTime was last changed
func YearsOld(f os.FileInfo) float64 {
	return time.Now().Sub(f.ModTime()).Hours() / 24 / 365
}

// WeeksOld returns the number of weeks f modTime was last changed
func WeeksOld(f os.FileInfo) float64 {
	return time.Now().Sub(f.ModTime()).Hours() / 24 / 7
}

// DaysOld returns the number of days f modTime was last changed
func DaysOld(f os.FileInfo) float64 {
	return time.Now().Sub(f.ModTime()).Hours() / 24
}
