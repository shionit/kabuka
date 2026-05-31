package fetcher

import (
	"strings"

	"github.com/shionit/kabuka/internal/app/kabuka/model"

	"github.com/PuerkitoBio/goquery"
)

const (
	selectorMarketNameSingle    = "[class*='PriceBoardMenu__label']"
	selectorMarketNameMulti     = "[class*='PriceBoardMenu__toggle']"
	selectorChangeAbs           = "[class*='PriceChangeLabel__primary'] [class*='StyledNumber__value']"
	selectorChangePct           = "[class*='PriceChangeLabel__secondary'] [class*='StyledNumber__value']"
	selectorDataListItemTerm    = "[class*='DataListItem__term']"
	selectorDataListItemName    = "[class*='DataListItem__name']"
	selectorDataListItemValue   = "[class*='DataListItem__value']"
	labelOpen                   = "始値"
	labelHigh                   = "高値"
	labelLow                    = "安値"
	labelVolume                 = "出来高"
)

var (
	fetchers []Fetcher
)

type Fetcher interface {
	// IsMarket returns whether the fetcher supports the document
	IsMarket(doc *goquery.Document) bool
	// Fetch parses the document and returns a parsed Stock model
	Fetch(doc *goquery.Document, symbol string, detail bool) (*model.Stock, error)
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

// FetchDetailFields fills detail fields on stock using selectors common to all markets.
func FetchDetailFields(doc *goquery.Document, stock *model.Stock) {
	stock.Change = doc.Find(selectorChangeAbs).First().Text()
	if pct := doc.Find(selectorChangePct).First().Text(); pct != "" {
		stock.ChangePct = pct + "%"
	}
	stock.Open = getDataListItemValue(doc, labelOpen)
	stock.High = getDataListItemValue(doc, labelHigh)
	stock.Low = getDataListItemValue(doc, labelLow)
	stock.Volume = getDataListItemValue(doc, labelVolume)
}

// getDataListItemValue finds the value paired with the given label in a DataListItem.
func getDataListItemValue(doc *goquery.Document, label string) string {
	var value string
	doc.Find(selectorDataListItemTerm).Each(func(_ int, dt *goquery.Selection) {
		if strings.TrimSpace(dt.Find(selectorDataListItemName).Text()) == label {
			value = FormatPrice(dt.Parent().Find(selectorDataListItemValue).First().Text())
		}
	})
	return value
}
