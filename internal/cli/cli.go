package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"time"
)

const (
	exitSuccess = 0
	exitFailure = 1
	exitUsage   = 2
)

// Clock supplies time and provides the wait seam used when minting a batch.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type realClock struct{}

func (realClock) Now() time.Time        { return time.Now() }
func (realClock) Sleep(d time.Duration) { time.Sleep(d) }

var (
	_ Clock = realClock{}
)

// Run executes the command-line interface.
func Run(
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
	version string,
	clock Clock,
) int {
	fs := flag.NewFlagSet("idgen", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		decode  bool
		number  int
		prefix  string
		showVer bool
	)
	fs.BoolVar(&decode, "decode", false, "decode identifiers")
	fs.IntVar(&number, "n", 1, "number of identifiers to mint")
	fs.IntVar(&number, "number", 1, "number of identifiers to mint")
	fs.StringVar(&prefix, "p", "R", "identifier prefix")
	fs.StringVar(&prefix, "prefix", "R", "identifier prefix")
	fs.BoolVar(&showVer, "version", false, "print version")
	fs.BoolVar(&showVer, "V", false, "print version")
	fs.Usage = func() {}

	switch err := fs.Parse(args); {
	case errors.Is(err, flag.ErrHelp):
		writeUsage(stdout)
		return exitSuccess
	case err != nil:
		writeUsage(stderr)
		return exitUsage
	}

	if showVer {
		fmt.Fprintln(stdout, version)
		return exitSuccess
	}
	if decode {
		return runDecode(fs.Args(), stdin, stdout, stderr)
	}
	return runMint(number, prefix, fs.Args(), stdout, stderr, clock)
}
