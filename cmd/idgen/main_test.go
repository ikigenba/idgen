package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ai4mgreenly/idgen/internal/cli"
)

type testClock struct{}

func (testClock) Now() time.Time      { return time.Unix(0, 0) }
func (testClock) Sleep(time.Duration) {}

func TestUnstampedVersionDefaultsToDev(t *testing.T) {
	for _, arg := range []string{"--version", "-V"} {
		t.Run(arg, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(
				[]string{arg},
				strings.NewReader(""),
				&stdout,
				&stderr,
				version,
				testClock{},
			)

			// R-TJYQ-4NKS
			if code != 0 || stdout.String() != "dev\n" || stderr.String() != "" {
				t.Fatalf(
					"Run(%q) = (code %d, stdout %q, stderr %q), want (0, %q, %q)",
					arg,
					code,
					stdout.String(),
					stderr.String(),
					"dev\n",
					"",
				)
			}
		})
	}
}

func TestLinkerStampedVersionIsPrintedByBinary(t *testing.T) {
	const sentinel = "v9.8.7-linker-test"
	root := filepath.Clean(filepath.Join("..", ".."))
	binary := filepath.Join(t.TempDir(), "idgen")

	build := exec.Command(
		"go",
		"build",
		"-ldflags",
		"-X main.version="+sentinel,
		"-o",
		binary,
		"./cmd/idgen",
	)
	build.Dir = root
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build stamped binary: %v\n%s", err, output)
	}

	run := exec.Command(binary, "--version")
	var stdout, stderr bytes.Buffer
	run.Stdout = &stdout
	run.Stderr = &stderr
	err := run.Run()

	// R-TL6M-IFBH
	if err != nil || stdout.String() != sentinel+"\n" || stderr.String() != "" {
		t.Fatalf(
			"stamped binary = (err %v, stdout %q, stderr %q), want (nil, %q, %q)",
			err,
			stdout.String(),
			stderr.String(),
			sentinel+"\n",
			"",
		)
	}
}

func TestMakeBuildAndTestTargets(t *testing.T) {
	const childMarker = "IDGEN_MAKE_TEST_CHILD"
	if os.Getenv(childMarker) == "1" {
		return
	}

	root := filepath.Clean(filepath.Join("..", ".."))
	test := exec.Command("make", "-C", root, "test")
	test.Env = append(os.Environ(), childMarker+"=1")
	if output, err := test.CombinedOutput(); err != nil {
		t.Fatalf("make test: %v\n%s", err, output)
	}

	build := exec.Command("make", "-C", root, "-B", "build")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("make build: %v\n%s", err, output)
	}

	binary := filepath.Join(root, "bin", "idgen")
	info, err := os.Stat(binary)

	// R-XMM0-QR6E
	if err != nil {
		t.Fatalf("built binary %q: %v", binary, err)
	}
	if !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		t.Fatalf("built binary mode = %v, want executable regular file", info.Mode())
	}
}
