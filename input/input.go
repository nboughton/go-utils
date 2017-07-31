// Package input is for reading and parsing input from the CLI
package input

import (
	"bufio"
	"os"
	"strings"
)

// ReadLine reads a single line of input from the CLI (delimited by a newline char) and
// returns the trimmed string
func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
