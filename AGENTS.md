# AGENTS.md

`idgen` is a **spec-only** repository. It ships the contract and prompts that
*generate* a small Go CLI — not the source itself. The build is driven by the
[`ralph`](https://github.com/ikigenba/ralph) harness, which reads `project/` and
writes `cmd/`, `internal/`, and `go.mod` one phase at a time. Treat `project/` as
the source of truth; the code is downstream of it.

If `cmd/`, `internal/`, or `go.mod` are absent, that is expected — they have not
been generated yet. Do not scaffold them by hand; that is `ralph`'s job.

## Build & test

Requires Go 1.26+ and a `ralph`-supported agent. Always run from the repo root.

- `ralph project/prompts/gather.md project/prompts/build.md project/prompts/verify.md`
  — run the whole build loop (**gather → build → verify → …**) until `gather`
  reports `DONE`.
- `make build` — compile to `bin/idgen` (only after code has been generated).
- `make test` — `go test ./...`.
- `go test -race ./...` — the verify gate; must be green before a phase is marked done.

## Project layout

- `project/product/product.md` — *why* idgen exists (outcomes only, never mechanism).
- `project/design/` — *how*: one `DNN.md` per Decision, a `README.md` spine, and
  `INDEX.md`. Each Decision ends with a Verification list of minted `R-XXXX-XXXX`
  ids. **That id set is the requirement denominator** — there is no separate
  requirements doc.
- `project/plan/` — construction order: one `phase-NN.md` per package of work.
- `project/plan/STATUS.md` — the manifest and the **only** home of the `⬜`/`✅`
  phase markers.
- `project/prompts/{gather,build,verify}.md` — the three build-loop prompts.
- `.claude/commands/` — the `/*-mode` authoring commands, versioned on purpose.

## Conventions

- **The three authorities never restate each other.** Product owns *why*, design
  owns *how* + the id denominator, plan owns *order*. Keep a fact in exactly one place.
- **Tests are id-tagged.** Every checkable behavior carries its `R-XXXX-XXXX` id in
  the test as a `// R-XXXX-XXXX` comment, so coverage is a `grep`.
- **`STATUS.md` markers are load-bearing.** gather finds the next phase with
  `grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1`. Never add a bare
  `⬜`/`✅` glyph outside a phase line.
- **IDs anchor to a 2026 UTC epoch** and are minted only from already-elapsed instants.

## Boundaries — do not edit

- `project/prompts/brief.md` — gitignored, authored by gather and deleted by verify.
  It is `ralph`-owned transient state; never hand-edit or commit it.
- `bin/` — build output, gitignored.
- Generated `cmd/`/`internal/`/`go.mod` — change the design in `project/`, then let
  the loop regenerate, rather than editing generated code directly.

## Commits & PRs

- Commit messages are prefixed `idgen:` (e.g. `idgen: prompt-brief fix`).
- A phase is "done" only when `go test -race ./...` is green and every id in the
  phase's brief is covered by a genuine id-tagged test; verify flips its `STATUS.md`
  marker `⬜ → ✅` at that point — no earlier.
