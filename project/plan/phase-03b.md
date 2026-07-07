# Phase 03b — `internal/cli`: decode mode

*Realizes the decode-only slice of design Decision 4 (grammar, dispatch & exit
codes) and the decode slice of Decision 5 (input handling & validation); the
testing strategy is Decision 7. Depends on Phase 03a.*

The decode half of the CLI, added onto the `Run` dispatch skeleton Phase 03a
already built: the `--decode` branch (and `-n`/`-p` becoming inert under it);
`runDecode` with args-then-stdin precedence, partial-failure → exit 1, empty →
exit 0, UTC output regardless of `TZ`. Consumes `internal/idgen` only through
`TimeOf`. The observable end state: `cli.Run(args, stdin, stdout, stderr, clk)`
decodes correctly — positional args, stdin, mixed batches, empty input — and a
freshly minted id (Phase 03a) round-trips through `--decode` to its minting
instant.

**Done when** `go test -race ./...` exits 0 and each design Verification id below
is covered by a clearly-named, genuinely-asserting, id-tagged test (a `//` comment
naming the id) in `internal/cli/*_test.go` — table-driven over `Run(...)` with
in-memory buffers:

*Grammar, dispatch & exit codes — decode slice (D4):*
- R-X5JF-DYSO — `--decode` routes to the decode path.
- R-X6RB-RQJD — `-n`/`-p` with `--decode` are inert.

*Input handling & validation — decode slice (D5):*
- R-X974-JA0R — decode from positional args: one UTC line per id, in order.
- R-XAF0-X1RG — decode from stdin (mixed whitespace) matches the args case.
- R-XBMX-ATI5 — positionals win: with a stdin reader that fails if read, decode uses only positionals and stdin is never read (a `Read` call fails the test).
- R-XCUT-OL8U — one malformed token: good ids still decode, error names the bad token, exit 1.
- R-XE2Q-2CZJ — empty decode (no args, empty stdin): no output, exit 0.
- R-XFAM-G4Q8 — round-trip through `Run`: `--decode` of a freshly minted id returns its minting instant.
- R-XIYB-LFYB — decode output is UTC regardless of `TZ` (test sets `TZ=America/Chicago`).

*(The D6 build-smoke requirement is carried by Phase 04.)*
