---
name: add-fetcher
description: Add a new stock market fetcher to the kabuka project
disable-model-invocation: true
---
Add a new market fetcher for: $ARGUMENTS

Steps:

1. Identify the exact market name string Yahoo Finance uses in the DOM — inspect the element targeted by `selectorMarketNameSingle` in `fetcher/fetcher.go` on a real Yahoo Finance page for that market.

2. Create `internal/app/kabuka/fetcher/<market>/` with a `<market>_fetcher.go` file. Use `fetcher/jp/jp_fetcher.go` as the reference implementation.

3. Implement the `Fetcher` interface:
   - `IsMarket(doc *goquery.Document) bool` — use `fetcher.GetMarketName(doc)` and match against the market name(s) identified in step 1
   - `Fetch(doc *goquery.Document, symbol string) (*model.Stock, error)` — extract price using `doc.Find(selectorCurrentPrice).Text()` then wrap with `fetcher.FormatPrice()`

4. Register via `init()`:
   ```go
   func init() {
       fetcher.RegisterFetcher(&<market>Fetcher{})
   }
   ```

5. Add a blank import in `internal/app/kabuka/kabuka.go`:
   ```go
   _ "github.com/shionit/kabuka/internal/app/kabuka/fetcher/<market>"
   ```

6. Add integration test cases to `kabuka_integration_test.go` covering at least one symbol from the new market.

7. Verify:
   ```bash
   go build ./...
   go test -tags integration -run TestKabuka_fetch ./internal/app/kabuka/...
   ```
