# Output formats

`-f` / `--format` selects how the quote is printed: `text` (default, tab-separated), `json` (one object), or `csv` (header plus one row).

## Sub-features

- `format-text-default` is tab-separated text when no flag is given or with `-f=text`.
- `format-json` prints one JSON object with `symbol` and `current_price`.
- `format-csv` prints a header line and a data row.
- `format-spellings` accepts `-f=json`, `-f json` and `--format=json` alike.

## How to get to it (user POV)

- Run `kabuka <symbol> -f=<fmt>`, `kabuka <symbol> -f <fmt>` or `kabuka <symbol> --format=<fmt>`.

## Driving it with kabuka-verify.sh

Preconditions:

- Baseline from `README.md`: a built run and a healthy doctor.

- **Explicit text.** Run `$H run "$RUN" format-text -- 3994.T -f=text`. Exit `0`. stdout is `<price>\t3994.T`, identical in shape to running without `-f`.
- **JSON.** Run `$H run "$RUN" format-json -- 3994.T -f=json`. Exit `0`. stdout parses as JSON (`jq -e . <file>`), equals `{"symbol":"3994.T","current_price":"<price>"}`, and has no detail keys.
- **JSON, space-separated flag.** Run `$H run "$RUN" format-json-space -- 8604.T -f json`. Exit `0`. stdout is a JSON object with `"symbol":"8604.T"`.
- **JSON, long flag.** Run `$H run "$RUN" format-json-long -- NMR --format=json`. Exit `0`. stdout is a JSON object with `"symbol":"NMR"`.
- **CSV.** Run `$H run "$RUN" format-csv -- 3994.T -f=csv`. Exit `0`. Line 1 is `symbol,current_price,change,change_pct,open,high,low,volume`. Line 2 starts with `3994.T,<price>,` followed by six empty fields. A trailing blank line follows.
- **CSV with detail.** Run `$H run "$RUN" format-csv-detail -- -d 3994.T -f=csv`. Exit `0`. Same header, and line 2 has all 8 fields non-empty.
- **Proof.** Every `.sideeffects` file reads `none`. Report the JSON parse result alongside the evidence path.

## Gotchas

- The README shows non-detail CSV as `symbol,current_price` only. The binary always emits all 8 columns, leaving the detail ones empty. Trust the binary.
- CSV output ends with an extra empty line (gocsv newline plus `fmt.Println`). Don't treat it as a malformed row.
- Format values are case-sensitive. `-f=JSON` is rejected (see [errors](./errors.md)).
