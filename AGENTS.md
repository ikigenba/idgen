# AGENTS.md

`idgen` is a **spec-only** repository. It ships the contract and prompts that
*generate* a small Go CLI — not the source itself. The build is driven by the
[`ralph`](https://github.com/ikigenba/ralph) harness, which reads `project/` and
writes `cmd/`, `internal/`, and `go.mod` one phase at a time. Treat `project/` as
the source of truth; the code is downstream of it.

This repo also carries the spec-process workflow files used to maintain
`project/`. Those workflow files are versioned with the repo so anyone checking
out the project has the process knowledge needed to work on it.

If `cmd/`, `internal/`, or `go.mod` are absent, that is expected — they have not
been generated yet. Do not scaffold them by hand; that is `ralph`'s job.

## Build & test

Requires Go 1.26+ and a `ralph`-supported agent. Always run from the repo root.

- `project/loops/run` — run the whole build loop (**gather → build → verify → …**)
  until `gather` reports `DONE`. It is a shell wrapper that launches `ralph` on the
  three build-loop prompts with the project's chosen harness and model; invoke
  `ralph` directly to override either.
- `ralph project/loops/audit.md` — the optional, separate audit loop:
  adversarially re-checks that every requirement id is proven by a test that
  could actually fail. Read-only against the checkout; findings land in the
  gitignored `project/audit/REPORT.md`.
- `make build` — compile to `bin/idgen` (only after code has been generated).
- `make test` — `go test ./...`.
- `go test -race ./...` — the verify gate; must be green before a phase is marked done.

## Project layout

- `project/product/README.md` — *why* idgen exists (outcomes only, never mechanism).
- `project/research/` — optional, non-contractual: collected external ground truth
  that the design references.
- `project/design/` — *how*: one `DNN.md` per Decision, a `README.md` spine, and
  `INDEX.md`. Each Decision ends with a Verification list of minted `R-XXXX-XXXX`
  ids. **That id set is the requirement denominator** — there is no separate
  requirements doc.
- `project/plan/` — construction order: one `phase-NN.md` per package of work.
- `project/plan/STATUS.md` — the manifest and the **only** home of the `⬜`/`✅`
  phase markers.
- `project/loops/` — the generated loop prompts (the three build-loop prompts
  plus the single audit prompt), the `run` wrapper, and `README.md` describing
  the installed loops.
- Workflow files for agent-specific tooling may be versioned with the repo, but
  agent-specific discovery and configuration belong to that agent's own
  conventions.

## The workflow

Every spec change is the same three beats: discuss the goal in conversation,
sharpen it with `$grillme`, then `$codify` writes product/research/design/plan
in one pass (greenfield included).
`$create-gather-build-verify-prompts` and `$create-audit-prompts` generate
`project/loops/` once per project or when a loop's design changes. Launching
`ralph` is always an explicit operator action — normally `project/loops/run`
for the build loop.

## Conventions

- **The authorities never restate each other.** Product owns *why*, design owns
  *how* + the id denominator, plan owns *order*. Keep a fact in exactly one place.
  The artifact shapes live in the project workflow files; the loop mechanics
  live in `project/loops/README.md`.
- **Tests are id-tagged.** Every checkable behavior carries its `R-XXXX-XXXX` id in
  the test as a `// R-XXXX-XXXX` comment, so coverage is a `grep`.
- **`STATUS.md` markers are load-bearing.** gather finds the next phase with
  `grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1`. Never add a bare
  `⬜`/`✅` glyph outside a phase line.
- **IDs anchor to a 2026 UTC epoch** and are minted only from already-elapsed instants.

## Boundaries — do not edit

- `project/loops/brief.md` — gitignored, authored by gather and deleted by verify.
  It is `ralph`-owned transient state; never hand-edit or commit it.
- `project/audit/` — gitignored transient state of the audit loop (its manifest
  and report); a fresh audit starts fresh. Never hand-edit or commit it.
- `bin/` — build output, gitignored.
- Generated `cmd/`/`internal/`/`go.mod` — change the design in `project/`, then let
  the loop regenerate, rather than editing generated code directly.

## Commits & PRs

- Commit messages are prefixed `idgen:` (e.g. `idgen: prompt-brief fix`).
- A phase is "done" only when `go test -race ./...` is green and every id in the
  phase's brief is covered by a genuine id-tagged test; verify flips its `STATUS.md`
  marker `⬜ → ✅` at that point — no earlier.
