package us

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func newDocWithMarket(marketName string) *goquery.Document {
	html := fmt.Sprintf(`<html><body>
<div id="root"><main><div><section>
  <div>
    <span class="_PriceBoardMenu__label_92n65_18">%s</span>
  </div>
</section></div></main></div>
</body></html>`, marketName)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc
}

func newDocWithPrice(marketName, price string) *goquery.Document {
	html := fmt.Sprintf(`<html><body>
<div id="root"><main><div><section>
  <div>
    <span class="_PriceBoardMenu__label_92n65_18">%s</span>
  </div>
  <div>
    <span class="_StyledNumber__value_1arhg_9">%s</span>
  </div>
</section></div></main></div>
</body></html>`, marketName, price)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc
}

func TestUsFetcher_IsMarket(t *testing.T) {
	f := &usFetcher{}

	tests := []struct {
		market string
		want   bool
	}{
		{"NASDAQ", true},
		{"NYSE", true},
		{"東証PRM", false},
		{"東証GRT", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.market, func(t *testing.T) {
			doc := newDocWithMarket(tt.market)
			if got := f.IsMarket(doc); got != tt.want {
				t.Errorf("IsMarket(%q) = %v, want %v", tt.market, got, tt.want)
			}
		})
	}
}

func TestUsFetcher_Fetch(t *testing.T) {
	f := &usFetcher{}

	tests := []struct {
		name      string
		price     string
		symbol    string
		wantPrice string
	}{
		{"price with comma", "1,234.56", "AAPL", "1234.56"},
		{"price without comma", "175.00", "AAPL", "175.00"},
		{"empty price", "", "AAPL", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := newDocWithPrice("NASDAQ", tt.price)
			stock, err := f.Fetch(doc, tt.symbol)
			if err != nil {
				t.Fatalf("Fetch() returned error: %v", err)
			}
			if stock.Symbol != tt.symbol {
				t.Errorf("stock.Symbol = %q, want %q", stock.Symbol, tt.symbol)
			}
			if stock.CurrentPrice != tt.wantPrice {
				t.Errorf("stock.CurrentPrice = %q, want %q", stock.CurrentPrice, tt.wantPrice)
			}
		})
	}
}
