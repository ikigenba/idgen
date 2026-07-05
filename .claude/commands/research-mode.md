---
description: Research mode — dispatch sequential subagents to investigate, then write the research doc
---
We are in **research mode**.

When this mode is entered, first load the `ralph` skill (the family map) — read `.claude/library/ralph/SKILL.md`, falling back to `~/.claude/library/ralph/SKILL.md` — then acknowledge it and **stop — wait for my instructions** before doing anything. Don't read docs or dispatch subagents until I tell you to proceed.

Research is **optional and non-contractual**: it exists to inform *you* before you author the design. The autonomous build reads only product, design, and plan — never this doc — so `project/research/research.md` feeds your thinking and nothing downstream consumes it. It does not feed or update any other doc; design stays the single authority for *how*, authored through design-mode.

This mode is **re-usable**: research can evolve as the product does. If `project/research/research.md` already exists, read it first and treat it as the current state. As the goal changes, **edit that doc in place to align with the new goal — never append**. The doc is always a single, coherent statement of the current research, not a running log.

Read `project/product/README.md` for context. Decompose the open questions, then **fan out** — dispatch one subagent per question in parallel, each with a complete cold-start brief, each returning a distilled summary. Keep the main thread lean; never pull raw files or search dumps into it. Cover the codebase, prior art, and external sources.

Then **synthesize** every return into `project/research/research.md`: options, prior art, constraints, recommendations. Report the path.
