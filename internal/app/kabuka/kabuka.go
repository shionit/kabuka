package kabuka

import (
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/PuerkitoBio/goquery"
)

const (
	financeSiteUrl           = "https://info.finance.yahoo.co.jp/search/?query="
	financeSiteJpStockPrefix = "https://finance.yahoo.co.jp/quote"
	financeSiteUsStockPrefix = "https://stocks.finance.yahoo.co.jp/us/detail"

	selectorCurrentPriceJp = "#root > main > div > div > div.XuqDlHPN > div:nth-child(2) > section._1zZriTjI._2l2sDX5w > div._1nb3c4wQ > header > div.nOmR5zWz > span > span > span"
	selectorCurrentPriceUs = "#main > div.stocksDtlWp > div > div.forAddPortfolio > table > tbody > tr > td:nth-child(3)"
)

var (
	marketTypes = map[string]marketType{
		financeSiteJpStockPrefix: jp,
		financeSiteUsStockPrefix: us,
	}
)

// Fetch stock information from finance website.
func (k *Kabuka) Fetch() (*Stock, error) {
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
	market, err := parseMarketType(res)
	if err != nil {
		return nil, err
	}
	paths := strings.Split(res.Request.URL.Path, "/")
	symbol := paths[len(paths)-1]

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("goquery NewDocument failed, err: %w", err)
	}
	curPrice := ""
	switch market {
	case jp:
		curPrice = doc.Find(selectorCurrentPriceJp).Text()

	case us:
		curPrice = doc.Find(selectorCurrentPriceUs).Text()
	}

	return &Stock{
		Symbol:       symbol,
		CurrentPrice: formatPrice(curPrice),
	}, nil
}

func isSymbolNotFound(res *http.Response) bool {
	url := res.Request.URL.String()
	// If location is back to search page, result is "not found".
	return strings.HasPrefix(url, financeSiteUrl)
}

func parseMarketType(res *http.Response) (marketType, error) {
	url := res.Request.URL.String()
	for prefix, market := range marketTypes {
		if strings.HasPrefix(url, prefix) {
			return market, nil
		}
	}
	return unknown, xerrors.New("Unknown market type.")
}

func formatPrice(s string) string {
	return strings.ReplaceAll(s, ",", "")
}
