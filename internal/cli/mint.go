package cli

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/ai4mgreenly/idgen/internal/idgen"
)

var validPrefix = regexp.MustCompile(`^[A-Za-z0-9]+$`)

func runMint(
	number int,
	prefix string,
	args []string,
	stdout, stderr io.Writer,
	clock Clock,
) int {
	if len(args) != 0 {
		fmt.Fprintf(stderr, "idgen: unexpected argument(s): %v\n", args)
		return exitUsage
	}
	if !validPrefix.MatchString(prefix) {
		fmt.Fprintf(stderr, "idgen: invalid prefix %q\n", prefix)
		return exitUsage
	}
	if number <= 0 {
		fmt.Fprintf(stderr, "idgen: --number must be > 0, got %d\n", number)
		return exitUsage
	}

	var lastMS int64 = -1
	for range number {
		var now time.Time
		for {
			now = clock.Now()
			ms := now.Sub(idgen.Epoch).Milliseconds()
			if ms > lastMS {
				lastMS = ms
				break
			}
			clock.Sleep(time.Millisecond)
		}
		fmt.Fprintln(stdout, idgen.MintAt(prefix, now))
	}
	return exitSuccess
}
