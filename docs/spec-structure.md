# How the spec is structured

This repository is a spec, not a program. Before you pick a build method, it helps
to understand *what* you are handing an agent and *why* it is shaped the way it is.
This page is the tour. It is deliberately narrative; for the authoritative, terse
version of the same mechanics, see [`project/README.md`](../project/README.md),
which every build method reads from.

## The idea in one paragraph

You do not tell the agent *how* to write the code. You give it a contract (what the
tool must do, how its pieces fit, and the order to build them) plus a small set of
prompts that turn that contract into working, tested Go one phase at a time. The
code is downstream of the spec. Change the spec and rebuild, and the code changes
with it. Nothing about the implementation is held in the model's memory between
steps; everything it needs is on disk.

## What's in the box

There is no `cmd/`, `internal/`, or `go.mod` here yet, only the contract and the
prompts that turn it into code:

```
project/product/product.md     why idgen exists, for whom, and what it promises (outcomes)
project/design/README.md       the design spine: Conventions, the requirement-id denominator, layout
project/design/INDEX.md        manifest: each Decision → its file, plus a sorted R-id → Decision map
project/design/DNN.md          one self-contained Decision each: seams, interfaces, and its Verification ids
project/plan/README.md         the plan rules (one phase = one package; the done bar)
project/plan/STATUS.md         the manifest: one ⬜/✅ line per phase, the only home of status markers
project/plan/phase-NN.md       one body per phase: objective + the id slice it must cover
project/prompts/{gather,build,verify}.md   the three build-loop prompts
project/README.md              the workspace map + the full build-loop overview
.claude/commands/              the /product-, /research-, /design-, /plan-mode and
                               /create-gather-build-verify-prompts commands that authored the spec
.claude/library/ralph/         the ralph family map skill
```

The `.claude/` commands and skill are version-controlled and ship with the repo on
purpose: the same interactive commands that produced this contract travel with it,
so the *method* is reproducible, not just its output. Clone the project and you have
everything needed both to author the spec and to build from it.

## Three authorities, one fact each

The spec is split across three authorities, and the discipline that makes it work
is that **they never restate each other**. Each fact lives in exactly one place:

- **`project/product/product.md` owns *why*:** the problem, the users, and the
  user-facing promises, stated as outcomes. It never describes mechanism.
- **`project/design/` owns *how*:** seams, interfaces, types, and, crucially,
  **the denominator**. Each Decision ends with a Verification list, and every item
  on it carries a minted `R-XXXX-XXXX` id. That set of ids *is* the enumerated
  intent the test suite is measured against; there is no separate requirements
  document. The design is split for addressability (a spine, one `DNN.md` per
  Decision, and an `INDEX.md`) so a build phase reads only the one Decision it
  realizes.
- **`project/plan/` owns *construction order*:** an append-only list of phases,
  each one package's worth of work, naming the design Decisions and id slice it
  realizes. `STATUS.md` is the only home of the `⬜`/`✅` markers.

The payoff of tagging every checkable behavior with an id is that **coverage
becomes a grep**: each id-tagged test carries its `// R-XXXX-XXXX` id, so you can
mechanically confirm every requirement is proven rather than trusting a review.

## The build loop: gather → build → verify

The agent builds the project without holding anything between turns. Three prompts
run in a cycle, each in a fresh, isolated context:

1. **gather** is the only prompt that reads the big docs. It finds the next `⬜`
   phase in `STATUS.md` and writes a tiny, self-contained `project/prompts/brief.md`
   for just that phase: the objective and the id slice to cover.
2. **build** reads *only* that brief. It writes the code and the id-tagged tests
   (`// R-XXXX-XXXX`), then commits.
3. **verify** is the independent gate. It re-derives the denominator from the
   brief, proves every id is covered by a genuine test with `go test -race ./...`
   green, and only then flips the phase `⬜ → ✅`.

Then the cycle wraps back to gather for the next phase, and repeats until gather
finds no unbuilt phase and reports `DONE`. Because each step starts clean and reads
only what it needs from disk, the loop is reproducible and the model never has to
remember what happened last turn.

The build methods differ only in what *drives* this cycle: an unattended harness,
or you pasting a prompt into an interactive agent. The contract they read, and the
`DONE` they drive toward, are identical. See
[`project/README.md`](../project/README.md) for the authoritative build-loop
overview: the status contract, the state machine, and the brief schema.

## How the spec itself was authored

The contract was not written by hand in one sitting. It was authored interactively,
one decision at a time, by the `/*-mode` commands in
[`.claude/commands/`](../.claude/commands/), running product → research → design →
plan, with the loop prompts written by `/create-gather-build-verify-prompts`. Build
and verify are not authoring steps; they belong to the build loop. That the
authoring commands ship with the repo is the point: the *method* is meant to be
reused, not just the `idgen` it produced.
