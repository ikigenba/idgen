# idgen

This repository is **`idgen` with no source code.** It ships the spec that
*generates* it, and nothing else. Point an AI coding agent at that spec and it
writes the code, tests it, and proves it against the spec.

*Just want to build it? See the [build instructions](#how-to-build-it) below.*

So it is two things at once:

1. **A way to get and build `idgen`:** clone it, run one of the build methods
   below, end up with a working binary.
2. **A demonstration of spec-first construction:** a small project fully specified
   up front, then built from that spec rather than typed by hand. See
   [how the spec is structured](docs/spec-structure.md) for how that works.

## What idgen is (the end product)

Specs and the code that implements them need stable, traceable IDs to embed in
requirements, comments, and test names. These only need to be unique within one
project, not globally, so they can be much shorter than a UUID.

`idgen` is a small CLI that mints those IDs in the form `R-XXXX-XXXX` from the
millisecond it was minted, and decodes them back to that same millisecond. Every ID
is anchored to a **2026 UTC epoch**, so it stays traceable to the moment it was
minted.

The millisecond count isn't shown raw. It runs through a reversible affine bijection
(a modular multiply-and-add) before base-36 encoding, so consecutive milliseconds
land far apart and the IDs look random and scattered. Because the map is one-to-one,
`--decode` inverts it exactly to recover the original millisecond.

```sh
$ idgen
R-4K7P-9ZQ2

$ idgen --decode R-4K7P-9ZQ2
2026-06-01T12:00:00.000Z
```

- `-n N` / `--number N`: mint N IDs, one per line, all distinct.
- `-p PREFIX` / `--prefix PREFIX`: override the default `R` prefix.
- `--decode`: decode IDs (from args or whitespace-separated stdin) to their
  ISO-8601 UTC minting time, to the millisecond.
- `--help`, `--version`: usage and version (`0.1.0-pre+20260616`).

All times are UTC; IDs are minted only from a millisecond that has already elapsed.

These are the very same `R-XXXX-XXXX` ids that this project uses to track its own
requirements. The design mints one per checkable behavior, so **`idgen` is built
using `idgen`**.

## Prerequisites

Every build method produces the same Go program, so two things are required no
matter which path you take:

- **Go 1.26+**, which the spec targets; the generated code and tests need it to
  compile and run.
- **`git`**, to clone this repository.

Each method then needs its own agent tooling on top of that baseline. The
per-method requirements are listed in each document linked below.

## How to build it

The spec is the same for every path; the methods differ only in *who drives the
loop* and *how much is automated*. Pick one:

- **[`ralph`](docs/ralph.md)**: a purpose-built harness used across the ikigenba
  projects. Point it at the three prompts and walk away; it runs the spec's build
  loop unattended. If you don't already use it, one of the interactive paths below
  is likely the better choice.
- **[claude code](docs/claude-code.md)**: drive the same loop interactively with
  Claude Code. Paste one orchestration prompt and it cycles the phases for you.
  A good default for most people.
- **[codex cli](docs/codex-cli.md)**: the same interactive loop, driven by the
  Codex CLI instead. Also a good default if Codex is your agent.

New to the repo? Read **[how the spec is structured](docs/spec-structure.md)**
first. It explains what you are actually building, why the spec is split the way it
is, and how the gather → build → verify loop turns it into code. Every method above
assumes that background.
