# Quote a price

A user passes one stock symbol and kabuka prints the resolved symbol and the current price, tab-separated. Japanese codes resolve to their `.T` form. US tickers are printed as given.

## Sub-features

- `quote-jp-suffixed` quotes a Tokyo symbol given with `.T`.
- `quote-jp-bare` resolves a bare 4-digit code to `<code>.T` through Yahoo's search page.
- `quote-us` quotes a NASDAQ/NYSE ticker.
- `quote-multi-market` picks the right listing for a company listed in both Tokyo and New York.

## How to get to it (user POV)

- Run `kabuka <symbol>` in a terminal.

## Driving it with kabuka-verify.sh

Preconditions:

- Baseline from `README.md`: a built run and a healthy doctor.

- **JP with suffix.** Run `$H run "$RUN" quote-jp-suffixed -- 3994.T`. Exit `0`. stdout is one line, `3994.T\t<price>`, and the price matches `^[0-9]+(\.[0-9]+)?$`.
- **JP bare code.** Run `$H run "$RUN" quote-jp-bare -- 4373`. Exit `0`. stdout is `4373.T\t<price>`. The `.T` was added by resolution, not by the user.
- **US ticker.** Run `$H run "$RUN" quote-us -- AAPL`. Exit `0`. stdout is `AAPL\t<price>`, and the price has a decimal part (for example `343.46`).
- **Multi-market, Tokyo side.** Run `$H run "$RUN" quote-multi-market-jp -- 8604.T`. Exit `0`. stdout starts with `8604.T\t`, and the price is in yen (hundreds to thousands).
- **Multi-market, NYSE side.** Run `$H run "$RUN" quote-multi-market-us -- NMR`. Exit `0`. stdout starts with `NMR\t`, and the price is in dollars (single digits). It must differ in magnitude from the Tokyo price.
- **Proof.** Every `.sideeffects` file reads `none`. Report the five evidence names and the evidence dir.

## Gotchas

- Outside market hours Yahoo may serve `---` as the price. That is expected. Report it as blocked rather than failed.
- `Unknown market type.` on symbols that worked before means Yahoo changed the market-label DOM. That is not a symbol problem. Check `selectorMarketNameSingle` per `CLAUDE.md`.
- Supported JP markets are matched by prefix: 東証PRM/GRT/STD, 名証MN, 札証, 札幌ア, 福証, 福岡Q. The seeds above only cover 東証. A non-Tokyo-only symbol is unproven until someone finds one and adds it here.
- `kabuka 3994.T 4373` prints only `3994.T`. Pass one symbol per check.
