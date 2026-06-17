# idgen — Product

**Authority: intent.** This document owns *why idgen exists, for whom, what is in
and out of scope, and what we promise the user* — stated once, in outcome terms.
It does **not** state mechanism, exact formats, exit codes, or test assertions;
those belong to `idgen-design.md`. Where the two could overlap (behavior), product
states the *promise*; design states the *exact, checkable proof of that promise*.

## Problem

Specs and requirements need stable, unique identifiers you can embed in prose,
source comments, and test names so an item stays traceable to exactly one moment
of intent. Inventing IDs by hand drifts — collisions, inconsistent shapes, no way
back to "when was this written?" No small, standalone tool does only one thing:
mint an ID.

## Purpose

`idgen` is a single-purpose developer CLI that mints short, unique IDs of the
shape `R-XXXX-XXXX`, each derived from the wall-clock instant it was minted. It
also decodes an ID back to the UTC instant it came from, so an ID stays traceable
to its moment forever.

## Users

Developers and spec authors who want to drop a fresh, unique ID into a
requirement, a comment, or a test name — and occasionally ask "when was this
minted?" They run it from a terminal or pipe it into other tools.

## Scope

- **Single purpose.** Mint IDs; decode them back to instants. Nothing else.
- **Install from source only.** No binary distribution, no package registries; a
  `Makefile` builds, tests, and installs it.
- **No external runtime dependencies.** Standard library only.

## Contractual constants

These are promises, not implementation detail, so the product owns them; the
design *uses* them and never re-declares their values:

- **Epoch** — `2026-01-01T00:00:00Z` (UTC). Every ID is anchored to it.
- **Version** — starts at `0.1.0-pre+20260616`, following Semantic Versioning
  2.0.0.

## What we promise (user-facing behavior)

**Mint (default).** The bare command mints and prints one fresh ID and exits
successfully:

```
$ idgen
R-4K7P-9ZQ2
```

- `-n N` / `--number N` — mint N IDs, one per line, **all distinct**. Because an
  ID encodes the millisecond of minting, the tool produces at most one distinct
  ID per millisecond; asking for several therefore takes a brief moment, and the
  command does not return until the last ID's millisecond has elapsed.
- `-p P` / `--prefix P` — use `P` in place of the default `R` prefix. A prefix is
  a non-empty run of letters and/or digits with no separator; empty or otherwise
  invalid prefixes are rejected.

**Decode (`--decode`).** Turns IDs back into the UTC instant each was minted from:

```
$ idgen --decode R-4K7P-9ZQ2
2026-06-01T12:00:00.000Z
```

- IDs come from positional arguments, or from a stream on stdin (whitespace
  separated); arguments take precedence.
- Any prefix decodes — the instant lives in the body, so an ID is readable
  regardless of which prefix or tool produced it.
- A malformed ID is reported and the rest of the batch still decodes; the run
  ends with a non-zero status if anything was malformed.

**Help and version.** `--help` prints usage; `--version` prints the version
string.

**Time correctness.** Every printed time is UTC regardless of the operator's local
zone — local time never leaks into an ID or into decode output. IDs are minted
only from instants that have already elapsed, never the future.

## Success criteria (outcomes)

- A fresh `idgen` mints well-formed IDs and decodes them back, and a freshly
  minted ID round-trips to the instant it was minted from.
- Asking for several IDs yields distinct IDs, one per line.
- Behavior and output are identical regardless of the operator's time zone.
- A malformed ID on decode is reported without sinking the rest of the batch.
- The tool builds, tests, and installs from source via the `Makefile`.
