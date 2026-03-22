# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o bin/ -v ./...

# Run unit tests (no network, fast)
make test

# Run integration tests (requires live network access to Yahoo Finance)
make test-integration

# Run a single test
go test -v -run TestFunctionName ./internal/app/kabuka/...

# Run tests with coverage (full, includes integration)
make test-coverage
```

Linting runs via CI (golangci-lint, staticcheck via reviewdog). To run locally:
```bash
golangci-lint run
staticcheck ./...
```

## Architecture

**kabuka** is a CLI tool that fetches stock prices from Yahoo! Finance Japan (`finance.yahoo.co.jp`), supporting both Japanese (東証, 名証, etc.) and US (NASDAQ, NYSE) markets.

### Execution Flow

```
cmd/root.go (Cobra CLI) → Kabuka.Execute() → Kabuka.fetch() → SelectFetcher() → Fetcher.Fetch() → formatOutput()
```

1. **`cmd/root.go`** — Cobra command definition; parses symbol arg and `-f/--format` flag (text/json/csv), then calls `Kabuka.Execute()`
2. **`internal/app/kabuka/kabuka.go`** — Core orchestration: HTTP GET to Yahoo Finance search, handles search result pages by following targeted product links, delegates HTML parsing to the appropriate fetcher
3. **`internal/app/kabuka/fetcher/fetcher.go`** — Fetcher interface (`IsMarket(doc) bool`, `Fetch(doc, symbol) Stock`) plus registry (`RegisterFetcher` / `SelectFetcher`)
4. **`internal/app/kabuka/fetcher/jp/`** and **`fetcher/us/`** — Concrete fetcher implementations; registered via `init()`, detect their market from DOM, extract current price via CSS selectors

### Plugin Pattern for Fetchers

New market support is added by:
1. Implementing the `Fetcher` interface in a new package under `fetcher/`
2. Calling `fetcher.RegisterFetcher()` in the package's `init()`
3. Importing the package (blank import) in `kabuka.go`

### Key Design Notes

- Price extraction uses hardcoded CSS selectors against Yahoo Finance's DOM — changes to Yahoo Finance's HTML will break parsing
- `SanitizeInput()` strips newlines/carriage returns from user-supplied symbols
- `FormatPrice()` strips commas before returning numeric price strings

## Gotchas

- **CSS selector fragility**: All fetchers (`fetcher/jp/`, `fetcher/us/`) hardcode deep CSS selectors against Yahoo Finance's DOM. When fetching silently breaks, inspect `selectorCurrentPrice` and `selectorMarketNameSingle` against the current live page first — Yahoo Finance HTML changes will cascade to all fetchers at once.
- **Integration tests require live network**: `TestKabuka_fetch` (in `kabuka_integration_test.go`, build tag `integration`) makes real HTTP requests to Yahoo Finance. These tests fail without internet access and may return `"---"` during non-market hours — this is expected behaviour, not a bug. Use `make test` for offline development.
