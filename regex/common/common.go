// Package common supplies a few very commonly used patterns found in CLI/Web apps
// namely case insensitive matching for yes/no questions and matching html tags
package common

import (
	"regexp"
)

// Exported vars all common regex patterns
var (
	YN      = regexp.MustCompile(`(?i)^(y(es|)|n(o|))$`)
	Yes     = regexp.MustCompile(`(?i)^y(es|)$`)
	No      = regexp.MustCompile(`(?i)^n(o|)$`)
	HTMLTag = regexp.MustCompile(`(<[^>]*>)`)
	ASCII   = regexp.MustCompile(`[\x00-\x7F]`)
	ANSI    = regexp.MustCompile(`\x1b[^m]*m`)
)
