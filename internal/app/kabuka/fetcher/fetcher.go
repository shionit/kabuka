package fetcher

import (
	"strings"

	"github.com/shionit/kabuka/internal/app/kabuka/model"

	"github.com/PuerkitoBio/goquery"
)

const (
	selectorMarketNameSingle = "span._3sg2Atie"
	selectorMarketNameMulti  = "div._34dIY8Xd._2i88cY3H > button"
)

var (
	fetchers []Fetcher
)

type Fetcher interface {
	// IsMarket returns whether the fetcher supports the document
	IsMarket(doc *goquery.Document) bool
	// Fetch parses the document and returns a parsed Stock model
	Fetch(doc *goquery.Document, symbol string) (*model.Stock, error)
}

// RegisterFetcher registers a fetcher
func RegisterFetcher(f Fetcher) {
	fetchers = append(fetchers, f)
}

// SelectFetcher returns fetcher that can parse the document
func SelectFetcher(doc *goquery.Document) Fetcher {
	for _, f := range fetchers {
		if f.IsMarket(doc) {
			return f
		}
	}
	return nil
}

// FormatPrice formats price
func FormatPrice(s string) string {
	return strings.ReplaceAll(s, ",", "")
}

// GetMarketName returns the document's market name
func GetMarketName(doc *goquery.Document) string {
	selection := doc.Find(selectorMarketNameSingle)
	if strings.TrimSpace(selection.Text()) == "" {
		selection = doc.Find(selectorMarketNameMulti)
	}
	return selection.Text()
}
