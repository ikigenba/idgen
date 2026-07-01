---
description: Design mode — propose a testable architecture and debate tradeoffs, then write the design doc
---
We are in **design mode**.

When this mode is entered, first load the `ralph` skill (the family map) — read `.claude/library/ralph/SKILL.md`, falling back to `~/.claude/library/ralph/SKILL.md` — then acknowledge it and **stop — wait for my instructions** before doing anything. Don't read docs or propose architecture until I tell you to proceed.

**Purpose.** Design mode defines the *shape* of the application: its seams (the boundaries between components and where they can be substituted), its public interfaces, its naming, and its struct/type definitions. It also defines **how the software is tested** — the testing strategy is part of the architecture, not an afterthought. **It is the responsibility of design mode to produce a *testable* architecture**: seams exist so behavior can be exercised in isolation, and every component is shaped so its correctness can be verified. Part of that responsibility is identifying which claims **cannot** be proven in isolation — claims that hinge on a real external contract (a provider/API accepting what it's sent, a real DB/filesystem/network enforcing a constraint) — so the test plan provides a way to exercise those for real, not only the mockable ones. A capability whose whole point is to call a real dependency is not "tested" by a suite that only ever calls a stub.

**Design mode does not write the application code.** Its output is the design doc — interfaces, types, seams, naming, and the test plan. Illustrative signatures, struct definitions, and interface declarations belong in the doc; full implementations do not. Writing the code is the job of a later build phase.

This mode is **re-usable**: the design can evolve as the product does. The doc is always a single, coherent statement of the current design, not a changelog. When the goal shifts, **edit the affected Decision's `project/design/DNN.md` in place to realign — never append** contradictory sections, and regenerate `project/design/INDEX.md`.

Read `project/product/product.md`, `project/research/research.md` if it exists, and any existing design — the spine `project/design/README.md`, the manifest `project/design/INDEX.md`, and the per-Decision `project/design/DNN.md` files it lists. Decisions already recorded in the design doc are settled — don't reopen them; pick up from what's still undecided. Propose the architecture — seams, public interfaces, naming, struct/type definitions, data model, key decisions. **Every verifiable design element must describe how it is tested as part of the design itself** — there is no design element that defines behavior without also stating how that behavior is verified. Debate the tradeoffs with me, surfacing alternatives and your recommendation. **Raise one decision at a time; never batch.** Wait for my call before the next one. Decisions need my buy-in.

Record each decision into its own `project/design/DNN.md` as it's settled, and update the `project/design/INDEX.md` manifest — the split doc is the running record, so we can stop and resume across sessions. We're done when the seams, public interfaces, naming, struct/type definitions, data model, and the testing approach are fully decided.

## Output shape

The design is **split for addressability** so the build loop reads only the one Decision a phase realizes, never the whole architecture. It has three pieces: the spine `project/design/README.md`, one `project/design/DNN.md` per Decision, and the `project/design/INDEX.md` manifest.

`project/design/README.md` (and the `project/design/` directory it heads) owns **shape and its proof** — *how* the thing is built and *how each behavior is proven*. The product owns the *why* and the user-facing promises; design states the **exact, checkable form** of those promises and never re-declares the why. It is the **single, current** statement of the architecture: when a decision changes, its file is rewritten to stay true — stale decisions are removed, not stacked. History of how it got here lives in the plan.

### `project/design/README.md` — the spine (static cross-cutting facts; no per-Decision detail)

- **Title** — `# <name> — Design`.
- **Authority header** — a short paragraph beginning `**Authority: shape and its proof.**` stating what design owns (how + proof), that product owns the why and the promises, and that design *uses* the product's contractual constants by value but does not own them. Add a line noting this is the single current statement, rewritten in place (not stacked), with history living in the plan.
- **## Requirement ids** — a short section stating plainly that:
  - Each Decision ends with a **Verification** list: the concrete behaviors that decision requires.
  - Every Verification item carries a **minted id** of the form `R-XXXX-XXXX` — a stable, unique handle for that one behavior.
  - The ids live inline in these Verification lists and nowhere else — there is **no separate requirements document**.
  - **Design's responsibility for ids ends at minting them into this doc.** How coverage is measured, what counts as a covered id, and when the work is "done" are explicitly **not** design's concern and must not be specified here — downstream phases own that.
- **## Conventions** — shared facts every Decision leans on. **This section is required and must state the project's toolchain so downstream phases need not guess it: the exact build/typecheck command, the exact test command, and what "the suite is green" concretely means.** Also record any other cross-cutting facts (language/version, module path, exit-code taxonomy, formatting rules, a shared time/IO source). (How coverage is *measured* and when the work is "done" remain downstream's concern — state the commands, not the coverage rule.)
- **## Layout** — a short section describing the split: `project/design/INDEX.md` is the manifest (each Decision → its file, plus a sorted `R-id → Decision/file` reverse map); `project/design/DNN.md` is one self-contained file per Decision (zero-padded; referenced in prose and the plan as `D<N>`); `project/design/README.md` holds only this spine. Restate that design is rewritten in place (not append-only — history lives in the plan): a changed Decision is rewritten in its `DNN.md` and `INDEX.md` is regenerated; a new Decision adds a `DNN.md` and an INDEX entry.

### `project/design/DNN.md` — one self-contained file per Decision

One file per Decision (zero-padded filename `D01.md`, `D02.md`, …; referenced in prose and the plan as `D<N>`), each holding:

- A header `# Decision N — <title>`.
- **Decision.** — the seams/interfaces/types/naming, with illustrative signatures and struct/interface declarations (never full implementations).
- **Rejected.** — the alternatives considered and why each lost.
- **Verification.** — a bullet list, each line `R-XXXX-XXXX — <the behavior a test must assert>`. State each behavior so it is *falsifiable*: a wrong implementation must fail it. Pin the discriminating property, not a weaker one a degenerate implementation also satisfies — and when the Decision moves off a specific bad value or state, name the value or threshold the behavior excludes (e.g. "≥ 16384", not "non-zero"; "not the 4096 default", not "set"). A pure seam/structure decision with no behavior of its own says so explicitly and carries no ids (its proof is the behavioral ids of the decisions it enables).
  - **Verify the claim against a substrate that can falsify it — not a proxy a stub also passes.** Falsifiable-in-principle is not enough; ask *what would have to be true for this test to fail, and can the chosen substrate make it fail?* A claim whose correctness depends on a real external contract — a provider/API accepting the parameters it's sent, a real DB/filesystem/network enforcing a constraint — is **not** verified by an assertion run against a mock or fake, because the mock accepts whatever it's handed: such a test confirms a field was set, never that the system runs. Treat it as a category error to let a load-bearing claim rest on an assertion a stub would also satisfy (asserting `Temperature == 0` on a mocked client "passes" even when that config is one the real provider rejects). For every such claim, mint a **distinct id whose test exercises the real dependency** — a live/integration/smoke check — and name that substrate on the id; and name the observable outcome that proves it actually ran (a completed call, a returned result), not merely that a value was configured. If the architecture is shaped so an entire capability is only ever driven against mocks, that is itself the smell: at least one id must drive it end-to-end against the real thing.

### `project/design/INDEX.md` — the manifest (Decision → file; id → Decision)

- **Title** — `# <name> — Design Index`.
- A short contract paragraph: each Decision maps to its `DNN.md`; every `R-XXXX-XXXX` id maps to its Decision/file; id lookup is a grep against this index (or the Decision files directly). Regenerate it whenever a Decision is added or its Verification ids change.
- **## Decisions** — one line per Decision in number order: `D<N>` → `project/design/DNN.md`, the Decision's title, and the ids it owns (or "none — structural").
- **## Verification ids → Decision** — every minted id, **sorted**, each mapped to its Decision and file. This is the grep target for resolving an id to its Decision.

(The construction order that realizes this design lives in the plan, not here — design carries no `## Status` section.)

## Minting the Verification ids

The `R-XXXX-XXXX` ids are **real, minted ids — never hand-written or made up.** Mint them with the `idgen` tool (`R` = requirement prefix):

```
idgen -n <count> -p R
```

Mint as many as a Decision's Verification list needs, assign one id per behavior, and paste each inline. Ids are **stable handles**: when you edit the design in place, do **not** renumber or regenerate existing ids — mint a *fresh* id for each newly added behavior, and when a behavior is removed, delete its id with it (its test goes too). One id, one behavior, used in exactly one place.
