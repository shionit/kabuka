package us

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	selectorCurrentPrice = "[class*='StyledNumber__value']"
)

var (
	supportMarketNames = [...]string{"NASDAQ", "NYSE"}
)

func init() {
	fetcher.RegisterFetcher(&usFetcher{})
}

type usFetcher struct {
	fetcher.Fetcher
}

func (f *usFetcher) IsMarket(doc *goquery.Document) bool {
	marketName := fetcher.GetMarketName(doc)
	for _, name := range supportMarketNames {
		if marketName == name {
			return true
		}
	}
	return false
}

func (f *usFetcher) Fetch(doc *goquery.Document, symbol string, detail bool) (*model.Stock, error) {
	stock := &model.Stock{
		Symbol:       symbol,
		CurrentPrice: fetcher.FormatPrice(doc.Find(selectorCurrentPrice).First().Text()),
	}
	if detail {
		fetcher.FetchDetailFields(doc, stock)
	}
	return stock, nil
}
