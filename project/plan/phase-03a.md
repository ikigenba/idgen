# Phase 03a — `internal/cli`: mint mode

*Realizes design Decision 3 (`Clock` seam & wait loop) in full, the mint-side
slice of Decision 5 (input handling & validation), the mint-only slice of
Decision 4 (grammar, dispatch & exit codes), and the version/usage slice of
Decision 6; the testing strategy is Decision 7. Depends on Phase 02b.*

The mint half of the CLI, consuming `internal/idgen` only through `MintAt`:
`Run(...)` with one `flag.FlagSet`, enough dispatch to route a non-`--decode`
invocation to the mint path; the `Clock` seam and distinct-millisecond mint wait
loop (D3); default/`-n`/`-p` mint with prefix and number validation (D5 mint
slice); `--help`/`--version`/usage text (D6 version/usage slice); the `0`/`2`
exit-code taxonomy for these paths (D4 mint slice). The observable end state:
`cli.Run(args, stdin, stdout, stderr, clk)` mints correctly — default prefix,
custom prefix, `-n` repetition, validation errors, help/version — over in-memory
streams with an injectable clock. (`--decode` dispatch itself, and everything
downstream of it, is Phase 03b.)

**Done when** `go test -race ./...` exits 0 and each design Verification id below
is covered by a clearly-named, genuinely-asserting, id-tagged test (a `//` comment
naming the id) in `internal/cli/*_test.go` — table-driven over `Run(...)` with
in-memory buffers and a **fake `Clock`** (no subprocess, no real sleeps):

*Clock seam & wait loop (D3):*
- R-WTCF-K9DQ — fake clock (`Sleep` advances virtual `now`): `-n N` prints N distinct ids.
- R-WUKB-Y14F — under that clock, virtual time advanced ≥ N−1 ms (the per-id wait happened).
- R-WVS8-BSV4 — stalled clock still terminates; the N ids are pairwise distinct with the last ms ≥ N−1 beyond the first (not just success + one `Sleep`).
- R-WX04-PKLT — default/`N=1`: zero `Sleep` calls.
- R-WY81-3CCI — mint from an already-elapsed instant: the id decodes to the just-read `now`.
- R-WZFX-H437 — backward-clock tolerate: clock steps backward mid-sequence then recovers via `Sleep`; minted ms are non-decreasing **and** advance past the pre-step value (the dip is actually traversed).

*Grammar, dispatch & exit codes — mint slice (D4):*
- R-X0NT-UVTW — `--help`/`-h` → `Usage:` line counted **exactly once** on stdout, stderr empty, exit 0 (rules out the stdlib-`flag` double print).
- R-X1VQ-8NKL — `--version` → version string on stdout, exit 0.
- R-X33M-MFBA — unknown flag → exit 2, stderr non-empty (flag's own parse error).
- R-X4BJ-071Z — mint with a positional argument → exit 2 **and** stderr is
  non-empty (names the unexpected argument) — not merely a non-zero exit.
- R-7UL7-PF0O — bare invocation (`idgen`, no flags, no positionals): stdout is
  exactly one line matching `^R-[0-9A-Z]{4}-[0-9A-Z]{4}$` — the default action
  mints once with the default `R` prefix, never erroring and never requiring
  `-p`.
- R-PU67-68HE — mint with `-p X` (1-char and multi-char custom prefix): stdout
  is exactly one line matching `^X-[0-9A-Z]{4}-[0-9A-Z]{4}$` — the supplied
  prefix replaces the default outright, so `R` never appears in the output.

*Input handling & validation — mint slice (D5):*
- R-XGII-TWGX — prefix validation: `""`, `"  "`, `"R-X"`, `"S/"` → exit 2 and
  stderr contains `invalid prefix` (never a silent exit).
- R-XHQF-7O7M — number validation: `0`, `-3` → exit 2 and stderr contains
  `--number must be > 0` (never a silent exit).

*Version & usage (D6 slice):*
- R-XK67-Z7P0 — `--version` stdout is exactly `v0.1.0`.
- R-XLE4-CZFP — `--help` usage mentions each flag (`-n`, `-p`, `--decode`).

*(`--decode` dispatch, the decode-only D4 ids, the decode slice of D5, and the D6
build-smoke requirement are carried by Phase 03b / Phase 04.)*
