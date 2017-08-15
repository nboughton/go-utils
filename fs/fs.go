// Package fs is for utility functions related to dealing with Linux filesystems
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

// Exported constants for file size calculations
const (
	KB = 1024    // Kilobyte
	MB = KB * KB // Megabyte
	GB = MB * KB // Gigabyte
	TB = GB * KB // Terabyte
	PB = TB * KB // Petabyte
)

// NewDf returns the disk free information for specified mount point in 1KB blocks
// Values are stored as float64s in order to allow for a reasonable degree of accuracy when
// dealing with large value numbers i.e 1.7GB, 3.4TB etc...
func NewDf(mountPoint string) (f *DiskFree, err error) {
	s := syscall.Statfs_t{}
	if err = syscall.Statfs(mountPoint, &s); err != nil {
		return f, err
	}

	f = &DiskFree{
		total: float64(s.Blocks) * float64(s.Bsize),
		used:  float64(s.Blocks-s.Bfree) * float64(s.Bsize),
	}

	f.avail = f.total - f.used
	f.percentUsed = int(math.Ceil(f.used / f.total * 100))

	return f, err
}

// Total takes a value type (fs.MB, fs.GB etc) and returns the amount of total space
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Total(fs.GB))
func (df *DiskFree) Total(valType int) float64 {
	return df.total / float64(valType)
}

// Used takes a value type (fs.MB, fs.GB etc) and returns the amount of space used.
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Used(fs.GB))
func (df *DiskFree) Used(valType int) float64 {
	return df.used / float64(valType)
}

// Avail takes a value type (fs.MB, fs.GB etc) and returns the amount of space available.
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Avail(fs.GB))
func (df *DiskFree) Avail(valType int) float64 {
	return df.avail / float64(valType)
}

// PercentUsed returns the percentage of space used as an int
func (df *DiskFree) PercentUsed() int {
	return df.percentUsed
}

// BinPath returns the absolute path to the directory of the running binary
func BinPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// UID returns the Unix user id of the owner of file f
func UID(f os.FileInfo) int {
	return int(f.Sys().(*syscall.Stat_t).Uid)
}

// Uname returns the Unix username of the owner of file f
func Uname(f os.FileInfo) string {
	uid := strconv.Itoa(UID(f))
	u, _ := user.LookupId(uid)
	return u.Username
}

// YearsOld returns the number of years since file f modTime was last changed
func YearsOld(f os.FileInfo) float64 {
	return time.Now().Sub(f.ModTime()).Hours() / 24 / 365
}

// WeeksOld returns the number of weeks file f modTime was last changed
func WeeksOld(f os.FileInfo) float64 {
	return time.Now().Sub(f.ModTime()).Hours() / 24 / 7
}

// DaysOld returns the number of days since file f modTime was last changed
func DaysOld(f os.FileInfo) float64 {
	return time.Now().Sub(f.ModTime()).Hours() / 24
}
