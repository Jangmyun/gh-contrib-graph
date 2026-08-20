package render

import (
	"os"

	"golang.org/x/term"
)

// defaultTerminalWidth is used when stdout isn't a TTY (e.g. piped output)
// or the size can't be determined.
const defaultTerminalWidth = 80

// TerminalWidth returns the width of the terminal attached to stdout, or
// defaultTerminalWidth if it can't be determined.
func TerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return defaultTerminalWidth
	}
	return w
}
