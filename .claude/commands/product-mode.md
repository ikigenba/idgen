---
description: Product mode — interrogate to define the problem, then write the product doc
---
We are in **product mode**. Goal: pin down *what* we're building and *why*.

When this mode is entered, first load the `ralph` skill (the family map) — read `.claude/library/ralph/SKILL.md`, falling back to `~/.claude/library/ralph/SKILL.md` — then acknowledge it and **stop — wait for my instructions** before doing anything. Don't start interviewing until I tell you the feature or topic.

This mode is **re-usable**: the product can evolve. If `project/product/product.md` already exists, read it first and treat it as the current state. As the goal changes, **edit that doc in place to align with the new goal — never append**. The doc is always a single, coherent statement of the current product, not a changelog.

Interview me — problem, purpose, users, scope (and what's deliberately out), any contractual constants, the user-facing promises, and the success criteria — challenging assumptions and stating your recommendation with each. **Ask one question at a time; never batch.** Wait for my answer before the next one. Anything you can settle by reading the codebase, settle it yourself instead of asking. Keep going until the scope is sharp.

When we're done, write `project/product/product.md` in the shape below. Report the path.

## Output shape

`project/product/product.md` owns **intent** — *why* this exists, *for whom*, what is in and out of scope, and what we **promise** the user — stated once, in **outcome terms**. It must NOT state mechanism, exact formats, exit codes, or test assertions; those belong to `project/design/README.md`. Where the two could overlap (behavior), product states the *promise*; design states the *exact, checkable proof of that promise*. This boundary is load-bearing — it is what keeps product, design, and plan from overlapping.

Write these sections, in order:

- **Title** — `# <name> — Product`.
- **Authority header** — a short paragraph beginning `**Authority: intent.**` stating what this doc owns and, explicitly, what it does not (mechanism, formats, exit codes, test assertions → design), plus the promise-vs-proof boundary above.
- **## Problem** — the pain in the user's world; no solution yet.
- **## Purpose** — one paragraph: what the thing *is* and the single job it does.
- **## Users** — who runs it and what they are trying to get done.
- **## Scope** — what it does and, by exclusion, what it deliberately does not. Fold non-goals in here as bounded "nothing else" statements; only break out a separate `## Non-goals` section when the exclusions genuinely need their own emphasis.
- **## Contractual constants** *(only if any exist)* — fixed, promised values the design must use verbatim and never re-declare, such as a baseline constant, a starting version, or a protocol value. These are promises, not implementation detail. Omit the whole section when the product has none.
- **## What we promise (user-facing behavior)** — the observable behavior in outcome terms: what the user does and what they get back. Use short example invocations/output where they sharpen the promise. No mechanism, no exit codes, no internal formats.
- **## Success criteria (outcomes)** — a bullet list of user-observable outcomes, each phrased as a *result* the user could confirm, never as a test assertion or mechanism. The verification gate runs the built artifact against exactly this list, so every item must be outcome-shaped and checkable end-to-end against the real thing.
