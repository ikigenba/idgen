# idgen — the installed build loop

This directory holds the three prompts an unattended harness (`ralph`) re-invokes
to build idgen one phase at a time, plus this overview. It lives beside the prompts
it describes so it can never drift from the loop that is actually on disk. The
`project/` workspace map (`project/README.md`) only points here; the spec shapes
live in the `spec-shapes` skill and are not restated here.

## Invocation

```
ralph project/loops/gather.md project/loops/build.md project/loops/verify.md
```

By convention the committed wrapper `project/loops/run` issues exactly this
invocation with the operator's chosen harness/model flags baked in — a
typing-saver, nothing more.

`ralph` runs from the **service root** (its working directory), hands the three
prompt paths in order, and re-invokes each in a **fresh, isolated context** every
turn. All loop state lives in the workspace — chiefly the ephemeral
`project/loops/brief.md` seam — not in any prompt's memory.

## The status contract

Each turn ends with a `status` and a one-sentence `message`. `ralph` reads only the
**final** message of a turn and advances on its terminal status:

- **`NEXT`** — terminal: advance to the next prompt, wrapping `verify → gather`.
  **build** and **verify** always report `NEXT`.
- **`DONE`** — terminal: stop the loop. **Only `gather` ever reports `DONE`**, and
  only when no `⬜` phase remains in `STATUS.md`.
- **`CONTINUE`** — **non-terminal**: the status a streaming model tags the progress
  messages it emits *before* its terminal message. It never advances or ends the
  loop. (Under codex, gpt-5.5-style backends coerce every streamed message into the
  schema, so a non-terminal value is required for mid-turn narration; `ralph` reads
  only the last message and acts on the terminal `NEXT`/`DONE`.)

The `{status, message}` schema is supplied by the harness out of band (codex via
`--output-schema`; claude via `--json-schema`, surfaced as a synthetic
`StructuredOutput` tool) — the prompts describe only the contract, never a
transport.

## Per-step reads / writes / commits / deletions

| step | reads | writes | commits | retires phase |
|---|---|---|---|---|
| **gather** | `STATUS.md`, one `phase-NN.md`, `INDEX.md`, realized `D0N.md` | brief **contract** region (fresh phase only) | no | no |
| **build** | `brief.md` only | package source + co-located `*_test.go` | yes (each turn) | no |
| **verify** | `brief.md`, the test suite | brief **feedback** region *or* deletes brief; deletes `STATUS.md` line + `phase-NN.md` on pass | yes (phase deletion / stall log) | **yes, on pass only** |

## The brief lifecycle

`project/loops/brief.md` is the seam that keeps build's context scoped to one
phase. It is **gitignored**, **single-phase**, and **phase-scoped** (not
per-cycle):

- **gather** authors the contract region **once**, when a phase first becomes the
  active `⬜` phase. While that phase stays `⬜`, gather **no-ops** on the in-flight
  brief — it never regenerates it, so the contract and verify's feedback persist.
- **build** consumes the brief (contract + feedback) each cycle and commits its
  increment, never touching the brief.
- **verify** either **passes** the phase (delete its `STATUS.md` line and its
  `phase-NN.md` file, commit, **delete the brief**) or records a **gap**
  (overwrite the feedback region with the open gaps, **keep the brief**). The
  brief thus persists across cycles until a pass or a stall reset.

## Why the loop converges

verify can neither halt the loop nor advance a phase on a gap, so an incomplete
phase just stays `⬜` and is re-attacked next cycle — now with verify's grounded,
command-tied feedback in front of build, and without gather re-reading the big docs
(it no-ops on the in-flight brief). The persisted feedback gives verify cross-cycle
memory: it distinguishes *slow convergence* (the open-gap id set shrinking) from a
*true stall* (the **same** gap ids unsatisfied for **3** consecutive attempts with
**no new build commit**). On a true stall it does a **trajectory reset** — discards
the brief and logs the stall — so the next gather rebuilds the contract fresh; this
stays inside "verify never halts / never advances on a gap." The only exit is
`gather → DONE`, which requires **zero `⬜` markers** — so the run ends only when
every phase is verified green (or a `ralph` budget rail trips). Every phase's done
bar is a deterministic predicate (`go test -race ./...` exit 0, `go build ./...`
exit 0, exact `.PHONY` target/file sets, id-tagged `*_test.go` match counts), so
each gate reliably reaches its passing state.

## The `project/loops/brief.md` schema

Two single-writer regions — gather owns the contract, verify owns the feedback —
so the writers never clobber each other.

**Contract region** (gather-authored, written once per phase):

```
# Brief — Phase NN

## Contract

- **Phase:** NN — <one-line objective>
- **Realizes Decision(s):** D<n>[, D<m>]
- **Decision file(s):** project/design/D0<n>.md[, …]

### Design prose — D<n> (<title>)
<the Decision's statement, shape/signatures, and Rejected alternatives, verbatim
from its DNN.md, with that Decision's Verification list omitted>

### Ids to cover
R-XXXX-XXXX — <full requirement text, verbatim from the Verification list>
…                       (one id per line; or "(none — structural phase)")

### Files to touch
<paths build creates/edits>

### Dependency interface signatures
<public signatures of consumed packages, verbatim>

### Done bar
<deterministic acceptance conditions>

## Verify feedback

_(empty — no gaps recorded yet)_
```

**Feedback region** (verify-authored, overwritten each gap cycle):

```
## Verify feedback — attempt <N>

- build commit observed: <sha>
- stall streak: <count>

Open gaps (each must be closed):
- R-XXXX-XXXX — <exact failing command> → <observed output> [file:line]
```

gather writes the feedback region empty on a fresh brief and never touches it
again; verify **overwrites** it (never appends) with only the currently-open gaps;
build reads it but never writes it.

# idgen — the installed audit loop

A **separate, single-prompt loop** that adversarially audits the test coverage of
the finished (or in-progress) build. It is independent of the build loop above:
different invocation, different state, run on demand. It **never modifies the live
checkout** — it reads the design and the tests and records findings.

## Invocation

```
ralph project/loops/audit.md
```

One prompt. `ralph` runs from the **service root** and re-invokes `audit.md` in a
**fresh, isolated context** every turn; `NEXT` wraps straight back to the same
prompt. All state lives in two **transient, gitignored** files under
`project/audit/` — never committed, so a fresh audit starts from a fresh
denominator.

## The status contract

The single prompt owns all three status values (it is the only prompt in its
loop, so — unlike build/verify — it legitimately owns `DONE`):

- **`NEXT`** — terminal: this turn's work is done; re-invoke the same prompt.
- **`DONE`** — terminal: the audit is complete (**no `⬜` Decision remains**) or
  **refused on a red baseline**. The loop stops. The finish turn echoes the
  report's absolute path.
- **`CONTINUE`** — **non-terminal**: mid-turn progress narration; never advances
  or ends the loop.

## The four turn cases

Each invocation is exactly one of:

- **Init** (`project/audit/STATUS.md` absent) — run the **baseline gate**
  (`go test -race ./...`); a **red baseline refuses** the audit (write the failure
  summary, `DONE`). Green → run the **structural sweep** (four deterministic set
  checks), write the report preamble, write the manifest (one line per id-owning
  Decision), `git worktree prune`, `NEXT`.
- **Staleness guard** — the manifest exists but the Decision/id sets re-derived
  from `project/design/INDEX.md` no longer match what it was built from: wipe
  `project/audit/` and re-init the same turn, noting `restarted: denominator
  changed`. (Implicit contract: the spec does not move while an audit runs.)
- **Audit one Decision** — grep the manifest for the first `⬜`, read only that
  `DNN.md`, judge every id in its Verification list, **append** the `## D<N>`
  report section, flip `⬜ → ✅`, `NEXT`.
- **Finish** — no `⬜` remains: append `## Summary`, `DONE`.

The only exits are the red-baseline refusal and the finish turn; every other turn
is `NEXT`, so an interrupted run resumes at the first `⬜` with all findings
intact.

## The two transient files (`project/audit/`, gitignored)

- **`STATUS.md`** — the manifest: one `- D<N> ⬜`/`✅` line per id-owning Decision,
  written by init in Decision order. Each turn greps the first `⬜`
  (`grep -nE '^- D[0-9]+ .* ⬜' project/audit/STATUS.md | head -1`), exactly like
  the build loop's phase lookup. No bare glyph outside a Decision line.
- **`REPORT.md`** — the deliverable, **append-only within a run**: a preamble
  (baseline + denominator + structural sweep), then one `## D<N>` section per
  audited Decision, then a final `## Summary`. Each turn appends its section
  **before** flipping its marker, so it survives any crash.

## The structural sweep (init, deterministic)

Four grep-and-set-compares with defined pass criteria: **orphan tags** (tagged ids
design never minted), **duplicate assignment** (one id in >1 Decision or >1 test),
**plan coverage drift** (design id set vs `project/plan/phase-*.md`), and **INDEX
staleness** (`DNN.md` id set vs `INDEX.md`). Failures are recorded as findings in
the preamble, not aborts.

## The verdict taxonomy (per id)

One verdict per `R-XXXX-XXXX`: **`covered`** (assertion pins the discriminating
property against a substrate that can falsify it), **`weak`** (a tag exists but
asserts a proxy, runs against a mock where a real substrate is named, a degenerate
impl would pass, is skipped/unreachable, or survived its mutation), **`missing`**
(no tag), **`mismatched`** (a tag exists but the test asserts a different behavior
— design/test drift). `weak` and `mismatched` stay separate because the fix
differs.

## Mutation escalation (the tiebreaker)

Static judgment is the baseline. Escalate **only** when the read suspects `weak`
but the test is plausible and "could this test actually fail?" can't be settled by
reading. Per escalation: `git worktree add "$(mktemp -d)" HEAD` **outside the repo
tree**, apply the **minimal mutation** that violates the id's behavior, run the
tagged test's **package** (not the full suite) — **fails → `covered`, survives →
`weak`** — then `git worktree remove --force` **unconditionally** the same turn.
No mutation ever touches the live checkout.

## The report is the deliverable

`REPORT.md` is verdict-first on each id line, so the gap list is greppable. Harvest
the audit's product with the summary's work-queue grep:

```
grep -E 'R-.* (weak|missing|mismatched)' project/audit/REPORT.md
```
