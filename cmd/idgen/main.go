package main

import (
	"os"
	"time"

	"github.com/ai4mgreenly/idgen/internal/cli"
)

var version = "dev"

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

func main() {
	os.Exit(cli.Run(
		os.Args[1:],
		os.Stdin,
		os.Stdout,
		os.Stderr,
		version,
		systemClock{},
	))
}
