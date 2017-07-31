// Package commonregex supplies a few very commonly used patterns found in CLI apps
// namely case insensitive matching for yes/no questions
package commonregex

import (
	"regexp"
)

// Exported vars all common regex patterns
var (
	YN  = regexp.MustCompile(`(?i)^[yn]$`)
	Yes = regexp.MustCompile(`(?i)^y(es|)$`)
	No  = regexp.MustCompile(`(?i)^n(o|)$`)
)
