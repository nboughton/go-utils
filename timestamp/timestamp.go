// Package timestamp contains utility functions for converting/creating common
// timestamp formats used in web development
package timestamp

import (
	"time"
)

// JS returns a Unix timestamp in milliseconds
func JS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// JSToUnix converts a unix timestamp (in milliseconds) to a standard unix date
func JSToUnix(tsMS int64) string {
	return time.Unix(0, tsMS*int64(time.Millisecond)).Format(time.UnixDate)
}

// JSToYYYYMMDD converts a timestamp (in milliseconds) to a date formatted YYYY-MM-DD
func JSToYYYYMMDD(tsMS int64) string {
	return time.Unix(0, tsMS*int64(time.Millisecond)).Format("2006-01-02")
}
