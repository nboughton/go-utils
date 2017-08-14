// Package fs is for utility functions related to dealing the filesystem
package fs

import (
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// DiskFree packages the same types of information one gets from the *nix df cmd
type DiskFree struct {
	Total       int64
	Used        int64
	Avail       int64
	PercentUsed int
}

// Exported constants for file sizes
const (
	BS = 1024    // Blocksize
	KB = BS      // Kilobyte
	MB = KB * BS // Megabyte
	GB = MB * BS // Gigabyte
	TB = GB * BS // Terabyte
	PB = TB * BS // Petabyte
)

// Df returns the disk free information for a single mounted device values in KB
func Df(mountPoint string) (f DiskFree, err error) {
	s := syscall.Statfs_t{}
	if err = syscall.Statfs(mountPoint, &s); err != nil {
		return f, err
	}

	f = DiskFree{
		Total: int64(s.Blocks) * MB / int64(s.Bsize),
		Used:  int64(s.Blocks-s.Bfree) * MB / int64(s.Bsize),
	}

	f.Avail = f.Total - f.Used
	f.PercentUsed = int(math.Ceil(float64(f.Used/f.Total) * 100))

	return f, err
}

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
