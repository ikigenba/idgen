# idgen — Design Index

Each Decision maps to its `DNN.md`; every `R-XXXX-XXXX` id maps to its
Decision/file. Resolve an id by grepping this index (or the Decision files
directly). Regenerate this file whenever a Decision is added or its Verification
ids change.

## Decisions

- **D1** → `project/design/D01.md` — Module path & package layout (top-level seams) — ids: (none — seam choice)
- **D2** → `project/design/D02.md` — `idgen` public API & prefix placement — ids: R-WH5F-QJYS, R-WIDC-4BPH, R-WJL8-I3G6, R-WKT4-VV6V, R-WM11-9MXK, R-WN8X-NEO9, R-WPOQ-EY5N, R-WQWM-SPWC, R-WS4J-6HN1
- **D3** → `project/design/D03.md` — `Clock` seam & the `-n` wait loop — ids: R-WTCF-K9DQ, R-WUKB-Y14F, R-WVS8-BSV4, R-WX04-PKLT, R-WY81-3CCI, R-WZFX-H437
- **D4** → `project/design/D04.md` — CLI grammar, dispatch & exit-code taxonomy — ids: R-X0NT-UVTW, R-X1VQ-8NKL, R-TIQT-QVU3, R-X33M-MFBA, R-X4BJ-071Z, R-X5JF-DYSO, R-X6RB-RQJD, R-7UL7-PF0O, R-PU67-68HE
- **D5** → `project/design/D05.md` — Input handling & validation (both modes) — ids: R-X974-JA0R, R-XAF0-X1RG, R-XBMX-ATI5, R-XCUT-OL8U, R-XE2Q-2CZJ, R-XFAM-G4Q8, R-XGII-TWGX, R-XHQF-7O7M, R-XIYB-LFYB
- **D6** → `project/design/D06.md` — Version, usage text, Makefile & release — ids: R-TJYQ-4NKS, R-TL6M-IFBH, R-XLE4-CZFP, R-XMM0-QR6E
- **D7** → `project/design/D07.md` — Overall testing strategy & test layout — ids: (none — strategy only)

## Verification ids → Decision

- R-7UL7-PF0O → D4 (`project/design/D04.md`)
- R-PU67-68HE → D4 (`project/design/D04.md`)
- R-TIQT-QVU3 → D4 (`project/design/D04.md`)
- R-TJYQ-4NKS → D6 (`project/design/D06.md`)
- R-TL6M-IFBH → D6 (`project/design/D06.md`)
- R-WH5F-QJYS → D2 (`project/design/D02.md`)
- R-WIDC-4BPH → D2 (`project/design/D02.md`)
- R-WJL8-I3G6 → D2 (`project/design/D02.md`)
- R-WKT4-VV6V → D2 (`project/design/D02.md`)
- R-WM11-9MXK → D2 (`project/design/D02.md`)
- R-WN8X-NEO9 → D2 (`project/design/D02.md`)
- R-WPOQ-EY5N → D2 (`project/design/D02.md`)
- R-WQWM-SPWC → D2 (`project/design/D02.md`)
- R-WS4J-6HN1 → D2 (`project/design/D02.md`)
- R-WTCF-K9DQ → D3 (`project/design/D03.md`)
- R-WUKB-Y14F → D3 (`project/design/D03.md`)
- R-WVS8-BSV4 → D3 (`project/design/D03.md`)
- R-WX04-PKLT → D3 (`project/design/D03.md`)
- R-WY81-3CCI → D3 (`project/design/D03.md`)
- R-WZFX-H437 → D3 (`project/design/D03.md`)
- R-X0NT-UVTW → D4 (`project/design/D04.md`)
- R-X1VQ-8NKL → D4 (`project/design/D04.md`)
- R-X33M-MFBA → D4 (`project/design/D04.md`)
- R-X4BJ-071Z → D4 (`project/design/D04.md`)
- R-X5JF-DYSO → D4 (`project/design/D04.md`)
- R-X6RB-RQJD → D4 (`project/design/D04.md`)
- R-X974-JA0R → D5 (`project/design/D05.md`)
- R-XAF0-X1RG → D5 (`project/design/D05.md`)
- R-XBMX-ATI5 → D5 (`project/design/D05.md`)
- R-XCUT-OL8U → D5 (`project/design/D05.md`)
- R-XE2Q-2CZJ → D5 (`project/design/D05.md`)
- R-XFAM-G4Q8 → D5 (`project/design/D05.md`)
- R-XGII-TWGX → D5 (`project/design/D05.md`)
- R-XHQF-7O7M → D5 (`project/design/D05.md`)
- R-XIYB-LFYB → D5 (`project/design/D05.md`)
- R-XLE4-CZFP → D6 (`project/design/D06.md`)
- R-XMM0-QR6E → D6 (`project/design/D06.md`)
