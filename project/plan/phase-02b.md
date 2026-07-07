# Phase 02b — `internal/idgen`: decode (inverse direction)

*Realizes the decode slice of design Decision 2 (`idgen` public API & prefix
placement); the testing strategy is Decision 7. Depends on Phase 02a.*

The inverse half of the bijection: `TimeOf`, `ErrInvalidID`, the precomputed
modular inverse of the affine multiplier, and the package `init` coprimality
fail-loud panic. Also closes the loop with the round-trip randomized property
sweep, which exercises `MintAt` (Phase 02a) together with `TimeOf`. The observable end
state: `idgen.TimeOf(id)` correctly inverts any id `MintAt` produced, is
prefix-agnostic, and returns an error wrapping `ErrInvalidID` on malformed input.

**Done when** `go test -race ./...` exits 0 and each design Verification id below
is covered by a clearly-named, genuinely-asserting, id-tagged test (a `//` comment
naming the id) in `internal/idgen/*_test.go`:
- R-WH5F-QJYS — round-trip randomized property sweep (ordinary `Test`, large PRNG-seeded `ms` sample) `TimeOf(MintAt(p, t)) == t` over `[Epoch, Epoch+cycle)` and several prefixes.
- R-WN8X-NEO9 — prefix-agnostic decode: `R`/`S`/`SPEC` bodies decode to the same instant.
- R-WPOQ-EY5N — malformed input returns an error wrapping `ErrInvalidID`.
- R-WQWM-SPWC — `TimeOf` robustness sweep (ordinary `Test`, large PRNG-seeded string sample): arbitrary strings never panic (only `ErrInvalidID` or a valid time).
- R-WS4J-6HN1 — coprimality fail-loud: package `init` panics if multiplier and `36⁸` lose coprimality.
