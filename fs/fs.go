// Package fs is for utility functions related to dealing with Linux filesystems
package fs

import (
	"math"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// DiskFree packages the same types of information one gets from the GNU Linux "df" cmd
type DiskFree struct {
	total       float64
	used        float64
	avail       float64
	percentUsed int
}

// size provides a specific type for data size category conversion
type size int64

// Exported constants for file size calculations
const (
	KB size = 1024    // Kilobyte
	MB      = KB * KB // Megabyte
	GB      = MB * KB // Gigabyte
	TB      = GB * KB // Terabyte
	PB      = TB * KB // Petabyte
)

// NewDf creates a DiskFree struct for the specified mount point. One can then use the
// struct methods to get the appropriate values in whatever size increment they wish.
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

// Total takes a size (fs.MB, fs.GB etc) and returns the amount of total space
//
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Total(fs.GB))
func (df *DiskFree) Total(s size) float64 {
	return df.total / float64(s)
}

// Used takes a size (fs.MB, fs.GB etc) and returns the amount of space used.
//
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Used(fs.GB))
func (df *DiskFree) Used(s size) float64 {
	return df.used / float64(s)
}

// Avail takes a size (fs.MB, fs.GB etc) and returns the amount of space available.
//
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Avail(fs.GB))
func (df *DiskFree) Avail(s size) float64 {
	return df.avail / float64(s)
}

// PercentUsed returns the percentage of space used as an int
func (df *DiskFree) PercentUsed() int {
	return df.percentUsed
}

// Mount packages the information gleaned from the GNU Linux command "mount" for a single device
type Mount struct {
	Point  string
	Device string
	Type   string
	Args   []string
}

// Mounts returns an array of Mount structs, replicating the data gleaned from the GNU Linux cmd "mount"
func Mounts() (m []Mount, err error) {
	o, err := exec.Command("mount").Output()
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(o), "\n") {
		f := strings.Fields(line)

		if len(f) == 6 {
			m = append(m, Mount{
				Device: f[0],
				Point:  f[2],
				Type:   f[4],
				Args:   strings.Split(strings.Trim(f[5], "()"), ","),
			})
		}
	}

	return m, nil
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
func Uname(f os.FileInfo) (string, error) {
	uid := strconv.Itoa(UID(f))
	u, err := user.LookupId(uid)
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

// duration provides a specific type for file age calculations
type duration float64

// Exported constants for file age calculation MONTH is omitted as months can
// have variable lengths. Technically years can vary but the difference is
// essentially a rounding error in the grand scheme of things.
const (
	HOURS duration = 1
	DAYS           = 24 * HOURS
	WEEKS          = 7 * DAYS
	YEARS          = 365.242199 * DAYS
)

// Age returns the age of file f in the given increment s.
//
// Example:
//    f, _ := os.Stat("/path/to/file")
//    fmt.Println(fs.Age(f, fs.WEEK))
func Age(f os.FileInfo, d duration) float64 {
	return time.Since(f.ModTime()).Hours() / float64(d)
}
