package us

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	financeSiteUsStockPrefix = "https://stocks.finance.yahoo.co.jp/us/detail"

	selectorCurrentPrice = "#main > div.stocksDtlWp > div > div.forAddPortfolio > table > tbody > tr > td:nth-child(3)"
)

func init() {
	fetcher.RegisterFetcher(&usFetcher{})
}

type usFetcher struct {
	fetcher.Fetcher
}

func (f *usFetcher) IsMarketUrl(url string) bool {
	return strings.HasPrefix(url, financeSiteUsStockPrefix)
}

func (f *usFetcher) Fetch(doc *goquery.Document, symbol string) (*model.Stock, error) {
	curPrice := doc.Find(selectorCurrentPrice).Text()

	return &model.Stock{
		Symbol:       symbol,
		CurrentPrice: fetcher.FormatPrice(curPrice),
	}, nil
}
