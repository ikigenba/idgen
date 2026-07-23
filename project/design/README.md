# idgen — Design

**Authority: shape and its proof.** This directory owns *how idgen is built* —
seams, public interfaces, naming, types, the data model, the encoding math — and
*how each behavior is proven*. `project/product/README.md` owns the *why* and the
user-facing promises; design states the **exact, checkable form** of those
promises and never re-declares the why. The contractual constants (the 2026 epoch
and the version string) are the product's; design uses their values, it does not
own them.

This is the **single, current** statement of the architecture: when a decision
changes, its `DNN.md` is rewritten in place to stay true (stale decisions are
removed, not stacked). History of how it got here lives in `project/plan/`.

## Requirement ids

- Each Decision ends with a **Verification** list: the concrete behaviors a test
  must assert for that decision to be considered built.
- Every Verification item carries a minted **idgen id** (`R-XXXX-XXXX`) — a
  stable, unique handle for that one behavior. The ids live inline in those lists
  and **nowhere else**; there is **no separate requirements document**. Ids are
  minted with `idgen` itself (`idgen -n <count> -p R`), never hand-written: a
  fresh id per newly added behavior, an existing id never renumbered, and a
  removed behavior takes its id (and its test) with it.
- **Design's responsibility for these ids ends at minting them.** This id set is
  the requirement denominator, but *how* coverage is measured against it and
  *when* the work is "done" are downstream concerns (the plan and the build
  loop) — they are not specified here.

## Conventions

Shared facts every Decision leans on:

- **Language / toolchain.** Go 1.26, module path `github.com/ai4mgreenly/idgen`.
  Standard library only — no third-party runtime dependency.
- **Build / typecheck command.** `go build ./...`. The `Makefile` (D6) wraps the
  toolchain: `make build` produces `bin/idgen`; `make test`, `make fmt`,
  `make clean`, `make install` are the other targets.
- **Test command — the green gate.** `go test -race ./...`. **"The suite is green"
  means it exits 0** with every package `ok` or no-test. The race detector is cheap
  CI insurance even though the tool is single-goroutine.
- **Exit codes.** `0` success · `1` decode data failure (≥1 malformed id in an
  otherwise valid invocation) · `2` usage/runtime error (bad flags, empty/invalid
  prefix, `number ≤ 0`, mint-with-positionals).
- **No silent failures.** Every path that exits `1` or `2` writes a non-empty,
  descriptive message to stderr before returning — an empty-output non-zero exit
  is always a defect, never an acceptable outcome. This applies to `flag`'s own
  parse errors (already routed to stderr via `fs.SetOutput(stderr)`) and to every
  hand-written validation error in `cli` alike; each Decision's Verification list
  for a `1`/`2` exit path must assert stderr is non-empty (or contains the
  expected fragment), not merely check the exit code.
- **Time format.** Every printed time is UTC, formatted `2006-01-02T15:04:05.000Z`.
- **Time source.** Standard-library `time` only — millisecond precision is portable
  to every target Go compiles for (`time.Now()` resolves far finer than a
  millisecond on every supported platform). `Epoch` is a constructed
  `time.Date(..., time.UTC)` value carrying **no** monotonic reading, so
  `time.Now().Sub(Epoch)` strips monotonic and yields a pure **wall-clock** elapsed
  duration — an absolute civil instant decodable forever, not host uptime (monotonic
  time has no fixed zero across reboots and could never anchor a decodable id). A
  backward wall-clock step (NTP, manual change) is tolerated, not detected — see D3.
- **Testability seams.** The `idgen` core is a pure function of its inputs (no
  clock, no I/O), and the whole CLI sits behind one injectable
  `cli.Run(args, stdin, stdout, stderr, clock) int`. Every behavior is reachable
  **in-process and deterministic** — no subprocess, no real sleeps: `idgen` proves
  its requirements with unit + randomized property tests at zero process setup;
  `cli` proves its
  requirements through in-memory `args`/`stdin`/`stdout` buffers, a return code, and
  a **fake `Clock`** whose `Sleep` advances virtual time; `main` has no logic and
  carries no requirement (a build smoke check stands in for it).

## Layout

The design is **split for addressability** so a build phase reads only the one
Decision it realizes:

- `project/design/INDEX.md` — the manifest: each Decision → its file, plus a sorted
  `R-id → Decision/file` reverse map. Regenerated whenever a Decision is added or
  its Verification ids change.
- `project/design/DNN.md` — one self-contained file per Decision (zero-padded
  `D01.md`, `D02.md`, …; referenced in prose and the plan as `D<N>`). Holds
  **Decision**, **Rejected**, and **Verification** (the minted ids).
- `project/design/README.md` — this spine: cross-cutting facts only, no
  per-Decision detail.

Design is **rewritten in place** (history lives in git, not in the spec): a
changed Decision is rewritten in its `DNN.md` and `INDEX.md` is regenerated; a new
Decision adds a `DNN.md` and an INDEX entry.
