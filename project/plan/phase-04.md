# Phase 04 — Wire `main`, install, release & smoke

*Realizes the version-seam and build-smoke slice of design Decision 6, the `main`
wiring + version linker seam of Decision 1, and the product's install-from-source
and prebuilt-release promises. Depends on Phase 03b.*

`cmd/idgen/main.go` wires real `os.Args`/stdio/`realClock` into `cli.Run` and
calls `os.Exit` — no branching logic beyond wiring — and owns the
`var version = "dev"` linker seam (D6), threading its stamped value into
`cli.Run`. With the core (Phase 02) and CLI (Phase 03) complete, the suite is now
fully green, so the version smokes and the build smoke are reachable.

This phase also lands the three repo-root release files (D6), targeting the
GitHub remote `ikigenba/idgen`:

```
.goreleaser.yaml               goreleaser v2: cross-compile linux/darwin × amd64/arm64,
                               versionless tar.gz + checksums, -X main.version={{.Tag}}
.github/workflows/release.yml  on push tag v*: setup-go 1.26 + goreleaser release --clean
install.sh                     curl installer: fetch releases/.../idgen_<os>_<arch>.tar.gz
```

The observable end state: `make install` produces a working, version-stamped
`idgen` under `PREFIX`; a `go build -ldflags` build reports its stamped version;
and pushing a `vMAJOR.MINOR.PATCH` tag is all a release takes.

**Done when** `go test -race ./...` exits 0 and:
- R-TJYQ-4NKS — version smoke (unstamped): a clearly-named, id-tagged test in
  `cmd/idgen/main_test.go` asserts the unstamped `go test` binary (its
  `main.version` defaulting to `dev`) prints exactly `dev` for `--version`/`-V`,
  exit 0.
- R-TL6M-IFBH — version smoke (real seam): a clearly-named, id-tagged test builds
  `cmd/idgen` with `go build -ldflags "-X main.version=<sentinel>"`, runs the
  produced binary with `--version`, and asserts stdout is exactly the sentinel —
  exercising the real linker seam end to end.
- R-XMM0-QR6E — build smoke: a clearly-named, id-tagged test asserts the suite is
  green and `go build` produces `bin/idgen` (the Makefile's `build`/`test` targets
  exercised).
- Deterministic install & release-plumbing checks (not behavioral ids — `main`
  carries no branching logic of its own):
  - `make install` succeeds and places `idgen` under `$PREFIX/bin`; the installed
    binary prints one well-formed id on a bare call and a mint → `--decode` round
    trip returns the minting instant.
  - the three release files exist at their paths and carry their contractual
    values (a `project/`-excluded grep): `.goreleaser.yaml` names
    `project_name: idgen`, `main: ./cmd/idgen`, `-X main.version=`, and
    `owner: ikigenba`; `install.sh` sets `REPO="ikigenba/idgen"`;
    `.github/workflows/release.yml` triggers on tag `v*`.
