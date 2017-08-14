// Package fs is for utility functions related to dealing the filesystem
package fs

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// GetBinPath returns the absolute path to the directory of the running binary
func GetBinPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// GetUID returns the Unix user id for f
func GetUID(f os.FileInfo) int {
	return int(f.Sys().(*syscall.Stat_t).Uid)
}

// GetUname returns the Unix username for f
func GetUname(f os.FileInfo) string {
	uid := strconv.Itoa(GetUID(f))
	u, _ := user.LookupId(uid)
	return u.Username
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
