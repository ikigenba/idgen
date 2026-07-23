package cli

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ai4mgreenly/idgen/internal/idgen"
)

type fakeClock struct {
	now    time.Time
	sleeps []time.Duration
}

func (c *fakeClock) Now() time.Time { return c.now }

func (c *fakeClock) Sleep(d time.Duration) {
	c.sleeps = append(c.sleeps, d)
	c.now = c.now.Add(d)
}

func runForTest(args []string, clock Clock) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := Run(args, strings.NewReader(""), &stdout, &stderr, "v-test", clock)
	return code, stdout.String(), stderr.String()
}

func outputIDs(t *testing.T, output string, count int) []string {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if len(lines) != count {
		t.Fatalf("got %d output lines, want %d: %q", len(lines), count, output)
	}
	return lines
}

func decodedTimes(t *testing.T, ids []string) []time.Time {
	t.Helper()
	times := make([]time.Time, len(ids))
	for i, id := range ids {
		instant, err := idgen.TimeOf(id)
		if err != nil {
			t.Fatalf("TimeOf(%q): %v", id, err)
		}
		times[i] = instant
	}
	return times
}

func TestMintBatchUsesDistinctMilliseconds(t *testing.T) {
	clock := &fakeClock{now: idgen.Epoch.Add(time.Second)}
	code, stdout, stderr := runForTest([]string{"-n", "4"}, clock)

	// R-WTCF-K9DQ
	if code != exitSuccess || stderr != "" {
		t.Fatalf("Run() = (%d, stderr %q), want success", code, stderr)
	}
	ids := outputIDs(t, stdout, 4)
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate minted id %q in %v", id, ids)
		}
		seen[id] = true
	}

	// R-WUKB-Y14F
	var elapsed time.Duration
	for _, sleep := range clock.sleeps {
		elapsed += sleep
	}
	if elapsed < 3*time.Millisecond {
		t.Fatalf("virtual time advanced %v, want at least 3ms", elapsed)
	}
}

func TestMintBatchTerminatesWhenSleepAdvancesStalledClock(t *testing.T) {
	clock := &fakeClock{now: idgen.Epoch.Add(10 * time.Second)}
	code, stdout, stderr := runForTest([]string{"--number", "5"}, clock)

	// R-WVS8-BSV4
	if code != exitSuccess || stderr != "" {
		t.Fatalf("Run() = (%d, stderr %q), want success", code, stderr)
	}
	ids := outputIDs(t, stdout, 5)
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("minted ids are not pairwise distinct: %v", ids)
		}
		seen[id] = true
	}
	times := decodedTimes(t, ids)
	if advance := times[len(times)-1].Sub(times[0]); advance < 4*time.Millisecond {
		t.Fatalf("last id advances %v beyond first, want at least 4ms", advance)
	}
}

func TestSingleMintDoesNotSleep(t *testing.T) {
	for _, test := range []struct {
		name string
		args []string
	}{
		{name: "default"},
		{name: "explicit one", args: []string{"-n", "1"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			clock := &fakeClock{now: idgen.Epoch.Add(time.Second)}
			code, stdout, stderr := runForTest(test.args, clock)

			// R-WX04-PKLT
			if code != exitSuccess || stderr != "" || stdout == "" {
				t.Fatalf("Run() = (%d, stdout %q, stderr %q), want one successful mint", code, stdout, stderr)
			}
			if len(clock.sleeps) != 0 {
				t.Fatalf("Sleep called %d times, want zero", len(clock.sleeps))
			}
		})
	}
}

func TestMintUsesJustReadInstant(t *testing.T) {
	now := idgen.Epoch.Add(123456 * time.Millisecond)
	clock := &fakeClock{now: now}
	code, stdout, stderr := runForTest(nil, clock)

	// R-WY81-3CCI
	if code != exitSuccess || stderr != "" {
		t.Fatalf("Run() = (%d, stderr %q), want success", code, stderr)
	}
	minted := outputIDs(t, stdout, 1)[0]
	got, err := idgen.TimeOf(minted)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(now) {
		t.Fatalf("minted time = %v, want just-read time %v", got, now)
	}
}

type backwardClock struct {
	now      time.Time
	nowCalls int
}

func (c *backwardClock) Now() time.Time {
	c.nowCalls++
	if c.nowCalls == 2 {
		c.now = c.now.Add(-5 * time.Millisecond)
	}
	return c.now
}

func (c *backwardClock) Sleep(d time.Duration) { c.now = c.now.Add(d) }

func TestMintWaitsOutBackwardClockStep(t *testing.T) {
	clock := &backwardClock{now: idgen.Epoch.Add(time.Second)}
	code, stdout, stderr := runForTest([]string{"-n", "2"}, clock)

	// R-WZFX-H437
	if code != exitSuccess || stderr != "" {
		t.Fatalf("Run() = (%d, stderr %q), want success", code, stderr)
	}
	times := decodedTimes(t, outputIDs(t, stdout, 2))
	if !times[1].After(times[0]) {
		t.Fatalf("times did not strictly cross backward dip: %v", times)
	}
	if clock.nowCalls < 7 {
		t.Fatalf("Now called %d times, backward excursion was not waited out", clock.nowCalls)
	}
}

func TestHelpPrintsUsageOnceToStdout(t *testing.T) {
	for _, arg := range []string{"--help", "-h"} {
		t.Run(arg, func(t *testing.T) {
			code, stdout, stderr := runForTest([]string{arg}, &fakeClock{})

			// R-X0NT-UVTW
			if code != exitSuccess {
				t.Fatalf("code = %d, want 0", code)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}
			if count := strings.Count(stdout, "Usage:"); count != 1 {
				t.Fatalf("Usage occurrence count = %d, want exactly 1 in %q", count, stdout)
			}
		})
	}
}

func TestVersionFlagsPrintInjectedVersion(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		id   string
	}{
		{name: "long", arg: "--version", id: "R-X1VQ-8NKL"},
		{name: "short", arg: "-V", id: "R-TIQT-QVU3"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, stdout, stderr := runForTest([]string{test.arg}, &fakeClock{})
			if test.id == "R-X1VQ-8NKL" {
				// R-X1VQ-8NKL
				if code != exitSuccess || stdout != "v-test\n" || stderr != "" {
					t.Fatalf("Run() = (%d, %q, %q), want (0, version, empty)", code, stdout, stderr)
				}
			} else {
				// R-TIQT-QVU3
				if code != exitSuccess || stdout != "v-test\n" || stderr != "" {
					t.Fatalf("Run() = (%d, %q, %q), want (0, version, empty)", code, stdout, stderr)
				}
			}
		})
	}
}

func TestUnknownFlagIsUsageError(t *testing.T) {
	code, _, stderr := runForTest([]string{"--wat"}, &fakeClock{})

	// R-X33M-MFBA
	if code != exitUsage || stderr == "" {
		t.Fatalf("Run() = (%d, stderr %q), want usage error with message", code, stderr)
	}
}

func TestMintRejectsPositionalArgument(t *testing.T) {
	code, _, stderr := runForTest([]string{"surprise"}, &fakeClock{})

	// R-X4BJ-071Z
	if code != exitUsage || !strings.Contains(stderr, "surprise") {
		t.Fatalf("Run() = (%d, stderr %q), want usage error naming argument", code, stderr)
	}
}

func TestBareInvocationMintsDefaultPrefix(t *testing.T) {
	clock := &fakeClock{now: idgen.Epoch.Add(time.Second)}
	code, stdout, stderr := runForTest(nil, clock)

	// R-7UL7-PF0O
	if code != exitSuccess || stderr != "" {
		t.Fatalf("Run() = (%d, stderr %q), want success", code, stderr)
	}
	if matched := regexp.MustCompile(`^R-[0-9A-Z]{4}-[0-9A-Z]{4}\n$`).MatchString(stdout); !matched {
		t.Fatalf("stdout = %q, want exactly one default-prefix id line", stdout)
	}
}

func TestCustomPrefixReplacesDefault(t *testing.T) {
	for _, prefix := range []string{"X", "TEAM"} {
		t.Run(prefix, func(t *testing.T) {
			clock := &fakeClock{now: idgen.Epoch.Add(time.Second)}
			code, stdout, stderr := runForTest([]string{"-p", prefix}, clock)

			// R-PU67-68HE
			want := regexp.MustCompile(`^` + prefix + `-[0-9A-Z]{4}-[0-9A-Z]{4}\n$`)
			if code != exitSuccess || stderr != "" || !want.MatchString(stdout) {
				t.Fatalf("Run() = (%d, stdout %q, stderr %q), want custom-prefix id", code, stdout, stderr)
			}
		})
	}
}

func TestInvalidPrefixIsUsageError(t *testing.T) {
	for _, prefix := range []string{"", "  ", "R-X", "S/"} {
		t.Run(prefix, func(t *testing.T) {
			code, _, stderr := runForTest([]string{"-p", prefix}, &fakeClock{})

			// R-XGII-TWGX
			if code != exitUsage || !strings.Contains(stderr, "invalid prefix") {
				t.Fatalf("Run() = (%d, stderr %q), want invalid-prefix usage error", code, stderr)
			}
		})
	}
}

func TestNonPositiveNumberIsUsageError(t *testing.T) {
	for _, number := range []string{"0", "-3"} {
		t.Run(number, func(t *testing.T) {
			code, _, stderr := runForTest([]string{"-n", number}, &fakeClock{})

			// R-XHQF-7O7M
			if code != exitUsage || !strings.Contains(stderr, "--number must be > 0") {
				t.Fatalf("Run() = (%d, stderr %q), want number usage error", code, stderr)
			}
		})
	}
}

func TestHelpMentionsPrimaryFlags(t *testing.T) {
	code, stdout, stderr := runForTest([]string{"--help"}, &fakeClock{})

	// R-XLE4-CZFP
	if code != exitSuccess || stderr != "" {
		t.Fatalf("Run() = (%d, stderr %q), want success", code, stderr)
	}
	for _, flag := range []string{"-n", "-p", "--decode"} {
		if !strings.Contains(stdout, flag) {
			t.Errorf("usage does not mention %q: %q", flag, stdout)
		}
	}
}
