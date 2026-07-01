# Build it with `ralph`

`ralph` is a purpose-built harness that runs the whole gather → build → verify loop
for you, unattended. You point it at the three prompts, and it cycles them in fresh,
isolated contexts until the tool is built. This is the reference path the spec was
designed around.

New here? Read [how the spec is structured](spec-structure.md) first — it explains
what the loop is doing and why. This page assumes it.

> **A note on tooling.** `ralph` is bespoke: it is its own harness with its own
> install and its own flags. If you don't already run it, the interactive
> [Claude Code](claude-code.md) or [Codex CLI](codex-cli.md) paths reach the same
> result with tooling you may already have. Reach for `ralph` when you want the
> loop to run fully hands-off.

## Requirements

On top of the repository-wide [prerequisites](../README.md#prerequisites) (Go 1.26+
and `git`), this method needs:

- The [`ralph`](https://github.com/ikigenba/ralph) harness, installed and on your
  `PATH`.
- A `ralph`-supported agent (for example `codex`) configured for `ralph` to drive.

## Steps

Clone the repo and, **from the repository root**, run `ralph` against the three
build-loop prompts:

```sh
git clone https://github.com/ikigenba/idgen.git
cd idgen
ralph project/prompts/gather.md project/prompts/build.md project/prompts/verify.md
```

That single command is the whole build. `ralph` cycles the three prompts in fresh,
isolated contexts — **gather → build → verify → …** — writing one package per phase
until `gather` finds no unbuilt phase and reports `DONE`. You end up with the
`cmd/`, `internal/`, and `go.mod` that weren't in the repo, with every design
requirement id covered by an id-tagged test.

You now have source — build and test it the normal way:

```sh
make build      # compile to bin/idgen
make test       # go test ./...
```

## Notes

- **Always invoke `ralph` from the repository root** so the prompt paths and the
  `project/` spec they read resolve correctly.
- `ralph` owns the lifecycle and the budget rails (`--max-spend`, `--max-time`, …);
  see [its README](https://github.com/ikigenba/ralph) for the flags.
