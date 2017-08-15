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
	total       float64
	used        float64
	avail       float64
	percentUsed int
}

// Exported constants for file sizes
const (
	KB = float64(1024) // Kilobyte
	MB = KB * KB       // Megabyte
	GB = MB * KB       // Gigabyte
	TB = GB * KB       // Terabyte
	PB = TB * KB       // Petabyte
)

// NewDf returns the disk free information for specified mount point in 1KB blocks
func NewDf(mountPoint string) (f *DiskFree, err error) {
	s := syscall.Statfs_t{}
	if err = syscall.Statfs(mountPoint, &s); err != nil {
		return f, err
	}

	f = &DiskFree{
		total: (float64(s.Blocks) * float64(s.Bsize)),
		used:  (float64(s.Blocks-s.Bfree) * float64(s.Bsize)),
	}

	f.avail = f.total - f.used
	f.percentUsed = int(math.Ceil(float64(f.used) / float64(f.total) * 100))

	return f, err
}

// Total takes a value type (fs.MB, fs.GB etc) and returns the appropriate value.
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Total(fs.GB))
//
// Floats are used in order to allow for a reasonable degree of accuracy when
// dealing with large value numbers i.e GB and TB etc...
func (df *DiskFree) Total(valType float64) float64 {
	return df.total / valType
}

// Used works the same way as Total but returns the amount of space currently used.
func (df *DiskFree) Used(valType float64) float64 {
	return df.used / valType
}

// Avail works as above but returns the amount of space available
func (df *DiskFree) Avail(valType float64) float64 {
	return df.avail / valType
}

// PercentUsed returns the percentage of space used as an int
func (df *DiskFree) PercentUsed() int {
	return df.percentUsed
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
