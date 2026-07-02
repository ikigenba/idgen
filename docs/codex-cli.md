# Build it with the Codex CLI

You can drive the gather → build → verify loop interactively with the Codex CLI
instead of a dedicated harness. You paste one orchestration prompt, and Codex spawns
a subagent per step, relays each step's status, and advances the cycle until it
reports `DONE`. This is a good default if Codex is your agent.

New here? Read [how the spec is structured](spec-structure.md) first — it explains
what the loop is doing and why. This page assumes it.

## Requirements

This method needs:

- **Go 1.26+** — to compile and test the generated code.
- **`git`** — to clone this repository.
- The Codex CLI, installed and authenticated.

Then clone the repo and start Codex from the repository root:

```sh
git clone https://github.com/ikigenba/idgen.git
cd idgen
codex
```

## Two ways to run the loop

Both prompts drive the same gather → build → verify cycle; they differ only in
*where the current step is remembered*.

- The **external-cursor** version keeps the current step in a file on disk
  (`project/prompts/cursor.md`). The step survives a context wipe, a crash, or you
  closing Codex and resuming later — the orchestrator just re-reads the file and
  picks up where it left off. **This is the preferred version** for a real build.
- The **internal-cursor** version keeps the current step in the orchestrator's own
  running context and continues the goal until `DONE`. There is no bookkeeping file
  to reason about, which makes it the simpler mental model — but the run has to
  complete in one sitting.

Pick one, paste it into Codex, and let it run.

### External cursor (preferred)

Keeps the loop position in `project/prompts/cursor.md`, so the build is resumable.

```
/goal Advance the gather -> build -> verify prompt cycle until DONE (verify wraps
  back to gather), using ./project/prompts/cursor.md as the durable current-step
  marker. Each iteration: read cursor.md for the current step; if it is missing or
  empty, the current step is gather. Spawn a Codex subagent to read and execute that
  step's prompt file (./project/prompts/{gather,build,verify}.md); it ends by
  emitting a single JSON object of the form {"status": ..., "message": ...}. Parse
  its "status" and relay its "message" so progress stays visible, then act:
  CONTINUE -- leave cursor.md unchanged and continue the goal with the same step.
  NEXT -- write the next step in the cycle to cursor.md and continue the goal.
  DONE -- delete cursor.md and stop the goal. If a subagent fails, or its final
  output is not a JSON object carrying a "status" of exactly CONTINUE, NEXT, or
  DONE, stop and surface the raw output -- never guess a status, and leave cursor.md
  unchanged. You are only the orchestrator: you read and write cursor.md for
  bookkeeping, but never read the prompt files themselves -- subagents read and
  execute those.
```

### Internal cursor (easier to reason about)

Keeps the loop position in the orchestrator's context and runs to `DONE` in one go.

```
/goal Cycle the prompt sequence gather -> build -> verify (verify wraps back to gather),
  starting at gather, until a subagent returns a DONE status. For each step, spawn a
  Codex subagent to read and execute the current prompt file
  (./project/prompts/{gather,build,verify}.md); it ends by emitting a single JSON
  object of the form {"status": ..., "message": ...}. Parse its "status" and relay
  its "message" so progress stays visible. CONTINUE -- re-run the same prompt file.
  NEXT -- advance to the next prompt file (verify wraps to gather). DONE -- stop. If
  a subagent fails, or its final output is not a JSON object carrying a "status" of
  exactly CONTINUE, NEXT, or DONE, stop and surface the raw output -- never guess a
  status. You are only the orchestrator -- never read the prompt files yourself;
  subagents read and execute them.
```

## When it finishes

Either way, when the loop reports `DONE` the `cmd/`, `internal/`, and `go.mod` that
weren't in the repo now exist, with every design requirement id covered by an
id-tagged test. Build and test the result the normal way:

```sh
make build      # compile to bin/idgen
make test       # go test ./...
make install    # go install ./cmd/idgen (onto your PATH via GOBIN)
```
