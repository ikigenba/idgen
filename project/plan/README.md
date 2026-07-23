# idgen — Plan

**Authority: construction order.** This document and the `project/plan/`
directory it heads own the **order idgen is built in** — a work queue of
**pending** phases only, nothing else (the *why* is product's, the *shape and its
proof* is design's). The plan holds no finished work: when a phase completes the
build loop **deletes** its `STATUS.md` line and its `project/plan/phase-NN.md` in
the completion commit, so the queue can never contradict a design that has since
moved on. **Completion is deletion — there is no `✅` state on disk; the record of
what was built lives in git** (the completion commits, and the deleted files
recoverable there), never here. To extend idgen: update
`project/product/README.md` and `project/design/` **in place** first, then
**append** a new phase — a new `project/plan/phase-NN.md` body plus a new
`project/plan/STATUS.md` line, numbered from the `Next phase` counter. Phase
numbers are **never renumbered and never reused**.

**Coverage invariant.** Every *current* design Verification id (`R-XXXX-XXXX`) is
either already **realized** — its id appearing as a tag in a test that runs under
the suite — or assigned to **exactly one** pending phase: none unassigned, none
split, none duplicated across pending phases.

**One phase = one package = one build-turn context.** Each phase is a single
coherent unit of work — almost always one Go package — scoped to that unit's
design Decisions and the *interfaces* (not internals) of the packages it depends
on, and **sized so the build loop can carry it in one fresh build-turn context**
and ideally finish it in a turn or two. The loop does *not* build a phase in one
long accumulating context — size to a single build turn, not an imagined single
sitting; sizing a phase as large as cleanly fits one turn is good (fewer cycles,
less context churn). Where a single Decision is too large for one context it is
split across phases, and each affected phase names the **slice** of that
Decision's Verification ids it carries.

**Done bar.** A phase is **done** when every Verification id (`R-XXXX-XXXX`) of the
design Decision(s) it realizes — or the explicit slice of those ids assigned to it
— is covered by a clearly-named test, and the suite is green. "Green" is defined by
design's *Conventions*: **`go test -race ./...` exits 0** (every package `ok` or
no-test). "Covered" is defined by each Decision's **Verification** list — the
concrete behavior the id names, asserted by a test that tags the id in a
`// R-XXXX-XXXX` comment, so coverage is a grep. Every phase's acceptance bar is
expressed as **deterministic exit conditions**: mechanically-checkable predicates
(a green suite, an exit code, an exact match count) reproducible on identical repo
state and whose passing state is actually reachable — never a subjective prose
judgment, and never a self-referential or unsatisfiable check. A purely
**structural** phase (no behavioral ids yet reachable) earns its own deterministic
check (a clean build, the exact set of named files/targets) rather than a prose
claim.

## Layout

The plan is physically split so the build loop reads only the one unit of work it
needs, never the whole history:

- **`project/plan/STATUS.md`** — the manifest: the `Next phase: NN` counter plus
  one line per **pending** phase in build order, the **only** home of the `⬜`
  markers. The loop finds its next work with
  `grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1`, reads only that
  phase's body file, and on completion deletes that phase's line.
- **`project/plan/phase-NN.md`** — one body file per phase (zero-padded; a
  sub-phase keeps its suffix, e.g. `phase-03a.md`). Carries **no** status marker
  of its own.
- **`project/plan/README.md`** — this file: the static rules above. It lists no
  phases and carries no status, so it never grows with the project.

**Completion is deletion, for this layout:** the build loop's only mutations are
removing a finished phase's `STATUS.md` line together with its `phase-NN.md` in the
completion commit. The `Next phase` counter is never decremented and a number is
never reused, so a phase number names one phase forever even after its files are
gone.
