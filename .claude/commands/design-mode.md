---
description: Design mode — propose a testable architecture and debate tradeoffs, then write the design doc
---
We are in **design mode**.

When this mode is entered, acknowledge it and **stop — wait for my instructions** before doing anything. Don't read docs or propose architecture until I tell you to proceed.

**Purpose.** Design mode defines the *shape* of the application: its seams (the boundaries between components and where they can be substituted), its public interfaces, its naming, and its struct/type definitions. It also defines **how the software is tested** — the testing strategy is part of the architecture, not an afterthought. **It is the responsibility of design mode to produce a *testable* architecture**: seams exist so behavior can be exercised in isolation, and every component is shaped so its correctness can be verified.

**Design mode does not write the application code.** Its output is the design doc — interfaces, types, seams, naming, and the test plan. Illustrative signatures, struct definitions, and interface declarations belong in the doc; full implementations do not. Writing the code is the job of a later build phase.

This mode is **re-usable**: the design can evolve as the product does. The doc is always a single, coherent statement of the current design, not a changelog. When the goal shifts, **edit the affected parts of the doc in place to realign — never append** contradictory sections.

Read `docs/product.md`, `docs/research.md` if it exists, and any existing `docs/design.md`. Decisions already recorded in the design doc are settled — don't reopen them; pick up from what's still undecided. Propose the architecture — seams, public interfaces, naming, struct/type definitions, data model, key decisions. **Every verifiable design element must describe how it is tested as part of the design itself** — there is no design element that defines behavior without also stating how that behavior is verified. Debate the tradeoffs with me, surfacing alternatives and your recommendation. **Raise one decision at a time; never batch.** Wait for my call before the next one. Decisions need my buy-in.

Record each decision into `docs/design.md` as it's settled — the doc is the running record, so we can stop and resume across sessions. We're done when the seams, public interfaces, naming, struct/type definitions, data model, and the testing approach are fully decided.

## Output shape

`docs/design.md` owns **shape and its proof** — *how* the thing is built and *how each behavior is proven*. The product owns the *why* and the user-facing promises; design states the **exact, checkable form** of those promises and never re-declares the why. It is the **single, current** statement of the architecture: when a decision changes, the doc is rewritten to stay true — stale decisions are removed, not stacked. History of how it got here lives in the plan.

Write the doc in this shape:

- **Title** — `# <name> — Design`.
- **Authority header** — a short paragraph beginning `**Authority: shape and its proof.**` stating what design owns (how + proof), that product owns the why and the promises, and that design *uses* the product's contractual constants by value but does not own them. Add a line noting this is the single current statement, rewritten in place (not stacked), with history living in the plan.
- **## Verification & "done" — the denominator** — the contract section. State plainly that:
  - Each Decision ends with a **Verification** list: the concrete behaviors a test must assert for that decision to be built.
  - Every Verification item carries a **minted id** of the form `R-XXXX-XXXX` — a stable, unique handle for that one behavior.
  - **That set of lists is the denominator** — the enumerated intent the test suite is measured against. There is **no separate requirements document**; the ids live inline here and nowhere else.
  - A behavior is **covered** when a test asserts it *and names its id in a `// R-XXXX-XXXX` comment*, so coverage is a grep — counted only inside a `//` comment, never a bare string literal (that keeps ids that appear as test *data* from masquerading as coverage).
  - The work is **done** when every Verification id is covered by exactly one genuine test and the suite is green.
- **## Conventions** *(optional)* — shared facts every Decision leans on (language/version, module path, exit-code taxonomy, formatting rules, a shared time/IO source). Omit if there's nothing cross-cutting to state.
- **## Decision N — <title>** — one per decision, in the order settled. Each contains:
  - **Decision.** — the seams/interfaces/types/naming, with illustrative signatures and struct/interface declarations (never full implementations).
  - **Rejected.** — the alternatives considered and why each lost.
  - **Verification.** — a bullet list, each line `R-XXXX-XXXX — <the behavior a test must assert>`. A pure seam/structure decision with no behavior of its own says so explicitly and carries no ids (its proof is the behavioral ids of the decisions it enables).
- **## Status** — one short section: what is decided, and that the construction order realizing this design lives in the plan.

## Minting the Verification ids

The `R-XXXX-XXXX` ids are **real, minted ids — never hand-written or made up.** Mint them with the `idgen` tool (`R` = requirement prefix):

```
idgen -n <count> -p R
```

Mint as many as a Decision's Verification list needs, assign one id per behavior, and paste each inline. Ids are **stable handles**: when you edit the design in place, do **not** renumber or regenerate existing ids — mint a *fresh* id for each newly added behavior, and when a behavior is removed, delete its id with it (its test goes too). One id, one behavior, used in exactly one place.
