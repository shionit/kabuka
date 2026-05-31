package fetcher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

func newDetailDoc(change, changePct, open, high, low, volume string) *goquery.Document {
	html := fmt.Sprintf(`<html><body>
<div class="_PriceChangeLabel_hse06_1">
  <span class="_PriceChangeLabel__primary_hse06_56"><span class="_StyledNumber__value_1arhg_9">%s</span></span>
  <span class="_PriceChangeLabel__secondary_hse06_62"><span class="_StyledNumber__value_1arhg_9">%s</span></span>
</div>
<dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">始値</span></dt><dd><span class="_DataListItem__value_1kf95_71">%s</span></dd></dl>
<dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">高値</span></dt><dd><span class="_DataListItem__value_1kf95_71">%s</span></dd></dl>
<dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">安値</span></dt><dd><span class="_DataListItem__value_1kf95_71">%s</span></dd></dl>
<dl class="_DataListItem_1kf95_1"><dt class="_DataListItem__term_1kf95_9"><span class="_DataListItem__name_1kf95_19">出来高</span></dt><dd><span class="_DataListItem__value_1kf95_71">%s</span></dd></dl>
</body></html>`, change, changePct, open, high, low, volume)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc
}

func TestFetchDetailFields(t *testing.T) {
	tests := []struct {
		name          string
		change        string
		changePct     string
		open          string
		high          string
		low           string
		volume        string
		wantChange    string
		wantChangePct string
		wantOpen      string
		wantHigh      string
		wantLow       string
		wantVolume    string
	}{
		{
			name:          "all fields present with commas",
			change:        "+120",
			changePct:     "+2.93",
			open:          "4,100",
			high:          "4,250",
			low:           "4,080",
			volume:        "341,200",
			wantChange:    "+120",
			wantChangePct: "+2.93%",
			wantOpen:      "4100",
			wantHigh:      "4250",
			wantLow:       "4080",
			wantVolume:    "341200",
		},
		{
			name:          "empty document yields empty fields",
			wantChange:    "",
			wantChangePct: "",
			wantOpen:      "",
			wantHigh:      "",
			wantLow:       "",
			wantVolume:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := newDetailDoc(tt.change, tt.changePct, tt.open, tt.high, tt.low, tt.volume)
			stock := &model.Stock{}
			FetchDetailFields(doc, stock)
			if stock.Change != tt.wantChange {
				t.Errorf("Change = %q, want %q", stock.Change, tt.wantChange)
			}
			if stock.ChangePct != tt.wantChangePct {
				t.Errorf("ChangePct = %q, want %q", stock.ChangePct, tt.wantChangePct)
			}
			if stock.Open != tt.wantOpen {
				t.Errorf("Open = %q, want %q", stock.Open, tt.wantOpen)
			}
			if stock.High != tt.wantHigh {
				t.Errorf("High = %q, want %q", stock.High, tt.wantHigh)
			}
			if stock.Low != tt.wantLow {
				t.Errorf("Low = %q, want %q", stock.Low, tt.wantLow)
			}
			if stock.Volume != tt.wantVolume {
				t.Errorf("Volume = %q, want %q", stock.Volume, tt.wantVolume)
			}
		})
	}
}

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no commas", "1234", "1234"},
		{"with commas", "1,234,567", "1234567"},
		{"empty string", "", ""},
		{"only commas", ",,,", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatPrice(tt.input); got != tt.want {
				t.Errorf("FormatPrice(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetMarketName_single(t *testing.T) {
	html := `<html><body>
<div id="root"><main><div><section>
  <div>
    <span class="_PriceBoardMenu__label_92n65_18">東証PRM</span>
  </div>
</section></div></main></div>
</body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}

	got := GetMarketName(doc)
	if got != "東証PRM" {
		t.Errorf("GetMarketName() = %q, want %q", got, "東証PRM")
	}
}

func TestGetMarketName_multi(t *testing.T) {
	html := `<html><body>
<div id="root"><main><div><section>
  <div>
    <button class="_PriceBoardMenu__toggle_92n65_18">NASDAQ<span aria-hidden="true"></span></button>
  </div>
</section></div></main></div>
</body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}

	got := GetMarketName(doc)
	if got != "NASDAQ" {
		t.Errorf("GetMarketName() = %q, want %q", got, "NASDAQ")
	}
}

func TestGetMarketName_empty(t *testing.T) {
	html := `<html><body><div id="root"></div></body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}

	got := GetMarketName(doc)
	if got != "" {
		t.Errorf("GetMarketName() = %q, want empty string", got)
	}
}

func TestSelectFetcher(t *testing.T) {
	// Save and restore fetchers to avoid polluting other tests
	original := fetchers
	defer func() { fetchers = original }()
	fetchers = nil

	html := `<html><body>
<div id="root"><main><div><section>
  <div>
    <span class="_PriceBoardMenu__label_92n65_18">TEST_MARKET</span>
  </div>
</section></div></main></div>
</body></html>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	t.Run("returns nil when no fetchers registered", func(t *testing.T) {
		if got := SelectFetcher(doc); got != nil {
			t.Errorf("SelectFetcher() = %v, want nil", got)
		}
	})

	t.Run("returns matching fetcher", func(t *testing.T) {
		mock := &mockFetcher{market: "TEST_MARKET"}
		RegisterFetcher(mock)

		got := SelectFetcher(doc)
		if got != mock {
			t.Errorf("SelectFetcher() = %v, want %v", got, mock)
		}
	})

	t.Run("returns nil when no fetcher matches", func(t *testing.T) {
		fetchers = nil
		RegisterFetcher(&mockFetcher{market: "OTHER_MARKET"})

		if got := SelectFetcher(doc); got != nil {
			t.Errorf("SelectFetcher() = %v, want nil", got)
		}
	})
}

type mockFetcher struct {
	Fetcher
	market string
}

func (m *mockFetcher) IsMarket(doc *goquery.Document) bool {
	return GetMarketName(doc) == m.market
}
