# idgen

This repository is **`idgen` with no source code.** It ships only the spec and
prompts needed to *generate* it: a product brief, a split design, a build plan, and
the three [`ralph`](https://github.com/ikigenba/ralph) build-loop prompts. Point
`ralph` at them and it writes the code, one phase at a time, until the tool is built
and every requirement is proven.

So it is two things at once:

1. **A way to get and build `idgen`** — clone it, run `ralph`, end up with a
   working binary.
2. **A worked demonstration of the `ralph` method** — how a small, fully-specified
   project is authored interactively with the `/*-mode` commands and then built
   autonomously from *spec + prompts*, with no state held in the model's memory
   between turns.

## Build it with `ralph`

Requires Go 1.26+, a `ralph`-supported agent (e.g. `codex`), and the
[`ralph`](https://github.com/ikigenba/ralph) harness.

Clone the repo and, **from the repository root**, run `ralph` against the three
build-loop prompts:

```sh
git clone https://github.com/ikigenba/idgen.git
cd idgen
ralph project/prompts/gather.md project/prompts/build.md project/prompts/verify.md
```

That single command is the whole build. `ralph` cycles the three prompts in fresh,
isolated contexts — **gather → build → verify → …** — writing one package per phase
until `gather` finds no unbuilt phase and reports `DONE`. You end up with the
`cmd/`, `internal/`, and `go.mod` that weren't in the repo, with every design
requirement id covered by an id-tagged test.

You now have source — build and test it the normal way:

```sh
make build      # compile to bin/idgen
make test       # go test ./...
```

Always invoke `ralph` from the repository root so the prompt paths and the
`project/` spec they read resolve correctly. `ralph` owns the lifecycle and the
budget rails (`--max-spend`, `--max-time`, …); see its README for the flags.

## What idgen is (the end product)

A small CLI that mints stable spec/requirement IDs of the form `R-XXXX-XXXX` from
the wall-clock instant of minting — and decodes them back to that instant. Every ID
is anchored to a **2026 UTC epoch**, so it stays traceable to the moment it was
minted.

```sh
$ idgen
R-4K7P-9ZQ2

$ idgen --decode R-4K7P-9ZQ2
2026-06-01T12:00:00.000Z
```

- `-n N` / `--number N` — mint N IDs, one per line, all distinct.
- `-p PREFIX` / `--prefix PREFIX` — override the default `R` prefix.
- `--decode` — decode IDs (from args or whitespace-separated stdin) to their
  ISO-8601 UTC minting instant.
- `--help`, `--version` — usage and version (`0.1.0-pre+20260616`).

All times are UTC; IDs are minted only from instants that have already elapsed.

These are the very same `R-XXXX-XXXX` ids that this project uses to track its own
requirements — the design mints one per checkable behavior, so **`idgen` is built
using `idgen`**.

## What's in the box

There is no `cmd/`, `internal/`, or `go.mod` here yet — only the contract and the
harness that turns it into code:

```
project/product/product.md     why idgen exists, for whom, and what it promises (outcomes)
project/design/README.md       the design spine: Conventions, the requirement-id denominator, layout
project/design/INDEX.md        manifest: each Decision → its file, plus a sorted R-id → Decision map
project/design/DNN.md          one self-contained Decision each — seams, interfaces, and its Verification ids
project/plan/README.md         the plan rules (one phase = one package; the done bar)
project/plan/STATUS.md         the manifest: one ⬜/✅ line per phase, the only home of status markers
project/plan/phase-NN.md       one body per phase — objective + the id slice it must cover
project/prompts/{gather,build,verify}.md   the three ralph build-loop prompts
project/README.md              the workspace map + the full build-loop overview
.claude/commands/              the /product-, /research-, /design-, /plan-mode and
                               /create-gather-build-verify-prompts commands that authored the spec
.claude/library/ralph/         the ralph family map skill
```

The `.claude/` commands and skill are version-controlled and ship with the repo on
purpose: the same interactive commands that produced this contract travel with it,
so the *method* is reproducible — not just its output. Clone the project and you
have everything needed both to author the spec and to build from it.

## How the method works

The spec under `project/` is the whole contract; the code is downstream of it. The
three authorities never restate each other:

- **`project/product/product.md`** owns *why* — the problem, the users, and the
  user-facing promises in outcome terms. It never states mechanism.
- **`project/design/`** owns *how* — seams, interfaces, types — and **the
  denominator**: each Decision ends with a Verification list, and every item carries
  a minted `R-XXXX-XXXX` id. That set of ids *is* the enumerated intent the test
  suite is measured against; there is no separate requirements document. The design
  is split for addressability (a spine + one `DNN.md` per Decision + an `INDEX.md`)
  so a build phase reads only the one Decision it realizes.
- **`project/plan/`** owns *construction order* — an append-only list of phases,
  each one package's worth of work, naming the design Decisions and id slice it
  realizes. `STATUS.md` is the only home of the `⬜`/`✅` markers.

`ralph` builds it without holding anything between turns. **gather** is the only
prompt that reads the big docs: it finds the next `⬜` phase and writes a tiny,
self-contained `project/prompts/brief.md`. **build** reads only that brief, writes
code and id-tagged tests (`// R-XXXX-XXXX`), and commits — so **coverage is a grep**.
**verify** is the independent gate: it re-derives the denominator from the brief,
proves every id is covered by a genuine test with `go test -race ./...` green, and
only then flips the phase `⬜ → ✅`. A full overview of the loop — the status
contract, the state machine, and the brief schema — is in
[`project/README.md`](project/README.md).

The spec itself was authored interactively, one decision at a time, by the
`/*-mode` commands in [`.claude/commands/`](.claude/commands/) — product → research
→ design → plan — and the loop prompts by `/create-gather-build-verify-prompts`.
Build and verify are not authoring steps; they belong to `ralph`.
