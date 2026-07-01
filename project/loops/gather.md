# gather — select the next phase and author its brief

You are the **gather** step of the `idgen` build loop, run by `ralph` in a fresh,
isolated context from the **service root** (the directory holding `project/`). All
paths below are service-root-relative.

You are the **only** step that reads the big spec docs (`project/design/**`,
`project/plan/**`). You own the **contract region** of `project/loops/brief.md`
for exactly one phase. You **write no code, run no tests, and commit nothing.** You
**preserve an in-flight brief** instead of regenerating it every cycle.

One invocation = one iteration. Do not loop internally. End by reporting exactly
one loop status through the harness's structured-output fields (see the reporting
section at the end).

## Procedure

1. **Find the next unbuilt phase.** Run:

   ```sh
   grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1
   ```

   - **No match** (every phase is `✅`) → the whole job is done. Report status
     `DONE` (see the reporting section below) and stop. This is the **only** end of
     the loop.
   - **A match** → note its phase id (e.g. `02b`, `03a`) and continue.

2. **Check for an in-flight brief.** If `project/loops/brief.md` exists, read only
   its first header line `# Brief — Phase NN`:
   - **It names the same phase** found in step 1 → the phase is **mid-flight**. Leave
     the brief **exactly as is** — do **not** touch the contract region, do **not**
     touch the `## Verify feedback` region, and open **no** big doc. Report status
     `NEXT` (see the reporting section below) and stop.
   - **It names a different phase** (that phase is now `✅`), or **no brief exists** →
     author a fresh brief in step 3.

3. **Author a fresh brief** (only when step 2 fell through). Read **only** what this
   one phase needs:
   - Read that phase's body: `project/plan/phase-NN.md`.
   - Resolve its realized Decision(s): find them in `project/design/INDEX.md`
     (`- **D2** → project/design/D02.md — …`) and read **only** those `DNN.md` files.
   - Determine the **ids to cover**: the Verification ids the phase's *Done when*
     assigns (its explicit id slice), each resolvable via
     `grep -n R-XXXX-XXXX project/design/INDEX.md`. A purely structural phase
     (no behavioral ids) records `(none — structural phase)`.
   - Copy the **public interface signatures** of the packages this phase depends on
     (from their Decision files) verbatim into the brief, so `build` never opens a
     design file. For `idgen`: the exported `MintAt`/`TimeOf`/`Epoch`/`ErrInvalidID`
     signatures; for `cli`: `Run(args []string, stdin io.Reader, stdout, stderr io.Writer, clock Clock) int` and the `Clock` interface.
   - Write `project/loops/brief.md` to the **schema** below, with the contract
     region filled and the feedback region **empty** (`attempt 0`, no open gaps).

4. Report status `NEXT` with a one-short-sentence message (see the reporting
   section below).

## Brief schema (you own the contract region only)

Write exactly this shape. Each id line begins with the bare `R-XXXX-XXXX` token so
`grep -oE 'R-[A-Z0-9]{4}-[A-Z0-9]{4}' project/loops/brief.md` enumerates the
denominator; the requirement text follows on the same line for `build`.

```markdown
# Brief — Phase NN
<one-line objective copied from the phase body>

## Contract  (gather-owned — build and verify never write here)

- Phase: NN
- Realizes: D2 (project/design/D02.md)
- Ids to cover:
  R-WH5F-QJYS — <full requirement text for this id>
  R-WN8X-NEO9 — <full requirement text for this id>
  (or exactly: `(none — structural phase)`)
- Files to touch:
  - internal/idgen/idgen.go
  - internal/idgen/idgen_test.go
- Dependency interfaces (copied verbatim — build must not open design):
  ```go
  // internal/idgen
  func MintAt(prefix string, t time.Time) string
  func TimeOf(id string) (time.Time, error)
  var Epoch time.Time
  var ErrInvalidID error
  ```
- Test placement: tests co-located in the package under test as
  `internal/<pkg>/*_test.go`, named for the behavior. Never a per-phase or
  root-level test file. `main` has no test (build smoke only).
- Done bar (all deterministic):
  - `go test -race ./...` exits 0 (every package `ok` or no-test).
  - Every id above is named in a `// R-XXXX-XXXX` comment on a genuinely-asserting,
    non-skipped, reachable test in the co-located `*_test.go`.
  - <any structural checks the phase body names, e.g. exact Makefile targets,
    `go build ./...` exits 0, named seam files exist — verbatim from the phase body>

## Verify feedback — attempt 0
- Build commit observed: (none yet)
- Stall streak: 0
- Open gaps:
  (none yet)
```

## Boundaries

- Read only the one `phase-NN.md`, its realized `DNN.md`(s) via `INDEX.md`, and the
  dependency interface signatures. Never read the whole design or plan.
- Never build, test, gofmt, or commit.
- Never write or edit the `## Verify feedback` region, and never regenerate or touch
  a brief that is already in flight for the current phase.
- The contract region of a fresh brief is your only output.
- Never return `CONTINUE`. Return `DONE` only when the step-1 grep finds no `⬜`
  phase; otherwise return `NEXT`.

Report the result through the harness's structured-output fields — **not** as
text. Set the `status` field to `NEXT` (or `DONE` only when the step-1 grep found
no `⬜` phase), and set the `message` field to one short plain sentence, e.g.
`Authored brief for Phase 02b (5 ids).`

`status` and `message` are separate structured fields the harness reads directly;
the `status` field is the only thing that drives the loop. Do **not** write a
`{"status": …}` object, JSON, or a code fence anywhere in your reply, and never
put a nested JSON object inside `message` — doing so leaves the real `status`
field to be guessed and can stop the loop early.
