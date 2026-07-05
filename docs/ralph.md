# Build it with `ralph`

`ralph` is a purpose-built harness used across the ikigenba projects. It runs a
spec's prompt sequence unattended, cycling the prompts in fresh, isolated contexts
until the build is done. This project uses a three-prompt sequence, gather → build →
verify. You point `ralph` at it and walk away.

If you're interested, read [how the spec is structured](spec-structure.md).

> **Why `ralph`.** Beyond running the build hands-off, `ralph` gives friendlier
> progress feedback, lets you mix harnesses and models per prompt file, and manages
> spend limits and max iterations. If you don't already run it, [Claude Code](claude-code.md)
> or [Codex CLI](codex-cli.md) will get you the same result with tooling you likely
> already have.

## Requirements

This method needs:

- **Go 1.26+**, to compile and test the generated code.
- **`git`**, to clone this repository.
- The [`ralph`](https://github.com/ikigenba/ralph) harness, installed and on your
  `PATH`.
- A `ralph`-supported agent (Claude Code or Codex CLI) installed.

## Steps

Clone the repo and, **from the repository root**, run `ralph` against this spec's
three build-loop prompts:

```sh
git clone https://github.com/ikigenba/idgen.git
cd idgen
ralph project/loops/gather.md project/loops/build.md project/loops/verify.md
```

That single command is the whole build. `ralph` cycles this spec's prompt sequence
in fresh, isolated contexts (gather → build → verify → …), building one phase at a
time until `gather` finds no unbuilt phase and reports `DONE`. Each phase is a
logically related, right-sized chunk of work the spec lays out. You end up with the
`cmd/`, `internal/`, and `go.mod` that weren't in the repo, with every design
requirement id covered by an id-tagged test.

You now have source. Build and test it the normal way:

```sh
make build      # compile to bin/idgen
make test       # go test ./...
make install    # go install ./cmd/idgen (onto your PATH via GOBIN)
```

## Notes

- **Always invoke `ralph` from the repository root** so the prompt paths and the
  `project/` spec they read resolve correctly.
- `ralph` owns the lifecycle and the budget rails (`--max-spend`, `--max-time`, …);
  run `ralph --help` for the flags.
