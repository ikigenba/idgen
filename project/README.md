# idgen/project — workspace layout

Everything needed to **design, plan, and build** `idgen` lives under `project/`.
This file is the workspace map; the application code (`cmd/`, `internal/`,
`go.mod`) does not exist yet — the unattended build loop writes it, one phase at a
time, from the spec below. Paths throughout the spec are written relative to the
**repository root**, which is also the directory the build loop runs from.

This repo is spec-driven: the source of truth is this `project/` tree. To change
behavior you change the spec (via a `/*-mode` command) and append a plan phase; the
loop implements it — never hand-edit generated code to add behavior.

## The folders

| folder | what's in it | owned by |
|---|---|---|
| `product/` | `README.md` — the *why*: problem, users, scope, user-facing promises, success criteria | `/product-mode` (rewritten in place) |
| `design/` | `README.md` (spine) + `INDEX.md` (manifest + sorted `R-id → Decision` map) + `DNN.md` (one per Decision) | `/design-mode` (rewritten in place) |
| `plan/` | `README.md` (rules) + `STATUS.md` (the manifest — the only home of each phase's `⬜`/`✅` marker) + `phase-NN.md` (one per phase) | `/plan-mode` (append-only) |
| `loops/` | the build-loop prompts `gather.md`, `build.md`, `verify.md` (+ the phase-scoped `brief.md`) | `/create-gather-build-verify-prompts` (generated) |

The three **spine documents** (`product/README.md`, `design/README.md`,
`plan/README.md`) are each singular and owned by a `/*-mode` command — that command
is the sanctioned way to change them. Product and design are **rewritten in place**
to stay authoritative; the plan is **append-only** (it doubles as construction
history). The `R-XXXX-XXXX` ids that the design mints are the **denominator** — the
enumerated intent the test suite is measured against; there is no separate
requirements doc.

## The build loop

Run it from the repository root:

```
ralph project/loops/gather.md project/loops/build.md project/loops/verify.md
```

It cycles the prompts in fresh, isolated contexts — `gather → build → verify → …` —
on a **two-status** contract: each prompt ends with one JSON object whose `status`
is either `NEXT` (advance to the next prompt, wrapping `verify → gather`) or `DONE`
(stop). State lives entirely in the workspace (the git tree, `project/plan/STATUS.md`,
and the phase-scoped `project/loops/brief.md`) — never in the agent's memory
between turns.

| step | reads | writes | commits | flips marker | returns |
|---|---|---|---|---|---|
| **gather** | the big docs (STATUS, one phase, its Decisions) — or just the brief header, if in-flight | `project/loops/brief.md` contract (fresh phase only) | no | no | `NEXT`, or `DONE` if no `⬜` |
| **build** | `project/loops/brief.md` only (contract + verify feedback) | code + co-located tests | the increment | no | `NEXT` |
| **verify** | the brief + the suite | brief's feedback region (gap) / deletes the brief (pass or stall) | only a marker flip (on pass) | yes (pass only) | `NEXT` |

- **gather** — the only step that opens `project/plan/` or `project/design/`. It
  greps `STATUS.md` for the first `⬜` phase; none → `DONE` (the sole end of the
  loop). If a brief for that same phase already exists it is left untouched (the
  phase is mid-flight — its contract and verify's feedback are preserved) and gather
  returns `NEXT` without opening a big doc. Only when no brief exists (or it belongs
  to an already-`✅` phase) does gather resolve the phase's Decision(s) via
  `project/design/INDEX.md` and write a tiny, self-contained `brief.md` with an empty
  feedback region.
- **build** — never opens the big docs. It consumes only the brief — including the
  dependency interface signatures copied into it **and verify's feedback region**,
  whose open gaps it addresses first — does a bounded, idempotent turn, writes
  id-tagged tests (`// R-XXXX-XXXX`), commits, and leaves the marker `⬜`. It reads
  the feedback but never writes it.
- **verify** — the independent gate and only step that flips a marker. It re-runs
  `go test -race ./...` and checks that every id is covered by a genuinely-asserting
  test, re-deriving the truth from scratch. Pass → flip that one `⬜ → ✅`, commit
  the flip, **and delete the brief**. Gap → leave `⬜`, change no source, and **write
  the open gaps into the brief's feedback region** (the brief persists so the next
  build sees them). If the same gaps stall for three no-progress cycles, verify
  discards the brief to reset the trajectory. Either way it returns `NEXT`.

### Why it is human-free and converges

`verify` can neither halt the loop nor advance a phase on a gap — its only powers
are "flip this phase green" (on full proof) or "leave it red." So an incomplete or
wrongly-built phase simply stays `⬜`, and the next cycle re-attacks it with a fresh
context — now with verify's grounded feedback in front of build, and without gather
re-reading the big docs (it no-ops on the in-flight brief). The loop ends only when
**every** phase is verified green (`gather` finds no `⬜` and returns `DONE`) — or
when a `ralph` budget rail trips. The marker is the sole completion signal and only
verify, only on proof, ever moves it.

### The brief is the seam

`project/loops/brief.md` keeps build's context tiny — the complete and only input
build and verify consume, so neither opens design or plan. It is **gitignored**,
**single-phase**, **phase-scoped** (authored once per phase, persists across cycles
while the phase stays `⬜`), and **region-owned by a single writer each**: a
gather-owned contract region (phase identity, ids, files, interface signatures, done
bar) and a verify-owned `## Verify feedback` region (the currently-open gaps). gather
writes the contract and never the feedback; verify writes the feedback and never the
contract; build reads both and writes neither.

#### The schema

The brief has a strict, grep-able shape so all three prompts read it mechanically —
the "Ids to cover" list is one bare `R-XXXX-XXXX` per line, so the denominator is
`grep -oE 'R-[A-Z0-9]{4}-[A-Z0-9]{4}' project/loops/brief.md | sort -u`:

```
# Brief — Phase NN

## Contract   (gather-owned — verify never writes here)

- **Phase:** NN — <one-line objective>
- **Realizes:** D<x>[, D<y>]        (or: — (structural, no ids))
- **Decision files:** project/design/D0X.md[, project/design/D0Y.md]

### Ids to cover
R-XXXX-XXXX
R-YYYY-YYYY
<one bare id per line — or exactly: (none — structural phase)>

### Files to touch
- <source file(s) + the co-located *_test.go that holds the id-tagged tests>

### Dependency interface signatures
```go
// from internal/idgen (D2)
func MintAt(prefix string, t time.Time) string
func TimeOf(id string) (time.Time, error)
```
<or exactly: (none — no earlier-package dependency)>

### Done bar
- `go test -race ./...` exits 0 (every package `ok` or no-test).
- Every id above is named in a `// R-XXXX-XXXX` comment on a genuinely-asserting,
  actually-run test in its co-located *_test.go.
- <any structural checks verbatim from phase-NN.md's Done when>

## Verify feedback — attempt N   (verify-owned — gather writes this empty)

- Build commit observed: <sha | (none yet)>
- Stall streak: <0..2>
- Open gaps: <(none yet — fresh brief), or a checklist, one line per open gap:
  R-id + the exact failing command + observed output (+ file:line when known)>
```

**Two writers, disjoint regions, no clobber.** gather authors the `## Contract`
region once (when the phase first becomes the active `⬜`) with an empty
`attempt 0` feedback stub, and thereafter no-ops while the phase is in flight.
verify **overwrites** (never appends) the `## Verify feedback` region each gap
cycle — bumping `attempt N`, recording the build commit it observed and the stall
streak — and deletes the whole brief on a pass or on a 3-cycle stall reset. build
reads both regions (open gaps first) and writes neither. The attempt counter and
build commit give verify cross-cycle memory: a shrinking/changing open-gap set is
*slow convergence*; the **same** gap ids with **no new build commit** for three
consecutive attempts is a *true stall*, which triggers the trajectory reset.

Test placement is fixed (design D7) — co-located `*_test.go`, named for the
behavior, never a per-phase or root-level test file: core →
`internal/idgen/idgen_test.go` + `fuzz_test.go`, CLI (incl. the mint↔decode
round-trip through `Run`) → `internal/cli/cli_test.go`, build smoke →
`cmd/idgen/main_test.go`.
