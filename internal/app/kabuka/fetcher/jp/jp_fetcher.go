package jp

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	financeSiteJpStockPrefix = "https://finance.yahoo.co.jp/quote"

	selectorCurrentPrice = "#root > main > div > div > div.XuqDlHPN > div:nth-child(2) > section._1zZriTjI._2l2sDX5w > div._1nb3c4wQ > header > div.nOmR5zWz > span > span > span"
)

func init() {
	fetcher.RegisterFetcher(&jpFetcher{})
}

type jpFetcher struct {
	fetcher.Fetcher
}

func (f *jpFetcher) IsMarketUrl(url string) bool {
	return strings.HasPrefix(url, financeSiteJpStockPrefix)
}

func (f *jpFetcher) Fetch(doc *goquery.Document, symbol string) (*model.Stock, error) {
	curPrice := doc.Find(selectorCurrentPrice).Text()

	return &model.Stock{
		Symbol:       symbol,
		CurrentPrice: fetcher.FormatPrice(curPrice),
	}, nil
}
