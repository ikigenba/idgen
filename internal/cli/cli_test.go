package cli

import (
	"bytes"
	"io"
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
	return runForTestWithInput(args, strings.NewReader(""), clock)
}

func runForTestWithInput(args []string, stdin io.Reader, clock Clock) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := Run(args, stdin, &stdout, &stderr, "v-test", clock)
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

func TestDecodeFlagRoutesToDecodePath(t *testing.T) {
	instant := idgen.Epoch.Add(1500 * time.Millisecond)
	id := idgen.MintAt("R", instant)
	code, stdout, stderr := runForTest([]string{"--decode", id}, &fakeClock{})

	// R-X5JF-DYSO
	if code != exitSuccess || stdout != "2026-01-01T00:00:01.500Z\n" || stderr != "" {
		t.Fatalf("Run() = (%d, %q, %q), want decoded instant", code, stdout, stderr)
	}
}

func TestDecodeIgnoresMintFlags(t *testing.T) {
	instant := idgen.Epoch.Add(2345 * time.Millisecond)
	id := idgen.MintAt("R", instant)
	baseCode, baseStdout, baseStderr := runForTest([]string{"--decode", id}, &fakeClock{})
	code, stdout, stderr := runForTest(
		[]string{"--decode", "-n", "99", "-p", "IGNORED", id},
		&fakeClock{},
	)

	// R-X6RB-RQJD
	if code != baseCode || stdout != baseStdout || stderr != baseStderr {
		t.Fatalf(
			"decode with mint flags = (%d, %q, %q), want unchanged (%d, %q, %q)",
			code, stdout, stderr, baseCode, baseStdout, baseStderr,
		)
	}
}

func TestDecodePositionalsPreservesInputOrder(t *testing.T) {
	ids := []string{
		idgen.MintAt("ONE", idgen.Epoch.Add(5*time.Millisecond)),
		idgen.MintAt("TWO", idgen.Epoch.Add(2*time.Second)),
		idgen.MintAt("THREE", idgen.Epoch.Add(999*time.Millisecond)),
	}
	code, stdout, stderr := runForTest(
		append([]string{"--decode"}, ids...),
		&fakeClock{},
	)

	// R-X974-JA0R
	want := "" +
		"2026-01-01T00:00:00.005Z\n" +
		"2026-01-01T00:00:02.000Z\n" +
		"2026-01-01T00:00:00.999Z\n"
	if code != exitSuccess || stdout != want || stderr != "" {
		t.Fatalf("Run() = (%d, %q, %q), want ordered decode output %q", code, stdout, stderr, want)
	}
}

func TestDecodeStdinMixedWhitespaceMatchesPositionals(t *testing.T) {
	ids := []string{
		idgen.MintAt("A", idgen.Epoch.Add(10*time.Millisecond)),
		idgen.MintAt("B", idgen.Epoch.Add(20*time.Millisecond)),
		idgen.MintAt("C", idgen.Epoch.Add(30*time.Millisecond)),
	}
	positionalCode, positionalStdout, positionalStderr := runForTest(
		append([]string{"--decode"}, ids...),
		&fakeClock{},
	)
	stdin := strings.NewReader(ids[0] + "\t" + ids[1] + "\n \r\n" + ids[2])
	code, stdout, stderr := runForTestWithInput([]string{"--decode"}, stdin, &fakeClock{})

	// R-XAF0-X1RG
	if code != positionalCode || stdout != positionalStdout || stderr != positionalStderr {
		t.Fatalf(
			"stdin decode = (%d, %q, %q), want positional result (%d, %q, %q)",
			code, stdout, stderr, positionalCode, positionalStdout, positionalStderr,
		)
	}
}

type forbiddenReader struct {
	t *testing.T
}

func (r forbiddenReader) Read([]byte) (int, error) {
	r.t.Helper()
	r.t.Fatal("stdin Read called despite positional decode input")
	return 0, nil
}

func TestDecodePositionalsDoNotReadStdin(t *testing.T) {
	instant := idgen.Epoch.Add(42 * time.Millisecond)
	id := idgen.MintAt("R", instant)
	code, stdout, stderr := runForTestWithInput(
		[]string{"--decode", id},
		forbiddenReader{t: t},
		&fakeClock{},
	)

	// R-XBMX-ATI5
	if code != exitSuccess || stdout != "2026-01-01T00:00:00.042Z\n" || stderr != "" {
		t.Fatalf("Run() = (%d, %q, %q), want positional-only decode", code, stdout, stderr)
	}
}

func TestDecodeMalformedTokenContinuesBatch(t *testing.T) {
	first := idgen.MintAt("R", idgen.Epoch.Add(7*time.Millisecond))
	last := idgen.MintAt("R", idgen.Epoch.Add(9*time.Millisecond))
	const malformed = "not-an-id"
	code, stdout, stderr := runForTest(
		[]string{"--decode", first, malformed, last},
		&fakeClock{},
	)

	// R-XCUT-OL8U
	wantStdout := "2026-01-01T00:00:00.007Z\n2026-01-01T00:00:00.009Z\n"
	if code != exitFailure || stdout != wantStdout {
		t.Fatalf("Run() = (%d, stdout %q), want partial output %q and exit 1", code, stdout, wantStdout)
	}
	if !strings.Contains(stderr, malformed) {
		t.Fatalf("stderr = %q, want malformed token %q named", stderr, malformed)
	}
}

func TestDecodeEmptyInputIsVacuousSuccess(t *testing.T) {
	code, stdout, stderr := runForTestWithInput(
		[]string{"--decode"},
		strings.NewReader(""),
		&fakeClock{},
	)

	// R-XE2Q-2CZJ
	if code != exitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("Run() = (%d, %q, %q), want silent success", code, stdout, stderr)
	}
}

func TestMintThenDecodeRoundTripThroughRun(t *testing.T) {
	instant := idgen.Epoch.Add(987654 * time.Millisecond)
	mintCode, mintStdout, mintStderr := runForTest(nil, &fakeClock{now: instant})
	if mintCode != exitSuccess || mintStderr != "" {
		t.Fatalf("mint Run() = (%d, %q, %q), want success", mintCode, mintStdout, mintStderr)
	}

	id := strings.TrimSpace(mintStdout)
	code, stdout, stderr := runForTest([]string{"--decode", id}, &fakeClock{})

	// R-XFAM-G4Q8
	if code != exitSuccess || stdout != "2026-01-01T00:16:27.654Z\n" || stderr != "" {
		t.Fatalf("decode Run() = (%d, %q, %q), want minting instant", code, stdout, stderr)
	}
}

func TestDecodeOutputIsUTCRegardlessOfTZ(t *testing.T) {
	t.Setenv("TZ", "America/Chicago")
	instant := time.Date(2026, time.January, 1, 6, 2, 3, 4*int(time.Millisecond), time.UTC)
	id := idgen.MintAt("R", instant)
	code, stdout, stderr := runForTest([]string{"--decode", id}, &fakeClock{})

	// R-XIYB-LFYB
	if code != exitSuccess || stdout != "2026-01-01T06:02:03.004Z\n" || stderr != "" {
		t.Fatalf("Run() = (%d, %q, %q), want UTC output", code, stdout, stderr)
	}
}
