package cli

import (
	"fmt"
	"io"
)

func writeUsage(w io.Writer) {
	fmt.Fprintln(w, `Usage: idgen [options] [ID ...]

Mint an identifier using the current time by default.

Options:
  -n, --number N       mint N identifiers (default 1)
  -p, --prefix PREFIX  use PREFIX (default "R")
      --decode         decode ID arguments, or whitespace-delimited IDs from stdin
  -h, --help           print this help
  -V, --version        print version`)
}
