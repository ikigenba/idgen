# How the spec is structured

This repository is a spec, not a program. Before you pick a build method, it helps
to understand *what* you are handing an agent and *why* it is shaped the way it is.
This page is the tour. It is deliberately narrative; the authoritative spec
shapes live in the `spec-shapes` skill, while
[`project/README.md`](../project/README.md) is the thin workspace map every build
method starts from.

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
project/product/README.md      why idgen exists, for whom, and what it promises (outcomes)
project/design/README.md       the design spine: Conventions, the requirement-id denominator, layout
project/design/INDEX.md        manifest: each Decision → its file, plus a sorted R-id → Decision map
project/design/DNN.md          one self-contained Decision each: seams, interfaces, and its Verification ids
project/plan/README.md         the plan rules (one phase = one package; the done bar)
project/plan/STATUS.md         the manifest: one ⬜/✅ line per phase, the only home of status markers
project/plan/phase-NN.md       one body per phase: objective + the id slice it must cover
project/loops/{gather,build,verify}.md   the three build-loop prompts
project/loops/audit.md         the single prompt of the separate, optional coverage-audit loop
project/loops/run              the operator entrypoint: a shell wrapper that launches ralph on the build loop
project/loops/README.md        the installed loops' overview: status contract, state machine, brief schema
project/README.md              the workspace map (thin: the folder table + pointers)
.claude/                       the spec-process skills and commands, packaged for Claude Code
.agents/skills/                the same spec-process skills, packaged for the Codex CLI
```

The workflow tooling is version-controlled and ships with the repo on purpose,
packaged for both Claude Code (`.claude/`) and the Codex CLI (`.agents/skills/`):
the same method that produced this contract travels with it, so the *method* is
reproducible, not just its output. Clone the project and you have everything
needed both to author the spec and to build from it, with either agent.

## Three authorities, one fact each

The spec is split across three authorities, and the discipline that makes it work
is that **they never restate each other**. Each fact lives in exactly one place:

- **`project/product/README.md` owns *why*:** the problem, the users, and the
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
   phase in `STATUS.md` and writes a tiny, self-contained `project/loops/brief.md`
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
[`project/loops/README.md`](../project/loops/README.md) for the authoritative
build-loop overview: the status contract, the state machine, and the brief schema.

## The audit loop: adversarial coverage checking (optional)

Verify already gates each phase, but it asks "is every id covered by a test?"
A separate, single-prompt loop (`project/loops/audit.md`) asks the harder
question after the build: **could each tagged test actually fail?** It walks the
design one Decision at a time, judges every `R-XXXX-XXXX` id's test against the
behavior the design states (escalating to a throwaway-worktree mutation test
when reading can't settle it), and writes its findings to a gitignored
`project/audit/REPORT.md`. It never modifies the live checkout. It is run on
demand — `ralph project/loops/audit.md` — and is documented alongside the build
loop in [`project/loops/README.md`](../project/loops/README.md).

## How the spec itself was authored

The contract was not written by hand in one sitting. It was settled interactively —
the goal discussed in conversation, then interrogated one question at a time with
the `grillme` workflow — and written in one pass by the `codify` workflow, to the
artifact shapes defined in the `spec-shapes` skill, with the loop prompts generated
by `create-gather-build-verify-prompts` and `create-audit-prompts`. Build and
verify are not authoring steps; they belong to the build loop. Those workflows ship
with the repo for both agents — as skills and commands under
[`.claude/`](../.claude/) for Claude Code, and as skills under
[`.agents/skills/`](../.agents/skills/) for the Codex CLI — and that is the point:
the *method* is meant to be reused, not just the `idgen` it produced.

To make a spec change of your own, the recipe is the same in either agent:

1. Start a session with the spec-process context loaded — `/skillset spec` in
   Claude Code, or `use $spec-shapes` in Codex. (Codex has no slash commands;
   its shipped skills are invoked by `$name` reference: `$grillme`, `$codify`.)
2. Describe the change in plain English.
3. Get grilled and answer the questions, one at a time, until the goal is
   settled. (In Claude Code the grillme workflow rides along with
   `/skillset spec` — just ask to be grilled; in Codex, invoke `$grillme`.)
4. Ask the agent to `codify` — it writes product, research, design, and plan in
   one pass.
5. From a shell, run `project/loops/run` to rebuild the code from the updated
   spec.

Worked examples of this recipe — both greenfield and adding a feature — live in
[authoring with Claude Code](authoring-claude-code.md) and
[authoring with the Codex CLI](authoring-codex-cli.md).
