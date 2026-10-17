# Detailed quote

With `-d` / `--detail`, kabuka prints eight fields: symbol, price, change, change%, open, high, low and volume. Without it, only price and symbol are shown.

## Sub-features

- `detail-text` prints eight tab-separated fields, symbol first.
- `detail-long-flag` makes `--detail` behave the same as `-d`.
- `detail-json` populates the optional JSON keys that plain JSON omits.

## How to get to it (user POV)

- Run `kabuka -d <symbol>` or `kabuka --detail <symbol>`. The flag may come before or after the symbol.

## Driving it with kabuka-verify.sh

Preconditions:

- Baseline from `README.md`: a built run and a healthy doctor.

- **Short flag, JP.** Run `$H run "$RUN" detail-text-jp -- -d 3994.T`. Exit `0`. stdout is one line with exactly 8 tab-separated fields: `3994.T`, a numeric price, a change starting with `+`/`-` (or `0`), a percent ending in `%`, then numeric open/high/low/volume.
- **Long flag, US.** Run `$H run "$RUN" detail-long-flag-us -- --detail AAPL`. Exit `0`. Same 8-field shape with `AAPL` first.
- **Flag after symbol.** Run `$H run "$RUN" detail-flag-after -- 3994.T -d`. Exit `0`. Same 8-field shape as the short-flag check.
- **Detail JSON.** Run `$H run "$RUN" detail-json -- -d AAPL -f=json`. Exit `0`. stdout is one JSON object with keys `symbol`, `current_price`, `change`, `change_pct`, `open`, `high`, `low` and `volume`, all non-empty strings.
- **Contrast without detail.** Run `$H run "$RUN" detail-off -- 3994.T`. stdout has exactly 2 fields. This proves the flag, not the symbol, is what adds the fields.
- **Proof.** Pair `detail-text-jp` with `detail-off`, plus `.sideeffects` reading `none` for each.

## Gotchas

- In text mode, `-d` reorders output to symbol-first. Non-detail output is price-first. Don't reuse the non-detail parser.
- Outside trading hours, open/high/low can be identical and volume tiny (seen: `6427.1 6427.1 6427.1 100`). That is real data, not a parsing bug.
