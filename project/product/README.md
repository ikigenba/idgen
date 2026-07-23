# idgen — Product

**Authority: intent.** This document owns *why idgen exists, for whom, what is in
and out of scope, and what we promise the user* — stated once, in outcome terms.
It does **not** state mechanism, exact formats, exit codes, or test assertions;
those belong to `project/design/README.md`. Where the two could overlap
(behavior), product states the *promise*; design states the *exact, checkable
proof of that promise*.

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
- **Two ways to install.** Build and install from source with the `Makefile`, or
  fetch a prebuilt release binary published for the tagged release. No package
  registries.
- **No external runtime dependencies.** Standard library only.

## Contractual constants

These are promises, not implementation detail, so the product owns them; the
design *uses* them and never re-declares their values:

- **Epoch** — `2026-01-01T00:00:00Z` (UTC). Every ID is anchored to it.
- **Version** — the first release is tagged `v0.1.0`; releases follow Semantic
  Versioning 2.0.0 and are named by their git tag (`vMAJOR.MINOR.PATCH`). The
  version is never a hand-maintained value in the sources — the git tag is its
  only source of truth.

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
- `-p P` / `--prefix P` — **optional.** `P` *replaces* the default prefix `R`
  outright — it never prepends to it, and the default `R` never appears
  alongside a supplied prefix. Omitting `-p` always mints with the `R` prefix,
  never an error:

  ```
  $ idgen -p X
  X-4K7P-9ZQ2
  ```

  A prefix is a non-empty run of letters and/or digits with no separator; an
  explicitly supplied empty or otherwise invalid prefix is rejected (and
  reported — see below).

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

**Error reporting.** Nothing fails silently. Any invalid invocation — a bad or
missing argument, an unrecognized flag, an invalid prefix, a malformed ID to
decode — always prints a message describing what was wrong before the tool
exits with a non-zero status. A user is never left staring at empty output
wondering whether anything ran.

**Help and version.** `--help` (or `-h`) prints usage. `--version` (or `-V`)
prints the version string, and that string identifies the exact source the
binary was built from: a release build reports its git tag (`v0.1.0`); a build
taken between tags or from modified sources says so; a build with no version
information stamped in reports `dev`. No version number is maintained by hand
anywhere in the sources.

**Time correctness.** Every printed time is UTC regardless of the operator's local
zone — local time never leaks into an ID or into decode output. IDs are minted
only from instants that have already elapsed, never the future.

## Success criteria (outcomes)

- A fresh `idgen` mints well-formed IDs and decodes them back, and a freshly
  minted ID round-trips to the instant it was minted from.
- Running `idgen` with no `-p` mints an ID with the `R` prefix; running it with
  `-p P` mints with `P` instead. Neither requires the other flag.
- Asking for several IDs yields distinct IDs, one per line.
- Behavior and output are identical regardless of the operator's time zone.
- A malformed ID on decode is reported without sinking the rest of the batch.
- Any invalid invocation (bad flag, missing required value, invalid prefix,
  malformed decode input) prints a message explaining the problem and exits
  non-zero — it never exits silently with no output.
- The tool builds, tests, and installs from source via the `Makefile`, and a
  tagged release publishes a prebuilt binary a user can fetch and run.
- Asking a release build for its version reports its git tag; asking a build
  with no version stamped in reports `dev`.
