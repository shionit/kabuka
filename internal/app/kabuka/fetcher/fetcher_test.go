package fetcher

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

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
  <div class="PriceBoardMenu__3fnA PriceBoard__menu__ISpY">
    <span>東証PRM</span>
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
  <div class="PriceBoardMenu__3fnA PriceBoard__menu__ISpY">
    <span></span>
    <div class="Popup__34dI PriceBoardMenu__popup__2i88">
      <button>NASDAQ</button>
    </div>
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
  <div class="PriceBoardMenu__3fnA PriceBoard__menu__ISpY">
    <span>TEST_MARKET</span>
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
