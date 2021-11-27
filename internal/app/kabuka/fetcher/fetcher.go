package fetcher

import (
	"strings"

	"github.com/shionit/kabuka/internal/app/kabuka/model"

	"github.com/PuerkitoBio/goquery"
)

var (
	fetchers []Fetcher
)

type Fetcher interface {
	IsMarketUrl(url string) bool
	Fetch(doc *goquery.Document, symbol string) (*model.Stock, error)
}

// RegisterFetcher registers a fetcher
func RegisterFetcher(f Fetcher) {
	fetchers = append(fetchers, f)
}

// SelectFetcher returns fetcher that matches the url
func SelectFetcher(url string) Fetcher {
	for _, f := range fetchers {
		if f.IsMarketUrl(url) {
			return f
		}
	}
	return nil
}

// FormatPrice formats price
func FormatPrice(s string) string {
	return strings.ReplaceAll(s, ",", "")
}
