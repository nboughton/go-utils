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

// DiskFree packages the same types of information one gets from the *nix "df" cmd
type DiskFree struct {
	total       float64
	used        float64
	avail       float64
	percentUsed int
}

// sizeIncrement provides a specific type for data size category conversion
type sizeIncrement int64

// Exported constants for file size calculations
const (
	KB sizeIncrement = 1024    // Kilobyte
	MB               = KB * KB // Megabyte
	GB               = MB * KB // Gigabyte
	TB               = GB * KB // Terabyte
	PB               = TB * KB // Petabyte
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

// Total takes a value type (fs.MB, fs.GB etc) and returns the amount of total space
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Total(fs.GB))
func (df *DiskFree) Total(valType sizeIncrement) float64 {
	return df.total / float64(valType)
}

// Used takes a value type (fs.MB, fs.GB etc) and returns the amount of space used.
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Used(fs.GB))
func (df *DiskFree) Used(valType sizeIncrement) float64 {
	return df.used / float64(valType)
}

// Avail takes a value type (fs.MB, fs.GB etc) and returns the amount of space available.
// Example:
//    df, _ := fs.NewDf("/mnt/fs")
//    fmt.Println(df.Avail(fs.GB))
func (df *DiskFree) Avail(valType sizeIncrement) float64 {
	return df.avail / float64(valType)
}

// PercentUsed returns the percentage of space used as an int
func (df *DiskFree) PercentUsed() int {
	return df.percentUsed
}

// Mount packages the information gleaned from the *nix command "mount" for a single device
type Mount struct {
	Point  string
	Device string
	Type   string
	Args   []string
}

// Mounts returns an array of Mount structs, replicating the data gleaned from the *nix cmd "mount"
func Mounts() (m []Mount, err error) {
	o, err := exec.Command("mount").Output()
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(o), "\n") {
		f := strings.Fields(line)

		if len(f) == 6 {
			m = append(m, Mount{
				Point:  f[2],
				Device: f[0],
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
func Uname(f os.FileInfo) string {
	uid := strconv.Itoa(UID(f))
	u, _ := user.LookupId(uid)
	return u.Username
}

type ageIncrement float64

// Exported constants for file age calculation MONTH is omitted as months can
// have variable lengths. Technically years can vary but the difference is
// essentially a rounding error in the grand scheme of things.
const (
	HOUR ageIncrement = 1
	DAY               = HOUR * 24
	WEEK              = DAY * 7
	YEAR              = WEEK * 52
)

// Age returns the age of file f in the given increment valType
func Age(f os.FileInfo, valType ageIncrement) float64 {
	return time.Now().Sub(f.ModTime()).Hours() / float64(valType)
}
