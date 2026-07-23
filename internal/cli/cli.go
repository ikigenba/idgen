package cli

import (
	"io"

	"github.com/ai4mgreenly/idgen/internal/idgen"
)

// Run executes the command-line interface.
func Run(
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
	version string,
	clock idgen.Clock,
) int {
	_ = args
	_ = stdin
	_ = stdout
	_ = stderr
	_ = version
	_ = clock
	return 0
}
