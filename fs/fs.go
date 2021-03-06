// Package fs is for utility functions related to dealing with Linux filesystems
package fs

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

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

// AbsPath returns the absolute path to the directory of the running binary
func AbsPath() (string, error) {
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

// GID returns the GID of the file
func GID(f os.FileInfo) int {
	return int(f.Sys().(*syscall.Stat_t).Gid)
}

// Gname returns the human readable group name of a file
func Gname(f os.FileInfo) (string, error) {
	gid := strconv.Itoa(GID(f))
	g, err := user.LookupGroupId(gid)
	if err != nil {
		return "", err
	}
	return g.Name, nil
}

// IsSymlink tests a filepath to see if it is a symlink or not. Returns err if
// it cannot stat the file
func IsSymlink(path string) (bool, error) {
	f, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	if f.Mode()&os.ModeSymlink != 0 {
		return true, nil
	}

	return false, nil
}

// Age returns the age of file f in the given increment s.
//
// Example:
//    f, _ := os.Stat("/path/to/file")
//    fmt.Println(fs.Age(f, fs.WEEK))
func Age(f os.FileInfo, d duration) (aTime, mTime, cTime float64) {
	stat := f.Sys().(*syscall.Stat_t)

	aTime = time.Since(time.Unix(stat.Atim.Sec, stat.Atim.Nsec)).Hours() / float64(d)
	mTime = time.Since(f.ModTime()).Hours() / float64(d)
	cTime = time.Since(time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)).Hours() / float64(d)

	return
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
