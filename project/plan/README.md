# idgen — Plan

**Authority: construction order and history.** This document and the
`project/plan/` directory it heads own the **order idgen is built in** and the
**record of what has been built** — nothing else (the *why* is product's, the
*shape and its proof* is design's). Unlike product and design — which are rewritten
in place to stay authoritative for the current state — the plan is **append-only**:
completed phases are never rewritten or deleted, so the plan doubles as the
construction history. To extend idgen: update `project/product/product.md` and
`project/design/` **in place** first, then **append** a new phase — a new
`project/plan/phase-NN.md` body plus a new line in `project/plan/STATUS.md`. Never
edit a finished phase except to flip its status marker in `STATUS.md`.

**One phase = one package = one accumulating context.** Each phase is a single
coherent unit of work — almost always one Go package — built in one accumulating
context against product and design, reading only that unit's design Decisions and
the *interfaces* (not internals) of the packages it depends on. That is what keeps
every phase the size of a small standalone tool no matter how large the project
grows. Where a single Decision is too large for one context it is split across
phases, and each affected phase names the **slice** of that Decision's
Verification ids it carries.

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

- **`project/plan/STATUS.md`** — the manifest: one line per phase in build order,
  and the **only** home of status markers (`✅` done / `⬜` not started). The loop
  finds its next work with
  `grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1` and reads only that
  phase's body file.
- **`project/plan/phase-NN.md`** — one body file per phase (zero-padded; a
  sub-phase keeps its suffix, e.g. `phase-03a.md`). Carries **no** status marker
  of its own.
- **`project/plan/README.md`** — this file: the static rules above. It lists no
  phases and carries no status, so it never grows with the project.

**Append-only, for this layout:** never rewrite or delete a `phase-NN.md`, never
delete a `STATUS.md` line. The only build-time mutation is flipping one phase's
`⬜ → ✅` in `STATUS.md`.
