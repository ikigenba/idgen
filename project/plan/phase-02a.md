# Phase 02a — `internal/idgen`: encode (mint direction)

*Realizes the encode slice of design Decision 2 (`idgen` public API & prefix
placement); the testing strategy is Decision 7. Depends on Phase 01.*

The forward half of the bijection, with no I/O and no clock: `Epoch`, the affine
forward map, base-36 4-4 body encode, padding, and clamping — `idgen.MintAt(prefix,
t)`. The observable end state: `MintAt` produces the exact, independently-derived
id string for any instant at or after `Epoch`, prefix supplied by the caller.

**Done when** `go test -race ./...` exits 0 and each design Verification id below
is covered by a clearly-named, genuinely-asserting, id-tagged test (a `//` comment
naming the id) in `internal/idgen/*_test.go` (golden vectors **derived
independently/offline**, not snapshotted from the code under test):
- R-WIDC-4BPH — golden vector at `Epoch` (ms 0) → its exact id string; pins the offset constant (does not lock the epoch — see R-WJL8).
- R-WJL8-I3G6 — golden vector at a fixed absolute post-epoch instant (literal, not the `Epoch` symbol) → exact string; locks the affine constants/4-4 split **and** the 2026 epoch.
- R-WKT4-VV6V — padding: small ms still yield 8 body chars.
- R-WM11-9MXK — clamping: pre-Epoch instants encode as Epoch (ms 0).
