# Phase 01 — Scaffolding: module, three seams & Makefile

*Realizes design Decisions 1 (package layout) and 6 (Makefile) — structural; no
behavioral ids. Depends on no earlier phase.*

Create the Go module (`github.com/ai4mgreenly/idgen`, Go 1.26) and the three seams
as a buildable skeleton, plus the build tooling. The observable end state is a
compiling skeleton with no behavior yet:

```
cmd/idgen/main.go      thin stub + `var version = "dev"` linker seam (no behavior yet)
internal/cli/          Run(args, stdin, stdout, stderr, version, clock) int  (stub returning 0)
internal/idgen/        pure encode/decode + Epoch + Clock seam       (stubs)
Makefile               build · fmt · test · install · clean (mirrors embed's)
```

`Epoch` and the exported `MintAt`/`TimeOf`/`Run` signatures exist (so later phases
fill bodies, not shapes), but carry only stub bodies; this phase asserts no
runtime behavior. `cmd/idgen/main.go` declares the `var version = "dev"` seam (D6)
so the Makefile's `-ldflags -X main.version=...` target resolves. The `Makefile`
mirrors `embed`'s (D6): `build` stamps the git-derived version and writes
`bin/idgen`; `install` copies the built binary under `PREFIX` (`$HOME/.local` by
default); `fmt` is `go fmt ./...`.

**Done when** — all deterministic, structural checks (this phase carries no
behavioral ids):
- `go build ./...` exits 0 (the skeleton compiles).
- `make build` produces `bin/idgen` (stamped via the ldflags seam).
- The `Makefile` defines exactly the five targets `build`, `fmt`, `test`,
  `install`, `clean` (each present in `.PHONY`).
- The three seam directories exist with the package files named above
  (`cmd/idgen/main.go`, `internal/cli/`, `internal/idgen/`).
