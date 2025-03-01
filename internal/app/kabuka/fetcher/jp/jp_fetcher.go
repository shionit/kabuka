package jp

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	selectorCurrentPrice = "#root > main > div > section > div.PriceBoardMain__1nb3 > div.PriceBoardMain__priceInformation__3YfB > div.PriceBoardMain__headerPrice__gbs7 > span > span > span"
)

var (
	supportMarketNamesPrefix = [...]string{"東証PRM", "東証GRT", "東証STD", "名証MN", "札証", "札幌ア", "福証", "福岡Q"}
)

func init() {
	fetcher.RegisterFetcher(&jpFetcher{})
}

type jpFetcher struct {
	fetcher.Fetcher
}

func (f *jpFetcher) IsMarket(doc *goquery.Document) bool {
	marketName := fetcher.GetMarketName(doc)
	for _, prefix := range supportMarketNamesPrefix {
		if strings.HasPrefix(marketName, prefix) {
			return true
		}
	}
	return false
}

func (f *jpFetcher) Fetch(doc *goquery.Document, symbol string) (*model.Stock, error) {
	curPrice := doc.Find(selectorCurrentPrice).Text()

	return &model.Stock{
		Symbol:       symbol,
		CurrentPrice: fetcher.FormatPrice(curPrice),
	}, nil
}
