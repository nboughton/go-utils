// Package common supplies a few very commonly used patterns found in CLI/Web apps
// namely case insensitive matching for yes/no questions and matching html tags
package common

import (
	"regexp"
)

// Exported vars
var (
	// Matchers yes/no strings
	YN  = regexp.MustCompile(`(?i)^(y(es|)|n(o|))$`)
	Yes = regexp.MustCompile(`(?i)^y(es|)$`)
	No  = regexp.MustCompile(`(?i)^n(o|)$`)
	// Matcher for HTML tags
	HTMLTag = regexp.MustCompile(`(<[^>]*>)`)
	// Matcher for ASCII chars
	ASCII = regexp.MustCompile(`[\x00-\x7F]`)
	// Matcher for ANSI escape chars (colors etc)
	ANSI = regexp.MustCompile(`\x1b[^m]*m`)
)
