# idgen — Plan

**Authority: construction order and history.** This document owns *the order idgen
is built in* and *the record of what has been built*. Unlike `product` and
`design` — which are rewritten in place to stay authoritative for the current
state — this plan is **append-only**: phases are added at the bottom and marked
done as they land; completed phases are never rewritten or deleted, so the plan
doubles as the construction history. To extend idgen, update `product` and
`design` in place, then append a new phase here.

**One phase = one package = one accumulating context.** Each phase is a single
coherent unit — almost always one package — built in one accumulating context
against `product` and `design`, reading only that package's design Decisions and
the *interfaces* (not internals) of the packages it depends on. This is what keeps
every phase the size of a small standalone tool no matter how large the whole
project grows; nothing ever has to hold more than one package at once.

A phase is **done** when every Verification item in the design Decisions it
realizes is covered by a clearly-named test and `go test -race ./...` is green
(see the design's *Verification & "done"* section). The phases below are ordered so
each depends only on earlier ones.

## Status

Not started. The workspace holds `product`, `design`, and this plan; no code yet.

## Phases

### Phase 1 — Scaffolding: module, seams & Makefile · ⬜ not started
*Realizes design Decision 1 (layout) and Decision 6 (Makefile).*

Create the Go module (`github.com/ai4mgreenly/idgen`, Go 1.26) and the three seams
as a buildable skeleton, plus the build tooling:

```
cmd/idgen/main.go      thin stub (no behavior yet)
internal/cli/          run(args, stdin, stdout, stderr, clock) int  (stub)
internal/idgen/        pure encode/decode + Epoch + Clock seam       (stub)
Makefile               build · test · clean · fmt · install
```

**Done when:** `go build ./...` is clean, `make build` yields `bin/idgen`, and all
five Makefile targets exist.

### Phase 2 — `internal/idgen` package: encoding core · ⬜ not started
*Realizes design Decisions 2 and 7 (idgen tier). Depends on Phase 1.*

The pure core, with no I/O and no clock: `MintAt`, `TimeOf`, `Epoch`,
`ErrInvalidID`, the affine bijection and base-36 4-4 body math, and the `init`
coprimality fail-loud panic. Built against Decision 2's interface alone.

**Done when** its Verification items are covered — independently-derived golden
vectors (lock epoch + constants), padding, pre-epoch clamping, prefix-agnostic
decode, malformed → `ErrInvalidID`, `FuzzRoundTrip`, `FuzzTimeOf` — and
`go test ./internal/idgen/...` is green.

### Phase 3 — `internal/cli` package: the CLI core · ⬜ not started
*Realizes design Decisions 3, 4, 5, and 6 (version/usage). Depends on Phase 2.*

The whole CLI in one context, consuming `internal/idgen` only through `MintAt` /
`TimeOf`: `run(...)` with one `flag.FlagSet`; `--decode` dispatch and the `0`/`1`/`2`
exit-code taxonomy (D4); the `Clock` seam and distinct-millisecond mint wait loop
(D3); default/`-n`/`-p` mint with prefix and number validation (D5); decode with
args-then-stdin precedence, partial-failure → 1, empty → 0, UTC output (D5); the
version `var` and usage text (D6).

**Done when** its Verification items are covered by table-driven `run(...)` tests
over in-memory buffers with a **fake `Clock`** (no subprocess, no real sleeps) —
dispatch/help/version/exit matrix, mint count/prefix/distinctness + virtual-time
advance, decode args-vs-stdin/partial/empty/`TZ`-independence, end-to-end
round-trip through `run` — and `go test ./internal/cli/...` is green.

### Phase 4 — Wire `main` & install · ⬜ not started
*Realizes design Decision 1 (main) and the product's install-from-source.
Depends on Phase 3.*

`cmd/idgen/main.go` wires real `os.Args`/stdio/`realClock` into `cli.Run` and calls
`os.Exit` — no logic beyond wiring.

**Done when:** `go install ./cmd/idgen` succeeds; the built binary prints one
`R-XXXX-XXXX` on a bare call, prints `0.1.0-pre+20260616` for `--version`, and a
mint → `--decode` round trip holds; `go test -race ./...` is green.
