# idgen

This repository is **`idgen` with no source code.** It ships the spec that
*generates* it, and nothing else. Point an AI coding agent at that spec and it
writes the code, tests it, and proves it against the spec.

*Just want to build it? See the [build instructions](#how-to-build-it) below.*

So it is two things at once:

1. **A way to get and build `idgen`:** clone it, run one of the build methods
   below, end up with a working binary.
2. **A demonstration of spec-first construction:** a small project fully specified
   up front, then built from that spec. See
   [how the spec is structured](docs/spec-structure.md).

## What idgen is (the end product)

Some projects, including the ikigenba projects, embed stable, traceable IDs in their
requirements, comments, and test names. Those IDs only need to be unique within the
project, not globally, so they can be much shorter than a UUID.

`idgen` is a small CLI that mints those IDs in the form `R-XXXX-XXXX`. Each ID
encodes the number of milliseconds from a **2026 UTC epoch** to the moment it was
minted.

The millisecond count isn't shown raw. It runs through a reversible affine bijection
(a modular multiply-and-add) before base-36 encoding, so consecutive milliseconds
land far apart and the IDs look random and scattered. Because the map is one-to-one,
every distinct millisecond yields a distinct ID.

The body uses base-36, digits `0-9` and uppercase `A-Z` only. That keeps an ID
case-insensitive, so it survives being lowercased in a URL, said aloud, or retyped
by hand without a normalization step. It stays identifier- and URL-safe, so it
embeds in comments, test names, and links with no escaping. And it stays short:
eight base-36 characters hold 36⁸, about 2.8 trillion values, which is roughly 89
years of milliseconds past the epoch, in a fixed-width body that splits cleanly
4-4.

```sh
$ idgen
R-4K7P-9ZQ2
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

Whichever path you pick, the spec also ships an optional follow-up: a
single-prompt **audit loop** that adversarially re-checks the finished build,
asking of every requirement id not just "is there a test?" but "could that test
actually fail?".

New to the repo? Read **[how the spec is structured](docs/spec-structure.md)**
first. It explains what you are actually building, why the spec is split the way it
is, and how the gather → build → verify loop turns it into code. Every method above
assumes that background.
