# idgen — Design

**Authority: shape and its proof.** This document owns *how idgen is built* —
seams, public interfaces, naming, types, the data model, the encoding math — and
*how each behavior is proven*. `docs/idgen-product.md` owns the *why* and the
user-facing promises; design states the **exact, checkable form** of those
promises and never re-declares the why. The contractual constants (the 2026 epoch
and the version string) are the product's; design uses their values, it does not
own them.

This is the **single, current** statement of the architecture: when a decision
changes, this document is rewritten to stay true (stale decisions are removed, not
stacked). History of how it got here lives in `docs/idgen-plan.md`.

## Verification & "done" — the denominator

Each Decision below ends with a **Verification** list: the concrete behaviors a
test must assert for that decision to be considered built. Every item carries a
minted **idgen id** (`R-XXXX-XXXX`) — a stable, unique handle for that one
requirement. **That set of lists is the denominator** — the enumerated intent the
test suite is measured against. There is **no separate requirements document**:
the ids live inline here and nowhere else. A behavior is **covered** when a test
asserts it *and names its id in a `// R-XXXX-XXXX` comment*, so coverage is a
grep, not a separate cross-reference. The work is **done** when every Verification
id is covered and `go test -race ./...` is green. The denominator stays honest the
same way the rest of the design does — by being pruned: remove a requirement here
and you remove its id and its test.

## Conventions

- Go 1.26, module path `github.com/ai4mgreenly/idgen`.
- Exit codes: `0` success, `1` decode data failure, `2` usage/runtime error.
- Every printed time is UTC, formatted `2006-01-02T15:04:05.000Z`.
- **Time source.** Standard-library `time` only — no third-party dependency, and
  millisecond precision is portable to every target Go compiles for (`time.Now()`
  resolves far finer than a millisecond on every supported platform). `Epoch` is a
  constructed `time.Date(..., time.UTC)` value carrying **no** monotonic reading,
  so `time.Now().Sub(Epoch)` strips monotonic and yields a pure **wall-clock**
  elapsed duration. That is what an ID needs: an absolute civil instant decodable
  forever, not host uptime (monotonic time has no fixed zero across reboots and
  could never anchor a decodable ID). A backward wall-clock step (NTP, manual
  change) is tolerated, not detected — see Decision 3.

---

## Decision 1 — Module path & package layout (top-level seams)

**Decision.** Three seams:

```
github.com/ai4mgreenly/idgen          (go.mod, go 1.26)
├── cmd/idgen/main.go                  thin: os.Args/stdio → cli.Run(); os.Exit
├── internal/cli/                       the testable CLI core
│   ├── cli.go      Run(args, stdin, stdout, stderr, clock) int  (exported)
│   ├── mint.go     mint path (default + -n + -p)
│   └── decode.go   decode path (--decode)
└── internal/idgen/                     pure encode/decode + Epoch + Clock seam
    └── idgen.go
```

- **`idgen`** — pure core: affine bijection, base-36 encoding, `Epoch`,
  constants, mint/decode. No I/O, no flags.
- **`cli`** — all flag parsing, stdin reading, stderr reporting, and exit codes
  sit behind one **exported** function `Run(args, stdin, stdout, stderr, clock) int`
  (exported because `package main` in `cmd/idgen` must call it across the seam).
- **`main`** — wires real `os` values into `cli.Run` and calls `os.Exit`; trivial.

**Rejected.**
- *Single `main` package (everything in `cmd/idgen`).* That suits a subcommand
  buried in a large binary; for a focused standalone tool an injectable `Run()`
  in `internal/cli` makes the whole CLI unit-testable for almost no cost.
- *Public `pkg/idgen`.* Product is install-from-source, single-purpose; keep the
  API private under `internal/` and free to change. Easy to promote later.

**Verification.** This decision is a seam choice; it has no behaviors of its own —
its payoff is *testability*, realized by the behavioral requirements of Decisions
2–6. Stated as the testing posture each seam buys: `idgen` proves its requirements
with zero process setup (unit + fuzz); `cli` proves its requirements through
in-memory `args`/`stdin`/`stdout` buffers + return code, no subprocess and no real
stdio; `main` has no logic and so carries no requirement.

---

## Decision 2 — `idgen` public API & prefix placement

**Decision.** `idgen` owns the format grammar in both directions and stays a
pure function of its inputs (no clock). The prefix is a parameter on mint and is
stripped-and-ignored on decode.

```go
package idgen

// Epoch is the zero point: 2026-01-01 00:00:00 UTC, constructed with
// time.Date(..., time.UTC) so it carries no monotonic reading.
var Epoch = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

// ErrInvalidID wraps the error TimeOf returns for a malformed id.
var ErrInvalidID = errors.New("invalid id")

// MintAt returns "<prefix>-XXXX-XXXX" for the given instant. Instants
// before Epoch are clamped to Epoch. The caller guarantees prefix is a
// non-empty run of letters/digits (cli validates); MintAt does not
// re-validate.
func MintAt(prefix string, t time.Time) string

// TimeOf inverts the body of any "<prefix>-XXXX-XXXX" id to the instant
// it was minted from. The prefix is accepted but ignored — the instant
// lives entirely in the body — so an id of any prefix decodes. Returns
// an error wrapping ErrInvalidID when id is not canonical.
func TimeOf(id string) (time.Time, error)
```

Body math, identical to the reference except the epoch:
`n = (ms·0x9E3779B1 + 0xC0FFEE) mod 36⁸` via `math/big`, uppercase base-36,
zero-padded to 8, split 4-4. `TimeOf` inverts with the precomputed modular
inverse of the multiplier; package `init` panics if multiplier and `36⁸` ever
lose coprimality (fail-loud — every existing id would otherwise be
irrecoverable). Decode grammar: `^[A-Za-z0-9]+-([0-9A-Z]{4})-([0-9A-Z]{4})$`.

- Prefix is a **parameter**, not baked in: `idgen` knows the grammar (body,
  base-36, dash split); the prefix string is supplied by the caller.
- Prefix **validation lives in `cli`** (validate at the flag boundary). A prefix
  containing `-` would corrupt the decode grammar; cli's letters/digits-only
  rule prevents it. `MintAt` documents and trusts the precondition.
- `idgen` holds **no clock** — the instant is always passed in; the `-n` wait
  concern moves to `cli`.

**Rejected.**
- *`idgen` deals only in the 8-char body; `cli` glues on `prefix-` and the
  dashes.* Leaks base-36/4-4 format into `cli` and splits decode-parsing across
  packages.
- *A `New`/`Generator` clock-carrying API.* Clock moves out to `cli`, so
  `idgen` needs neither. Kept `TimeOf` (reads as "the time of this id") over
  `Decode`.

**Verification.**
- `R-WH5F-QJYS` — Round-trip property/fuzz: `TimeOf(MintAt(p, t)) == t` for `t` in
  `[Epoch, Epoch+cycle)` truncated to ms, across arbitrary valid prefixes.
- `R-WIDC-4BPH` — Golden vector at `Epoch`: `Epoch` → its exact `R-XXXX-XXXX`
  string, derived independently/offline (locks the 2026 epoch; an epoch regression
  fails this).
- `R-WJL8-I3G6` — Golden vector mid-cycle: a known mid-cycle instant → its exact
  string, independently derived (locks the affine constants and the 4-4 split).
- `R-WKT4-VV6V` — Padding: small ms values still yield 8 body chars.
- `R-WM11-9MXK` — Clamping: instants before Epoch encode as Epoch (ms 0).
- `R-WN8X-NEO9` — Prefix-agnostic decode: ids with prefixes `R`, `S`, `SPEC` all
  decode to the same instant for the same body.
- `R-WPOQ-EY5N` — Malformed input: bad bodies/shapes return an error wrapping
  `ErrInvalidID`.
- `R-WQWM-SPWC` — `FuzzTimeOf`: arbitrary strings never panic — only `ErrInvalidID`
  or a valid time.
- `R-WS4J-6HN1` — Coprimality fail-loud: package `init` panics if the multiplier
  and `36⁸` ever lose coprimality (every existing id would otherwise be
  irrecoverable).

---

## Decision 3 — `Clock` seam & the `-n` wait loop

**Decision.** A small `Clock` interface, defined in `internal/cli` (its only
consumer) and injected into `Run`. Production wiring is `time.Now`/`time.Sleep`;
tests inject a fake.

```go
package cli

type Clock interface {
    Now() time.Time
    Sleep(d time.Duration)
}

type realClock struct{}
func (realClock) Now() time.Time        { return time.Now() }
func (realClock) Sleep(d time.Duration) { time.Sleep(d) }

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, clk Clock) int
```

Mint loop (driven through `clk`):

```go
var lastMs int64 = -1
for i := 0; i < count; i++ {
    var now time.Time
    for {
        now = clk.Now()
        ms := now.Sub(idgen.Epoch).Milliseconds()
        if ms > lastMs {            // strictly later ms than previous mint
            lastMs = ms
            break
        }
        clk.Sleep(time.Millisecond)
    }
    fmt.Fprintln(stdout, idgen.MintAt(prefix, now))
}
```

- **Mint from an already-elapsed instant**: mint the just-read `now`, never a
  future instant.
- **Distinct ms per id within a call**: the `lastMs` gate + `Sleep` put each id
  on its own strictly-later ms; `N` ids cost ≥ ~N−1 ms. `N=1`/default waits
  zero (the `lastMs == -1` gate admits immediately).
- **Backward-clock policy: tolerate.** A backward step yields `ms ≤ lastMs` and
  the loop *waits* until the clock climbs past `lastMs` — the minted sequence
  within one invocation never goes backward, with no special-case code. No
  cross-invocation protection (inherent to wall-clock ids). We document, not
  clamp/detect.

**Rejected.**
- *Two func fields `now`/`sleep`.* Works, but the interface bundles the
  shared-state `Now`+`Sleep` pair more cleanly.
- *Clamp/detect backward jumps.* Adds policy the product doesn't ask for;
  wait-until-caught-up is the simplest correct behavior.

**Verification.**
- `R-WTCF-K9DQ` — Fake `Clock` whose `Sleep(d)` advances virtual `now` by `d`:
  `-n N` prints N **distinct** ids (loop terminates with no real wall time).
- `R-WUKB-Y14F` — Under that same fake clock, virtual time advanced ≥ N−1 ms —
  proves the per-id wait actually happened.
- `R-WVS8-BSV4` — Stalled-clock: a `Clock` returning a fixed `Now` between sleeps
  still terminates and yields distinct ids only after `Sleep` advances it (locks
  "does not return before the last id's ms elapsed").
- `R-WX04-PKLT` — Default/`N=1`: zero `Sleep` calls (no wasted wait).
- `R-WY81-3CCI` — Mint from an already-elapsed instant: the id decodes to the
  just-read `now`, never a future instant.
- `R-WZFX-H437` — Backward-clock tolerate: with a clock that steps backward then
  recovers, the minted sequence within one invocation is non-decreasing in ms.

---

## Decision 4 — CLI grammar, dispatch & exit-code taxonomy

**Decision.** idgen is a single command with flags (no subcommands); `--decode`
flips mode. One `flag.FlagSet` is parsed in `Run`, then branched.

```go
fs := flag.NewFlagSet("idgen", flag.ContinueOnError)
fs.SetOutput(stderr)
var (
    decode  bool   // --decode
    number  int    // -n / --number, default 1
    prefix  string // -p / --prefix, default "R"
    showVer bool   // --version
)
// -n/--number and -p/--prefix each bound twice (long + short).
fs.Usage = func() { writeUsage(stderr) }

switch err := fs.Parse(args); {
case errors.Is(err, flag.ErrHelp): writeUsage(stdout); return exitSuccess // --help/-h
case err != nil:                   return exitUsage                       // bad flags → 2
}
if showVer { fmt.Fprintln(stdout, version); return exitSuccess }
if decode  { return runDecode(fs.Args(), stdin, stdout, stderr) }
return runMint(number, prefix, fs.NArg(), stdout, stderr, clk)
```

- **Stdlib `flag`**, no third-party CLI lib (research: no external deps). `flag`
  requires flags before positionals — matches the product's `idgen --decode
  R-...` examples.
- **`--help`/`-h`** → usage to **stdout**, exit 0. **`--version`** → bare version
  string to stdout, exit 0 (version string itself: Decision 6).
- **Mint takes no positionals**: `runMint` errors (exit 2) if `NArg() != 0`.
- **`-n`/`-p` in decode mode: accepted but ignored** (mint concerns, no decode
  meaning; keeps the grammar simple).
- **Exit codes:** `0` success · `2` usage error (bad flags, empty/invalid
  prefix, `number ≤ 0`, mint-with-positionals) · `1` decode data failure (≥1
  malformed id in an otherwise valid invocation). Constants `exitSuccess=0`,
  `exitFailure=1`, `exitUsage=2`.

**Rejected.**
- *Subcommands (`idgen mint`/`idgen decode`).* Contradicts the documented
  default-action-plus-`--decode` UX.
- *Single non-zero code for everything.* The 1-vs-2 split
  separates "called it wrong" from "a batch value was bad" at no cost.
- *Error on `-n`/`-p` in decode mode.* Gold-plating the product doesn't ask for.

**Verification.** Table-driven `Run` tests over arg vectors assert (stdout, stderr,
exit code):
- `R-X0NT-UVTW` — `--help`/`-h` → usage on **stdout**, exit 0.
- `R-X1VQ-8NKL` — `--version` → version string on **stdout**, exit 0.
- `R-X33M-MFBA` — unknown flag → exit 2.
- `R-X4BJ-071Z` — mint with a positional argument → exit 2.
- `R-X5JF-DYSO` — `--decode` routes to the decode path.
- `R-X6RB-RQJD` — `-n`/`-p` supplied with `--decode` are inert (decode output
  unchanged).

---

## Decision 5 — Input handling & validation (both modes)

**Decision (decode, `runDecode`).**

```go
func runDecode(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
    ids := args
    if len(ids) == 0 {                 // stdin only when no positionals
        sc := bufio.NewScanner(stdin)
        sc.Split(bufio.ScanWords)      // spaces, tabs, newlines all delimit
        for sc.Scan() { ids = append(ids, sc.Text()) }
    }
    exit := exitSuccess
    for _, id := range ids {
        t, err := idgen.TimeOf(id)
        if err != nil {
            fmt.Fprintf(stderr, "idgen: %s\n", err) // err already names the id
            exit = exitFailure                        // 1
            continue
        }
        fmt.Fprintln(stdout, t.UTC().Format("2006-01-02T15:04:05.000Z"))
    }
    return exit
}
```

- **Positional args take precedence; stdin is read only when there are none.**
  Each input yields exactly one line — a UTC instant on stdout or an error on
  stderr — in input order within each stream.
- **stdin tokenized by `bufio.ScanWords`** (spaces/tabs/newlines).
- **Partial failure tolerated**: a bad id reports to stderr (the `ErrInvalidID`
  message quotes the token) and the batch continues; exit becomes `1` if any
  failed.
- **No ids at all** (no args, empty stdin) → exit `0`, no output (vacuous
  success). *Rejected*: treating empty decode as usage error `2`.

**Decision (mint validation, `runMint`, before the loop).**
- **prefix** must match `^[A-Za-z0-9]+$` (non-empty, letters/digits only). Empty,
  whitespace, or any separator → `idgen: invalid prefix ...` on stderr, exit
  `2`. (A separator would corrupt the decode grammar; this is what guards it.)
- **number** must be `> 0` → else `idgen: --number must be > 0, got N`, exit `2`.

**Rejected.**
- *`strings.Fields` after slurping all stdin.* `bufio.Scanner` streams and is the
  idiomatic token reader; same splitting, no need to read it all into memory.
- *Validating prefix inside `idgen`.* Decision 2 placed validation at the flag
  boundary; `MintAt` trusts the precondition.

**Verification.**
- `R-X974-JA0R` — Decode from positional args: each positional yields one UTC line
  on stdout in input order.
- `R-XAF0-X1RG` — Decode from stdin (mixed whitespace) yields output identical to
  the positional-args case.
- `R-XBMX-ATI5` — Positionals win when both positionals and stdin are present
  (stdin is read only when there are no positionals).
- `R-XCUT-OL8U` — Batch with one malformed token: good ids still decode (in order),
  the error on stderr names the bad token, exit `1`.
- `R-XE2Q-2CZJ` — Empty decode (no args, empty stdin): no output, exit `0`.
- `R-XFAM-G4Q8` — Round-trip through `Run`: `--decode` of a freshly minted id
  returns the minting instant (ties Decisions 2/3/5 together end-to-end).
- `R-XGII-TWGX` — Prefix validation: `""`, `"  "`, `"R-X"`, `"S/"` → exit `2`.
- `R-XHQF-7O7M` — Number validation: `0`, `-3` → exit `2`.
- `R-XIYB-LFYB` — Decode output is UTC regardless of `TZ` env (test sets
  `TZ=America/Chicago`).

---

## Decision 6 — Version, usage text & Makefile

**Decision.** Version is a single in-source `var`, no build-time injection:

```go
// internal/cli/version.go
package cli
// version is the single source of truth, bumped by editing this line.
var version = "0.1.0-pre+20260616"
```

`--version` prints it bare. **Usage text** lives in `internal/cli/usage.go`
(`func writeUsage(io.Writer)`), to stdout for `--help`/`-h`, to stderr on a usage
error; covers default mint, `-n/--number`, `-p/--prefix`, `--decode`
(args-or-stdin note), `--help`, `--version`.

**Makefile** — exactly the product's five targets:

```make
BIN := bin/idgen
build:   ; go build -o $(BIN) ./cmd/idgen
test:    ; go test ./...
clean:   ; rm -rf bin
fmt:     ; gofmt -w .
install: ; go install ./cmd/idgen
.PHONY: build test clean fmt install
```

**Rejected.**
- *`-ldflags -X` version injection / `VERSION` file.* Two sources of truth and
  build coupling, with no provenance need for a source-only dev tool.
- *A separate `internal/version` package.* Over-engineered for one string.

**Verification.**
- `R-XK67-Z7P0` — `--version` asserts stdout is exactly `0.1.0-pre+20260616` — the
  product's success criterion; pins the constant.
- `R-XLE4-CZFP` — `--help` asserts usage mentions each flag (`-n`, `-p`,
  `--decode`).
- `R-XMM0-QR6E` — Build smoke: `go test ./...` is green and `go build` produces
  `bin/idgen` (the Makefile's `build`/`test` targets exercised).

---

## Decision 7 — Overall testing strategy & test layout

**Decision.** Three tiers, all in-process and deterministic (no real sleeps, no
subprocesses) — the `idgen`-pure / `cli.Run`-injectable seams make every
behavior reachable without launching a binary.

**1. `idgen` unit + fuzz** (`internal/idgen/idgen_test.go`, `fuzz_test.go`):
- **Golden vectors** — known instants → exact `R-XXXX-XXXX`, **derived
  independently/offline** (not snapshotted from the code under test) so they
  guard the *2026 epoch + affine constants*; an epoch regression breaks them.
  At minimum: `Epoch`→its id and one mid-cycle instant.
- Padding (small ms → 8 chars), clamping (pre-Epoch → ms 0), prefix-agnostic
  decode (`R`/`S`/`SPEC`), malformed → `ErrInvalidID`.
- `FuzzRoundTrip`: fuzzed `ms ∈ [0, 36⁸)` ⇒ `TimeOf(MintAt(p, Epoch+ms)) ==
  Epoch+ms`. `FuzzTimeOf`: arbitrary strings never panic — only `ErrInvalidID`
  or a valid time.

**2. `cli` table-driven** (`internal/cli/cli_test.go`) over
`Run(args, stdin, stdout, stderr, fakeClock)`, asserting `(stdout, stderr,
exit)`:
- Dispatch/help/version/exit-code matrix (D4); mint count & prefix; decode
  args-vs-stdin, partial-failure→1, empty→0 (D5); `TZ`-independence.
- **Wait-loop** via a **fake `Clock`** whose `Sleep(d)` advances virtual `now`:
  `-n N` → N distinct ids and virtual time advanced ≥ N−1 ms; stalled-clock
  still terminates; `N=1` issues zero `Sleep`s. Fake clock is a `cli` test
  helper.

**3. `main`** — no unit tests (no logic); covered only by `make build` producing
`bin/idgen` (a build smoke check).

**Cross-cutting.**
- **Determinism**: fake clock throughout `cli`; suite never touches the real
  wall clock — `go test ./...` is fast and reproducible.
- **`-race`**: tool is single-goroutine (detector moot) but `go test -race
  ./...` is cheap CI insurance.

**Rejected.**
- *Snapshot golden vectors generated by the current code* — would pass even if
  the epoch silently regressed; independently-derived vectors are the point.
- *Real-time `-n` timing tests* (wall-clock elapsed ≥ N−1 ms) — flaky/slow; the
  fake-clock virtual-advance assertion proves the property deterministically.
- *`os/exec` black-box harness* — slower, flakier, redundant given the `Run()`
  seam; in-memory buffers suffice.

---

## Status

Seams, public interfaces, naming, struct/type definitions, the data model, and
the testing approach are all decided (Decisions 1–7). The construction order that
realizes this design is `docs/idgen-plan.md`.
