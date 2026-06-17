---
description: Plan mode — decompose the design into sequential phases for build
---
We are in **plan mode**.

When this mode is entered, acknowledge it and **stop — wait for my instructions** before doing anything. Don't read docs or draft phases until I tell you to proceed.

Read `docs/design.md` and decompose it into an ordered sequence of phases that the build loop will execute one at a time. This is iterative, not a single pass — draft the phases, then refine them with me: resize, reorder, split, and merge until each is right. Don't assume the plan is finished in one shot. Order is dependency order: each phase depends only on earlier ones. The plan says *what*; build supplies *how*. When we're settled, write `docs/plan.md` in the shape below and report the path.

## Output shape

`docs/plan.md` owns **construction order and history** — the order the thing is built in and the record of what has been built. Unlike product and design (which are rewritten in place to stay authoritative for the current state), the plan is **append-only**: phases are added at the bottom and marked done as they land; completed phases are never rewritten or deleted, so the plan doubles as the construction history. To extend the project later, update product and design in place, then append a new phase here.

Write the doc in this shape:

- **Title** — `# <name> — Plan`.
- **Authority header** — a paragraph beginning `**Authority: construction order and history.**` stating that the plan owns the build order and the record of what's built, that it is **append-only** (completed phases are never rewritten or deleted, so it doubles as history), and how to extend it (update product and design in place, then append a new phase).
- **One phase = one package = one accumulating context** — a short paragraph: each phase is a single coherent unit — almost always one package — built in one accumulating context against product and design, reading only that package's design Decisions and the *interfaces* (not internals) of the packages it depends on. This is what keeps every phase the size of a small standalone tool no matter how large the project grows.
- **Done bar** — a line stating a phase is **done** when every Verification item in the design Decisions it realizes is covered by a clearly-named test and the suite is green (point at the design's *Verification & "done"* section).
- **## Status** — one line: overall progress (e.g. *not started; the workspace holds product, design, and this plan; no code yet*, or which phases are done).
- **## Phases** — the ordered list. Each phase is:
  - A header `### Phase N — <one cohesive objective> · ⬜ not started` — the `⬜ not started` / `✅ done` token is the **single status marker the build loop flips**; nothing else in a completed phase is ever edited.
  - A `*Realizes design Decision <n> (<short label>)[ and <m> (...)]. Depends on Phase <k>.*` line — name exactly which design Decisions this phase builds and which earlier phase(s) it needs. This link is what lets the build read only the relevant slice of design, not the whole doc.
  - A short body: what gets built (the package/seam and its paths), stated as the observable end state, not an implementation recipe.
  - **Done when:** the acceptance bar for this phase — its design Verification ids covered by genuine tests and the suite green.

Each phase must be completable in a single context — read, change, and verify — and leave the build green. If a phase won't fit one package's worth of context, it's too big; split it before writing.
