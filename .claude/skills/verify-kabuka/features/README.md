# kabuka verification map

This directory is the maintained source for verifying the user-facing behavior of the kabuka CLI. Read this index before driving the binary, then use the matching feature file as the recipe.

## Baseline preconditions

- `H=.claude/skills/verify-kabuka/scripts/kabuka-verify.sh` from the repo root, running under bash.
- `RUN=$($H build)` was run against the checkout under test.
- `$H doctor "$RUN"` reports `healthy`. That includes a clean rev match and HTTP 200 from `finance.yahoo.co.jp`.
- Never verify with a `kabuka` found on `PATH`. Only the per-run binary counts.

## Driving conventions

- Every check is `$H run "$RUN" <evidence-name> -- <args>`. Treat args as literal, including `=` in `-f=json`.
- Use a unique evidence name per check. The feature files suggest names. Keep them, so runs are comparable across time.
- Prices are live. Assert shape (symbol, numeric price, field presence, exit code), never exact values.
- kabuka reads only the first positional arg. Extra symbols are silently ignored, so pass one symbol per run.

## Proof and skip reporting

- CLI proof is `.cmd` + `.stdout` + `.stderr` + `.exit` + `.sideeffects` (must be `none`) from the evidence dir.
- Record the feature ID with each evidence name (the suggested names already start with it).
- If the network is down or Yahoo serves `---`, report the check as blocked with the doctor output. Do not mark it verified.
- Do not report a skipped entry point (say, the `--format` long flag) as verified through a different one (`-f`).

## Feature entry contract

Each feature file starts with an H1 title and one paragraph describing the user-visible behavior. It then uses exactly four H2 sections in this order.

1. `Sub-features` lists short IDs with one line for each behavior.
2. `How to get to it (user POV)` lists every user entry point.
3. `Driving it with kabuka-verify.sh` starts with `Preconditions:` and uses labeled bullets that pair each user action with an exact command and observable result.
4. `Gotchas` lists traps that can waste or invalidate a verification run.

## Features

- [Quote a price](./quote-price.md) covers JP codes with and without `.T`, US tickers, and multi-market symbols.
- [Detailed quote](./detail.md) covers `-d` / `--detail` and its eight fields.
- [Output formats](./output-format.md) covers text, `json` and `csv` via `-f` / `--format`.
- [Error paths](./errors.md) covers unknown symbols, missing args, bad formats and bad flags.
