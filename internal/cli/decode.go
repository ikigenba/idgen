package cli

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ai4mgreenly/idgen/internal/idgen"
)

func runDecode(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	ids := args
	if len(ids) == 0 {
		scanner := bufio.NewScanner(stdin)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			ids = append(ids, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(stderr, "idgen: reading input: %v\n", err)
			return exitFailure
		}
	}

	exit := exitSuccess
	for _, id := range ids {
		instant, err := idgen.TimeOf(id)
		if err != nil {
			fmt.Fprintf(stderr, "idgen: %q: %v\n", id, err)
			exit = exitFailure
			continue
		}
		fmt.Fprintln(stdout, instant.UTC().Format("2006-01-02T15:04:05.000Z"))
	}
	return exit
}
