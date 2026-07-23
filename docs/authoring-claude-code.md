# Author a spec with Claude Code

Building `idgen` from its spec is one half of this repo; the other half is the
method that wrote the spec, and it ships in the repo too. This page walks that
method twice with Claude Code: once to **create a new project from an empty
directory**, and once to **add a feature to an existing spec**. The two walks
are nearly identical on purpose; the second is the first minus the
bootstrapping.

The example app is deliberately sketchy — a stand-in you'd replace with your
own idea. What matters here is the beats, not the app.

## Requirements

- [Claude Code](https://claude.com/claude-code), installed and authenticated.
- The **`idgen` binary on your `PATH`**: sealing the spec mints requirement ids
  with it. Build it from this repo first (any method on the
  [README](../README.md)), then `make install`.
- The [`ralph`](https://github.com/ikigenba/ralph) harness for the unattended
  run — or drive the loop interactively instead, as
  [claude-code.md](claude-code.md) shows.

## Part 1 — a new project from scratch

### 1. Make a repo and carry the method in

The spec-process skills travel by copying the two agent trees out of this repo:

```sh
mkdir grump && cd grump && git init
cp -r ../idgen/.claude .
cp -r ../idgen/.agents .    # optional: the same skills packaged for Codex
```

### 2. Open the spec session

Start `claude` in the new repo, then open a spec session:

```
/open-spec
```

It loads the spec contracts (the `ikispec` shapes) and scopes the session to
`project/*`.

### 3. Describe the goal in plain English

Just talk. A paragraph is plenty to start:

> I want a small CLI called grump. It reads text on stdin and prints the same
> text back with complaints attached — one grumble per line, deterministic so
> tests can pin it. Single static Go binary, standard library only.

### 4. Get grilled

Sharpen the goal one question at a time — invoke `/grill-me`:

```
/grill-me
```

The agent interrogates one question at a time, each with a recommendation
attached, until the goal is settled:

> **Q:** Should grump preserve the input ordering, or sort lines before
> grumbling? I recommend preserving order — it keeps the tool composable in
> pipes. — *preserve it.*
>
> **Q:** What happens on empty stdin? I recommend exiting 0 with no output,
> matching how `cat` behaves. — *agreed.*

Answer until it runs out of questions, not until it "seems good enough".

### 5. Seal the spec

```
/seal-spec
```

One automated pass, no interviewing. When it finishes, `project/` exists:
`product/README.md` (the why, as outcomes), `design/` (one Decision per file,
every checkable behavior carrying a **minted `R-XXXX-XXXX` id** — this is where
your `idgen` binary gets used), `plan/` (a queue of dependency-ordered pending
phases, all `⬜` in `STATUS.md`, with a `Next phase` counter), and the workspace
map.

### 6. Generate the build loop (once per project)

```
/create-gather-build-verify-prompts
```

This writes `project/loops/{gather,build,verify}.md` plus the loop's
`README.md`. Optionally also run `/create-audit-prompts` for the adversarial
coverage-audit prompt.

### 7. First run

From a shell, launch `ralph` on the three prompts. By convention, put the
invocation behind a committed wrapper so every future run is one short command.
This is the minimal default; add your chosen harness/model flags here if needed:

```sh
cat > project/loops/run <<'EOF'
#!/bin/bash
exec ralph \
  project/loops/gather.md \
  project/loops/build.md \
  project/loops/verify.md
EOF
chmod +x project/loops/run
project/loops/run
```

Watch phases disappear from `project/plan/STATUS.md` one by one — verify deletes
each finished phase's line (and its `phase-NN.md`) as it goes. The run ends when
gather finds no pending phase and reports `DONE` — and the source that didn't
exist now does, with every requirement id covered by an id-tagged test.

## Part 2 — adding a feature

Same beats, existing spec — this works identically on the project above or on
`idgen` itself:

1. Start `claude` at the repo root and open the session: `/open-spec`.
2. Describe the change in plain English: *"grump should also take a
   `--mood MOOD` flag that picks the complaint register."*
3. `/grill-me`, and answer the questions until settled.
4. `/seal-spec`. The differences from greenfield are worth noticing: product and
   design are **rewritten in place** to the new current truth, **fresh ids are
   minted only for the new behaviors**, and the plan is **appended** — a new
   `phase-NN.md` and a new `⬜` line in `STATUS.md`, numbered from the `Next phase`
   counter, with existing phases left untouched.
5. `project/loops/run`. gather finds only the new `⬜` phase (completed phases are
   already gone), so the loop builds just the delta and reports `DONE`.

There is no step 6: the loop is already installed. Regenerate it only when the
loop design itself changes.

Optionally close either walk with the audit loop —
`ralph project/loops/audit.md` — to adversarially re-check that every id's test
could actually fail.
