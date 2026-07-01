---
name: ralph
description: Map of the ralph autonomous-build family — how the executor, the doc-authoring pipeline (product/research/design/plan), /create-gather-build-verify-prompts, the gather→build→verify loop, and the supporting skills fit together. Load this first when doing any ralph/loop work.
---

# Ralph

`ralph` builds a software project **unattended**: you author a spec interactively
with an agent, then hand a generated prompt loop to the `ralph` binary, which
drives an agent in a fresh context every turn until the work declares itself done.

This skill is the **map** — what each part is and how they connect. The detailed
contracts live in the individual commands/skills named below; load those when you
actually do the work.

## The two halves

The family splits cleanly in two:

- **Authoring (interactive, human-in-the-loop)** — you and an agent build the
  `project/` spec together, one decision at a time. These are **slash commands**
  (`/*-mode`, `/create-gather-build-verify-prompts`): imperative session modes you invoke.
- **Execution (unattended, no human)** — the `ralph` binary re-runs the generated
  prompt files in a loop until done. No interaction; all state lives in the
  workspace.

Intelligence lives at the top (composing the spec); everything below is
mechanical execution.

## The executor — `ralph` (`~/projects/ralph`)

A Go binary that runs a sequence of prompt files, each as a clean-context
`codex exec`, cycling until a prompt reports `DONE` or a budget rail trips. It is
**work-agnostic**: it owns only the lifecycle and the budget rails
(`--max-iterations/-time/-spend/-tokens`); the prompts own the work.

`ralph`'s only assumptions are that it runs from the **root directory of the code
it's changing** — in this mono-repo that's the **service root** (e.g. `wiki/`),
not the repo root — and that it's handed the **full paths to the prompt files**.
It knows nothing about what those prompts are named or where they live: that is a
per-project convention, documented in the project's `project/README.md`. (Our
convention is `project/prompts/{gather,build,verify}.md`, with all workspace paths
written relative to the service root.)

Each prompt ends its final message with one JSON object whose `status` drives the
loop:

- `NEXT` — advance to the next prompt (wrapping past the last back to the first).
- `DONE` — the whole job is complete; stop, exit 0.
- (`CONTINUE` re-runs the same prompt; the build loop below doesn't use it.)

State never lives in `ralph` or the model's memory between turns — only in files.
See `~/projects/ralph/README.md` for flags, the status schema, budget rails, exit
codes, and the run ledger at `~/.ralph/runs.jsonl`.

## The authoring pipeline (interactive `/*-mode` commands)

Run in order; each produces part of `project/`. Every mode is **re-usable** (re-enter
to evolve the spec) and **edits its doc in place** (never append) — except the
plan, which is append-only history.

| Command | Produces | Owns |
| --- | --- | --- |
| `/product-mode` | `project/product/product.md` | **why / intent** — problem, users, scope, promises, success criteria (outcomes only; no mechanism) |
| `/research-mode` | `project/research/research.md` | optional, **non-contractual** background that informs *you* before design; nothing downstream reads it |
| `/design-mode` | `project/design/README.md` (spine) + `project/design/INDEX.md` (manifest) + `project/design/DNN.md` (one per Decision) | **shape + its proof** — seams, interfaces, types, the test strategy, and the minted `R-XXXX-XXXX` Verification ids |
| `/plan-mode` | `project/plan/README.md` (rules) + `project/plan/STATUS.md` (manifest) + `project/plan/phase-NN.md` (one per phase) | **construction order + history** — dependency-ordered phases; append-only; `⬜`/`✅` markers live only in STATUS.md |

The design and plan are **split for addressability**: the loop greps a manifest
for the next unit of work and reads only the one file it needs, never the whole
architecture. That split is the precondition for the loop below.

The authority boundary is load-bearing: product states the *promise*, design
states the *exact checkable proof* of that promise, plan states the *order it's
built in*. They never restate each other.

## Generating the loop — `/create-gather-build-verify-prompts`

Once `project/` exists, `/create-gather-build-verify-prompts` emits the **three-prompt build loop** the
executor runs:

- `project/prompts/gather.md` — the **only** prompt that reads the big docs. Greps STATUS.md
  for the first `⬜` phase (none → `DONE`, the sole exit); if a brief for that same
  phase already exists it leaves it untouched (no-op — the phase is mid-flight),
  otherwise it resolves the phase's Decisions via INDEX.md and writes a tiny,
  self-contained `project/prompts/brief.md` (contract region + empty feedback region). Returns `NEXT`.
- `project/prompts/build.md` — reads **only** the brief (contract *and* `verify`'s feedback
  region — open gaps first). Builds the package, writes id-tagged tests, runs the
  suite, commits. Leaves the marker untouched. `NEXT`.
- `project/prompts/verify.md` — the independent gate; the **only** prompt that flips a
  marker. Pass → flip `⬜→✅` + commit + delete the brief; gap → leave `⬜` and
  **write its open gaps into the brief's feedback region** (the brief persists so
  the next build sees them), unless the same gaps have stalled for 3 cycles, in
  which case it discards the brief to reset the trajectory. `NEXT`.

The author-facing overview of the loop — the invocation, the status contract, the
state machine, and the `brief.md` schema — lives in the project's `project/README.md`,
not in any generated file.

`project/prompts/brief.md` is the **seam** that keeps build's context tiny:
gitignored, single-phase, and **phase-scoped** — gather authors its contract
region once per phase and then no-ops while that phase is in flight; verify
deletes it on a pass (or a stall reset) and otherwise persists it, overwriting a
single feedback region with its open gaps so the next build is informed. The two
writers own disjoint regions (gather → contract, verify → feedback), so neither
clobbers the other. Because gather is the only big-doc reader and verify is the
only marker-flipper, the loop is human-free and converges — an incomplete phase
just stays `⬜` and gets re-attacked next cycle with verify's feedback in hand;
the only stops are gather's `DONE` or a budget rail.

Run it (from the service root): `ralph project/prompts/gather.md project/prompts/build.md project/prompts/verify.md`

## Supporting reference skills

- `repo` — create a new repo with the three-tier git layout (GitHub push-mirror,
  bare source-of-truth at `/mnt/store/git/<org>/<repo>`, local clone at
  `~/projects/<repo>`).

## The flow at a glance

```
/product-mode ─► /research-mode ─► /design-mode ─► /plan-mode ─► /create-gather-build-verify-prompts
   product/         research/         design/*         plan/*        prompts/{gather,build,verify}.md
                                                                         │
                       (from the service root)  ralph project/prompts/gather.md project/prompts/build.md project/prompts/verify.md
                                                  (unattended: gather ─► build ─► verify ─► …)
```
