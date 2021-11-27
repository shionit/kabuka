package kabuka

import (
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/PuerkitoBio/goquery"
	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	financeSiteUrl = "https://info.finance.yahoo.co.jp/search/?query="
)

// Fetch stock information from finance website.
func (k *Kabuka) Fetch() (*model.Stock, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	// Yahoo! Finance Web scraping
	res, err := client.Get(financeSiteUrl + k.Symbol)
	if err != nil {
		return nil, xerrors.Errorf("Http client Get failed, err: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v\n", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("Http client Get status code error: %d %s, err: %w",
			res.StatusCode, res.Status, err)
	}
	if isSymbolNotFound(res) {
		return nil, xerrors.New("Symbol is not found.")
	}
	f := fetcher.SelectFetcher(res.Request.URL.String())
	if f == nil {
		return nil, xerrors.New("Unknown market type.")
	}
	paths := strings.Split(res.Request.URL.Path, "/")
	symbol := paths[len(paths)-1]

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("goquery NewDocument failed, err: %w", err)
	}
	return f.Fetch(doc, symbol)
}

func isSymbolNotFound(res *http.Response) bool {
	url := res.Request.URL.String()
	// If location is back to search page, result is "not found".
	return strings.HasPrefix(url, financeSiteUrl)
}
