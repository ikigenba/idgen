# Phase 04 — Wire `main`, install & build smoke

*Realizes the build-smoke slice of design Decision 6, the `main` wiring of
Decision 1, and the product's install-from-source promise. Depends on Phase 03b.*

`cmd/idgen/main.go` wires real `os.Args`/stdio/`realClock` into `cli.Run` and
calls `os.Exit` — no logic beyond wiring. With the core (Phase 02) and CLI (Phase
03) complete, the suite is now fully green, so the build-smoke requirement is
reachable. The observable end state: `go install ./cmd/idgen` produces a working
`idgen` binary that mints and decodes end-to-end.

**Done when** `go test -race ./...` exits 0 and:
- R-XMM0-QR6E — build smoke: a clearly-named, id-tagged test asserts the suite is
  green and `go build` produces `bin/idgen` (the Makefile's `build`/`test` targets
  exercised).
- Deterministic install/round-trip checks (not behavioral ids — `main` carries no
  unit test of its own): `go install ./cmd/idgen` succeeds; the built binary prints
  one well-formed id on a bare call, prints `v0.1.0` for `--version`, and
  a mint → `--decode` round trip returns the minting instant.
