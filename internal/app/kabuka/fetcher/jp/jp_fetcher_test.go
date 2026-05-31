package jp

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
  <div class="_PriceChangeLabel_hse06_1">
    <span class="_StyledNumber__item_1arhg_6 _PriceChangeLabel__primary_hse06_56"><span class="_StyledNumber__value_1arhg_9">%s</span></span>
    <span class="_StyledNumber__item_1arhg_6 _PriceChangeLabel__secondary_hse06_62"><span class="_StyledNumber__value_1arhg_9">%s</span></span>
  </div>
  <dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">始値</span></dt><dd class="_DataListItem__description_1kf95_50"><span class="_StyledNumber__value_1arhg_9 _DataListItem__value_1kf95_71">%s</span></dd></dl>
  <dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">高値</span></dt><dd class="_DataListItem__description_1kf95_50"><span class="_StyledNumber__value_1arhg_9 _DataListItem__value_1kf95_71">%s</span></dd></dl>
  <dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">安値</span></dt><dd class="_DataListItem__description_1kf95_50"><span class="_StyledNumber__value_1arhg_9 _DataListItem__value_1kf95_71">%s</span></dd></dl>
  <dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">出来高</span></dt><dd class="_DataListItem__description_1kf95_50"><span class="_StyledNumber__value_1arhg_9 _DataListItem__value_1kf95_71">%s</span></dd></dl>
</section></div></main></div>
</body></html>`, marketName, price, change, changePct, open, high, low, volume)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc
}

func TestJpFetcher_IsMarket(t *testing.T) {
	f := &jpFetcher{}

	tests := []struct {
		market string
		want   bool
	}{
		{"東証PRM", true},
		{"東証PRM スタンダード", true},
		{"東証GRT", true},
		{"東証STD", true},
		{"名証MN", true},
		{"札証", true},
		{"札幌ア", true},
		{"福証", true},
		{"福岡Q", true},
		{"NASDAQ", false},
		{"NYSE", false},
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

func TestJpFetcher_Fetch(t *testing.T) {
	f := &jpFetcher{}

	tests := []struct {
		name      string
		price     string
		symbol    string
		wantPrice string
	}{
		{"price with comma", "1,234", "3994.T", "1234"},
		{"price without comma", "4208", "3994.T", "4208"},
		{"empty price", "", "3994.T", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := newDocWithPrice("東証PRM", tt.price)
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

func TestJpFetcher_FetchDetail(t *testing.T) {
	f := &jpFetcher{}
	doc := newDocWithDetail("東証PRM", "4,208", "+120", "+2.93", "4,100", "4,250", "4,080", "341,200")
	stock, err := f.Fetch(doc, "3994.T", true)
	if err != nil {
		t.Fatalf("Fetch() returned error: %v", err)
	}
	if stock.CurrentPrice != "4208" {
		t.Errorf("CurrentPrice = %q, want %q", stock.CurrentPrice, "4208")
	}
	if stock.Change != "+120" {
		t.Errorf("Change = %q, want %q", stock.Change, "+120")
	}
	if stock.ChangePct != "+2.93%" {
		t.Errorf("ChangePct = %q, want %q", stock.ChangePct, "+2.93%")
	}
	if stock.Open != "4100" {
		t.Errorf("Open = %q, want %q", stock.Open, "4100")
	}
	if stock.High != "4250" {
		t.Errorf("High = %q, want %q", stock.High, "4250")
	}
	if stock.Low != "4080" {
		t.Errorf("Low = %q, want %q", stock.Low, "4080")
	}
	if stock.Volume != "341200" {
		t.Errorf("Volume = %q, want %q", stock.Volume, "341200")
	}
}
