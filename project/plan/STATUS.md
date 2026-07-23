# idgen — Plan Status

This is the manifest: one line per **pending** phase in build order, and the
**only** place a phase's status marker lives. Each phase line is a Markdown bullet
beginning with `- Phase` and its zero-padded number, then `⬜` (pending), then
`realizes <Decision ids>` (or `realizes —` for a pure structural phase), then
`— <objective>`. The build loop finds its next work with
`grep -nE '^- Phase .* ⬜' project/plan/STATUS.md | head -1`, reads only that
phase's `project/plan/phase-NN.md`, and **on completion deletes that phase's line
here and its body file** — there is no done marker; done is gone. This file
deliberately carries **no** bare status glyph outside these phase lines, so the
anchored grep matches only phase lines.

Next phase: 05

- Phase 02b ⬜  realizes D2         — internal/idgen: decode (inverse direction) + randomized property sweeps
- Phase 03a ⬜  realizes D3, D4, D5, D6  — internal/cli: mint mode (Run dispatch, wait-loop, validation, version/usage)
- Phase 03b ⬜  realizes D4, D5  — internal/cli: decode mode
- Phase 04  ⬜  realizes D6         — wire cmd/idgen/main + install + release (goreleaser/workflow/install.sh) + version & build smokes
