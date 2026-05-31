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

func newDocWithDetail(marketName, price, change, changePct, open, high, low, volume string) *goquery.Document {
	html := fmt.Sprintf(`<html><body>
<div id="root"><main><div><section>
  <div>
    <span class="_PriceBoardMenu__label_92n65_18">%s</span>
  </div>
  <span class="_StyledNumber__value_1arhg_9">%s</span>
  <div class="PriceChangeLabel__2Kf0">
    <span class="StyledNumber__item__1-yu PriceChangeLabel__primary__Y_ut"><span class="StyledNumber__value__3rXW">%s</span></span>
    <span class="StyledNumber__item__1-yu PriceChangeLabel__secondary__3BXI"><span class="StyledNumber__value__3rXW">%s</span></span>
  </div>
  <dl class="DataListItem__38iJ"><dt class="DataListItem__term__30Fb"><span class="DataListItem__name__3RQJ">始値</span></dt><dd class="DataListItem__description__a5Lp"><span class="StyledNumber__value__3rXW DataListItem__value__2wUI">%s</span></dd></dl>
  <dl class="DataListItem__38iJ"><dt class="DataListItem__term__30Fb"><span class="DataListItem__name__3RQJ">高値</span></dt><dd class="DataListItem__description__a5Lp"><span class="StyledNumber__value__3rXW DataListItem__value__2wUI">%s</span></dd></dl>
  <dl class="DataListItem__38iJ"><dt class="DataListItem__term__30Fb"><span class="DataListItem__name__3RQJ">安値</span></dt><dd class="DataListItem__description__a5Lp"><span class="StyledNumber__value__3rXW DataListItem__value__2wUI">%s</span></dd></dl>
  <dl class="DataListItem__38iJ"><dt class="DataListItem__term__30Fb"><span class="DataListItem__name__3RQJ">出来高</span></dt><dd class="DataListItem__description__a5Lp"><span class="StyledNumber__value__3rXW DataListItem__value__2wUI">%s</span></dd></dl>
</section></div></main></div>
</body></html>`, marketName, price, change, changePct, open, high, low, volume)
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
			stock, err := f.Fetch(doc, tt.symbol, false)
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

func TestUsFetcher_FetchDetail(t *testing.T) {
	f := &usFetcher{}
	doc := newDocWithDetail("NASDAQ", "175.00", "-0.45", "-0.14", "174.50", "176.00", "173.80", "70,026,752")
	stock, err := f.Fetch(doc, "AAPL", true)
	if err != nil {
		t.Fatalf("Fetch() returned error: %v", err)
	}
	if stock.CurrentPrice != "175.00" {
		t.Errorf("CurrentPrice = %q, want %q", stock.CurrentPrice, "175.00")
	}
	if stock.Change != "-0.45" {
		t.Errorf("Change = %q, want %q", stock.Change, "-0.45")
	}
	if stock.ChangePct != "-0.14%" {
		t.Errorf("ChangePct = %q, want %q", stock.ChangePct, "-0.14%")
	}
	if stock.Open != "174.50" {
		t.Errorf("Open = %q, want %q", stock.Open, "174.50")
	}
	if stock.High != "176.00" {
		t.Errorf("High = %q, want %q", stock.High, "176.00")
	}
	if stock.Low != "173.80" {
		t.Errorf("Low = %q, want %q", stock.Low, "173.80")
	}
	if stock.Volume != "70026752" {
		t.Errorf("Volume = %q, want %q", stock.Volume, "70026752")
	}
}
