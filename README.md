# idgen

This repository is **`idgen` with no source code.** It ships only the documents
and prompts needed to *generate* it: a product brief, a design, a build plan, and
the two [`ralph`](https://github.com/ikigenba/ralph) prompts that drive the build.
Point `ralph` at them and it writes the code, one phase at a time, until the tool
is built and every requirement is proven.

So it is two things at once:

1. **A way to get and build `idgen`** — clone it, run `ralph`, end up with a
   working binary.
2. **A worked demonstration of the `ralph` method** — how a small, fully-specified
   project is built autonomously from *docs + prompts*, with no state held in the
   model's memory between turns.

## Build it with `ralph`

Requires Go 1.26+, `codex`, and the [`ralph`](https://github.com/ikigenba/ralph)
harness.

Clone the repo and, **from the repository root**, run `ralph` against the build
prompt:

```sh
git clone https://github.com/ikigenba/idgen.git
cd idgen
ralph prompts/LOOP.md
```

That single command is the whole build. `ralph` runs `LOOP.md` with a fresh
context per phase, writing one package per turn until the plan reports done — and
you end up with the `cmd/`, `internal/`, and `go.mod` that weren't in the repo.

Then, still from the root, gate the result with an independent fresh-context pass
that proves every design requirement id is covered by a real, id-tagged test and
the suite is green:

```sh
ralph prompts/VERIFY.md
```

You now have source — build and test it the normal way:

```sh
make build      # compile to bin/idgen
make test
```

Always invoke `ralph` from the repository root, so the prompt paths
(`prompts/LOOP.md`, `prompts/VERIFY.md`) and the `docs/` it reads resolve
correctly. `ralph` owns the lifecycle and the budget rails (`--max-spend`,
`--max-time`, …); see its README for the flags. `make fmt`, `make clean`, and
`make install` are available once the source exists.

## What idgen is (the end product)

A small CLI that mints stable spec/requirement IDs of the form `R-XXXX-XXXX` from
the wall-clock instant of minting — and decodes them back to that instant. Every
ID is anchored to a **2026 UTC epoch**, so it stays traceable to the moment it was
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
requirements — the design mints one per checkable behavior (see below), so `idgen`
is built using `idgen`.

## What's in the box

There is no `cmd/`, `internal/`, or `go.mod` here yet — only the contract and the
harness that turns it into code:

```
docs/product.md     why idgen exists, for whom, and what it promises (outcomes)
docs/design.md      how it's built — seams, interfaces, and the denominator:
                    every checkable behavior carries a minted R-XXXX-XXXX id
docs/plan.md        the construction order — one phase = one package
prompts/LOOP.md     ralph builds the next unbuilt phase each turn, tests tagged // R-XXXX-XXXX
prompts/VERIFY.md   ralph proves every requirement id is covered and the suite is green
.claude/commands/   the /product-, /research-, /design-, /plan-mode commands that
                    authored the docs in the first place
```

The `.claude/commands/` modes are version-controlled and ship with the repo on
purpose: the same interactive commands that produced this contract travel with it,
so the *method* is reproducible — not just its output. Clone the project and you
have everything needed both to author the docs and to build from them.

## How the method works

The three documents are the whole contract; the code is downstream of them.

- **`docs/product.md`** owns *why* — the problem, the users, and the user-facing
  promises in outcome terms. It never states mechanism.
- **`docs/design.md`** owns *how* — seams, interfaces, types — and **the
  denominator**: each design Decision ends with a Verification list, and every
  item carries a minted `R-XXXX-XXXX` id. That set of ids *is* the enumerated
  intent the test suite is measured against; there is no separate requirements
  document.
- **`docs/plan.md`** owns *construction order* — an append-only list of phases,
  each one package's worth of work, naming the design Decisions it realizes.

`ralph` (driving `codex`) builds it without holding anything between turns:
`LOOP.md` builds the next unbuilt phase each fresh turn and tags every test with
its requirement id in a `// R-XXXX-XXXX` comment, so **coverage is a grep**.
`VERIFY.md` is the independent gate — it re-derives the denominator from the
design and proves every id is covered by a genuine test with the suite green,
trusting nothing but the files on disk.

The docs themselves were authored interactively, one decision at a time, by the
`/*-mode` commands in [`.claude/commands/`](.claude/commands/) — product →
research → design → plan. Build and verify are not authoring steps; they belong
to `ralph`.
