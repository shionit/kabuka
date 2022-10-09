package kabuka

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
	"github.com/goccy/go-json"
	"github.com/shionit/kabuka/internal/app/kabuka/fetcher"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

const (
	financeSiteUrl = "https://finance.yahoo.co.jp/search/?query="
)

func (k *Kabuka) Execute() error {
	result, err := k.fetch()
	if err != nil {
		return err
	}
	output, err := k.formatOutput(result)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

// fetch stock information from finance website.
func (k *Kabuka) fetch() (*model.Stock, error) {
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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("goquery NewDocument failed, err: %w", err)
	}

	f := fetcher.SelectFetcher(doc)
	if f == nil {
		return nil, xerrors.New("Unknown market type.")
	}

	paths := strings.Split(res.Request.URL.Path, "/")
	symbol := paths[len(paths)-1]

	return f.Fetch(doc, symbol)
}

func isSymbolNotFound(res *http.Response) bool {
	url := res.Request.URL.String()
	// If location is back to search page, result is "not found".
	return strings.HasPrefix(url, financeSiteUrl)
}

// formatOutput format output string.
func (k *Kabuka) formatOutput(stock *model.Stock) (string, error) {
	var result string

	switch k.Option.Format {
	case OutputFormatTypeText:
		result = fmt.Sprintf("%s\t%s", stock.CurrentPrice, stock.Symbol)
	case OutputFormatTypeJson:
		b, err := json.Marshal(stock)
		if err != nil {
			return "", xerrors.Errorf("json Marshal failed, err: %w", err)
		}
		result = string(b)
	case OutputFormatTypeCsv:
		stocks := []*model.Stock{stock}
		csvContent, err := gocsv.MarshalString(&stocks)
		if err != nil {
			return "", xerrors.Errorf("csv Marshal failed, err: %w", err)
		}
		result = csvContent
	default:
		return "", xerrors.New("Unknown formatOutput format.")
	}
	return result, nil
}
