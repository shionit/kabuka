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
	financeSiteUrl = "https://finance.yahoo.co.jp/quote/"

	selectorCurrentPrice = "#root > main > div > div > div.XuqDlHPN > div:nth-child(2) > section._1zZriTjI._2l2sDX5w > div._1nb3c4wQ > header > div.nOmR5zWz > span > span > span"
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
	paths := strings.Split(res.Request.URL.Path, "/")
	symbol := paths[len(paths)-1]

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("goquery NewDocument failed, err: %w", err)
	}
	curPrice := doc.Find(selectorCurrentPrice).Text()

	return &Stock{
		Symbol:       symbol,
		CurrentPrice: curPrice,
	}, nil
}
