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
- R-WIDC-4BPH — golden vector at `Epoch` → its exact id string (locks the 2026 epoch).
- R-WJL8-I3G6 — golden vector mid-cycle (locks the affine constants and the 4-4 split).
- R-WKT4-VV6V — padding: small ms still yield 8 body chars.
- R-WM11-9MXK — clamping: pre-Epoch instants encode as Epoch (ms 0).
