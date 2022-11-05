package jp

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	selectorCurrentPrice = "#root > main > div > div > div.XuqDlHPN div section._1zZriTjI._2l2sDX5w > div._1nb3c4wQ > header > div.nOmR5zWz > span > span > span"
)

var (
	supportMarketNamesPrefix = [...]string{"東証", "名証", "札幌", "福岡"}
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
