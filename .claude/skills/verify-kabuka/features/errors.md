# Error paths

When kabuka cannot produce a quote, it prints a message to stderr, prints nothing to stdout, and exits `1`. Usage errors come from cobra and include the usage text. Runtime errors are a single timestamped log line.

## Sub-features

- `error-notfound` reports a symbol Yahoo cannot resolve as `Symbol is not found.`.
- `error-noarg` makes a missing symbol print a usage error.
- `error-badformat` rejects an unsupported `-f` value with `Unsupported format: <value>`.
- `error-badflag` makes an unknown flag print a usage error.

## How to get to it (user POV)

- Run `kabuka` with a bad symbol, no symbol, a bad `-f` value or an unknown flag.

## Driving it with kabuka-verify.sh

Preconditions:

- Baseline from `README.md`: a built run and a healthy doctor. The not-found check needs the network, because "not found" is Yahoo's answer.

- **Unknown symbol.** Run `$H run "$RUN" error-notfound -- ZZZZNOTREAL`. Exit `1`. stdout is empty. stderr is one line, `YYYY/MM/DD HH:MM:SS Symbol is not found.`
- **No symbol.** Run `$H run "$RUN" error-noarg --`. Exit `1`. stdout is empty. stderr starts with `Error: requires at least 1 arg(s), only received 0` and contains `Usage:`.
- **Bad format.** Run `$H run "$RUN" error-badformat -- 3994.T -f=xml`. Exit `1`. stdout is empty. stderr ends with `Unsupported format: xml`.
- **Bad format case.** Run `$H run "$RUN" error-badformat-case -- 3994.T -f=JSON`. Exit `1`. stderr ends with `Unsupported format: JSON`.
- **Unknown flag.** Run `$H run "$RUN" error-badflag -- 3994.T --version`. Exit `1`. stderr starts with `Error: unknown flag: --version`.
- **Proof.** Each check shows empty `.stdout`, `.exit` = `1`, and `.sideeffects` = `none`.

## Gotchas

- The format is validated before any network call, so `error-badformat` passes offline. Don't use it as a network-health signal.
- An HTTP failure surfaces as `Http client Get status code error: 400 ...` on stderr. Seen when the whole argv was passed as one string from zsh. If you see a 400 on a normal symbol, check your quoting before suspecting Yahoo.
- `Unknown market type.` is not a user error. It means the fetcher selectors no longer match Yahoo's DOM. Treat it as a regression in [quote-price](./quote-price.md).
- `--version` does not exist. Don't use it as a liveness probe. `doctor` uses `--help`.
