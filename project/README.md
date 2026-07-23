# idgen/project — workspace layout

Everything needed to **design, plan, and build** `idgen` lives under `project/`.
This file is the workspace map — a map, not a manual. The application code
(`cmd/`, `internal/`, `go.mod`) does not exist yet — the unattended build loop
writes it, one phase at a time, from the spec below. Paths throughout the spec
are written relative to the **repository root**, which is also the directory the
build loop runs from.

This repo is spec-driven: the source of truth is this `project/` tree. To change
behavior you change the spec (open a session with `$open-spec`, settle the goal in
conversation, then `$seal-spec` writes it) and let the loop implement it — never
hand-edit generated code to add behavior.

## The folders

| folder | what's in it | written by |
|---|---|---|
| `product/` | `README.md` — the *why*: problem, users, scope, user-facing promises, success criteria | `$seal-spec` (rewritten in place) |
| `research/` | `research.md` — collected external ground truth the design references; optional | `$seal-spec` (rewritten in place) |
| `design/` | `README.md` (spine) + `INDEX.md` (manifest + sorted `R-id → Decision` map) + `DNN.md` (one per Decision) | `$seal-spec` (rewritten in place) |
| `plan/` | `README.md` (rules) + `STATUS.md` (the manifest — `Next phase` counter + the only home of the `⬜` markers) + `phase-NN.md` (one per **pending** phase) | `$seal-spec` (appends); the build loop deletes completed phases |
| `loops/` | the generated build-loop prompts + `README.md` describing the installed loop | `$create-gather-build-verify-prompts` |

The artifact shapes, the authority boundaries between them (product owns *why*,
design owns *how* + the `R-XXXX-XXXX` requirement-id denominator, plan owns
*order*), and the hard invariants are defined once, in the `ikispec` skill
(`.claude/skills/ikispec/SKILL.md`) — not restated here. The design's
Verification ids are the requirement denominator; there is no separate
requirements doc.

## The build loop

How the installed loop works — its prompts, the status contract, the state
machine, and the `brief.md` schema — is documented beside the prompts, in
[`project/loops/README.md`](loops/README.md), and is run from the repository root
via the committed `project/loops/run` wrapper. This map restates none of those
mechanics.
