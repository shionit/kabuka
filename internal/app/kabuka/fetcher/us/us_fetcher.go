package us

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	selectorCurrentPrice = "#root > main > div > section > div._1nb3c4wQ > header > div.gbs7pEV2 > span > span > span"
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

func (f *usFetcher) Fetch(doc *goquery.Document, symbol string) (*model.Stock, error) {
	curPrice := doc.Find(selectorCurrentPrice).Text()

	return &model.Stock{
		Symbol:       symbol,
		CurrentPrice: fetcher.FormatPrice(curPrice),
	}, nil
}
