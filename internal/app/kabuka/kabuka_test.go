package kabuka

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/shionit/kabuka/internal/app/kabuka/fetcher/jp"
	_ "github.com/shionit/kabuka/internal/app/kabuka/fetcher/us"

	"github.com/shionit/kabuka/internal/app/kabuka/model"
)

// jpStockHTML returns minimal HTML that satisfies the JP fetcher's CSS selectors.
func jpStockHTML(marketName, price string) string {
	return fmt.Sprintf(`<html><body>
<div id="root"><main><div><section>
  <div><span class="_PriceBoardMenu__label_92n65_18">%s</span></div>
  <div><span class="_StyledNumber__value_1arhg_9">%s</span></div>
</section></div></main></div>
</body></html>`, marketName, price)
}

// usStockHTML returns minimal HTML that satisfies the US fetcher's CSS selectors.
func usStockHTML(marketName, price string) string {
	return jpStockHTML(marketName, price) // identical structure, different market name
}

func TestKabuka_fetch_unit(t *testing.T) {
	tests := []struct {
		name       string
		symbol     string
		html       string
		wantSymbol string
		wantPrice  string
		wantErr    bool
	}{
		{
			name:       "JP stock (東証PRM)",
			symbol:     "3994.T",
			html:       jpStockHTML("東証PRM", "1,234.56"),
			wantSymbol: "3994.T",
			wantPrice:  "1234.56",
		},
		{
			name:       "US stock (NASDAQ)",
			symbol:     "AAPL",
			html:       usStockHTML("NASDAQ", "189.30"),
			wantSymbol: "AAPL",
			wantPrice:  "189.30",
		},
		{
			name:    "unknown market returns error",
			symbol:  "UNKNOWN",
			html:    `<html><body><div id="root"></div></body></html>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The server mirrors Yahoo Finance: /search/?query=SYM redirects to /quote/SYM,
			// which is what isSearchResultsPage relies on to detect non-search pages.
			html := tt.html
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/quote/"+tt.symbol {
					if _, err := fmt.Fprint(w, html); err != nil {
						t.Errorf("failed to write response: %v", err)
					}
				} else {
					http.Redirect(w, r, "/quote/"+tt.symbol, http.StatusFound)
				}
			}))
			defer srv.Close()

			k := &Kabuka{
				Option:  Option{Symbol: tt.symbol},
				baseURL: srv.URL + "/search/?query=",
			}
			got, err := k.fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			want := &model.Stock{Symbol: tt.wantSymbol, CurrentPrice: tt.wantPrice}
			if got.Symbol != want.Symbol {
				t.Errorf("Symbol = %q, want %q", got.Symbol, want.Symbol)
			}
			if got.CurrentPrice != want.CurrentPrice {
				t.Errorf("CurrentPrice = %q, want %q", got.CurrentPrice, want.CurrentPrice)
			}
		})
	}
}

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    OutputFormatType
		wantErr bool
	}{
		{"text", OutputFormatTypeText, false},
		{"", OutputFormatTypeText, false},
		{"json", OutputFormatTypeJson, false},
		{"csv", OutputFormatTypeCsv, false},
		{"xml", "", true},
		{"JSON", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOutputFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFindProductLinkFromSearchResults(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		html    string
		want    string
		wantErr bool
	}{
		{
			name:   "matching link found",
			symbol: "3994.T",
			html:   `<html><body><a href="https://finance.yahoo.co.jp/quote/3994.T">MoneyForward</a></body></html>`,
			want:   "https://finance.yahoo.co.jp/quote/3994.T",
		},
		{
			name:    "link for different symbol does not match",
			symbol:  "3994.T",
			html:    `<html><body><a href="https://finance.yahoo.co.jp/quote/7203.T">Toyota</a></body></html>`,
			wantErr: true,
		},
		{
			name:    "no links in document",
			symbol:  "3994.T",
			html:    `<html><body></body></html>`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				Body: io.NopCloser(strings.NewReader(tt.html)),
			}
			got, err := findProductLinkFromSearchResults(res, tt.symbol, "https://finance.yahoo.co.jp/quote/")
			if (err != nil) != tt.wantErr {
				t.Errorf("findProductLinkFromSearchResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findProductLinkFromSearchResults() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestKabuka_fetch_searchResultsPage(t *testing.T) {
	symbol := "3994.T"

	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/quote/"+symbol {
			fmt.Fprint(w, jpStockHTML("東証PRM", "4208"))
		} else {
			// Return search results HTML with a link — no redirect so isSearchResultsPage returns true
			fmt.Fprintf(w, `<html><body><a href="%s/quote/%s">stock</a></body></html>`, srvURL, symbol)
		}
	}))
	defer srv.Close()
	srvURL = srv.URL

	k := &Kabuka{
		Option:       Option{Symbol: symbol},
		baseURL:      srv.URL + "/search/?query=",
		quoteBaseURL: srv.URL + "/quote/",
	}
	got, err := k.fetch()
	if err != nil {
		t.Fatalf("fetch() unexpected error: %v", err)
	}
	if got.CurrentPrice != "4208" {
		t.Errorf("CurrentPrice = %q, want %q", got.CurrentPrice, "4208")
	}
}

func TestKabuka_fetch_searchResultsPage_noLink(t *testing.T) {
	symbol := "3994.T"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Search results page with no matching link
		fmt.Fprint(w, `<html><body><a href="https://finance.yahoo.co.jp/other">other</a></body></html>`)
	}))
	defer srv.Close()

	k := &Kabuka{
		Option:       Option{Symbol: symbol},
		baseURL:      srv.URL + "/search/?query=",
		quoteBaseURL: srv.URL + "/quote/",
	}
	if _, err := k.fetch(); err == nil {
		t.Error("fetch() expected error when no matching link found, got nil")
	}
}

func TestKabuka_fetch_searchResultsPage_productError(t *testing.T) {
	symbol := "3994.T"

	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/quote/"+symbol {
			w.WriteHeader(http.StatusNotFound)
		} else {
			fmt.Fprintf(w, `<html><body><a href="%s/quote/%s">stock</a></body></html>`, srvURL, symbol)
		}
	}))
	defer srv.Close()
	srvURL = srv.URL

	k := &Kabuka{
		Option:       Option{Symbol: symbol},
		baseURL:      srv.URL + "/search/?query=",
		quoteBaseURL: srv.URL + "/quote/",
	}
	if _, err := k.fetch(); err == nil {
		t.Error("fetch() expected error when product page returns 404, got nil")
	}
}

func TestKabuka_fetch_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	k := &Kabuka{
		Option:  Option{Symbol: "3994.T"},
		baseURL: srv.URL + "/search/?query=",
	}
	if _, err := k.fetch(); err == nil {
		t.Error("fetch() expected error for HTTP 500, got nil")
	}
}

func TestKabuka_fetch_networkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close() // closed immediately so the connection is refused

	k := &Kabuka{
		Option:  Option{Symbol: "3994.T"},
		baseURL: srv.URL + "/search/?query=",
	}
	if _, err := k.fetch(); err == nil {
		t.Error("fetch() expected error for closed server, got nil")
	}
}

func TestKabuka_Execute_fetchError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	k := &Kabuka{
		Option:  Option{Symbol: "3994.T", Format: OutputFormatTypeText},
		baseURL: srv.URL + "/search/?query=",
	}
	if err := k.Execute(); err == nil {
		t.Error("Execute() expected error when fetch fails, got nil")
	}
}

func TestKabuka_Execute(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/quote/3994.T" {
			fmt.Fprint(w, jpStockHTML("東証PRM", "4208"))
		} else {
			http.Redirect(w, r, "/quote/3994.T", http.StatusFound)
		}
	}))
	defer srv.Close()

	k := &Kabuka{
		Option:  Option{Symbol: "3994.T", Format: OutputFormatTypeText},
		baseURL: srv.URL + "/search/?query=",
	}
	if err := k.Execute(); err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}
}

func TestKabuka_formatOutput(t *testing.T) {
	type fields struct {
		Option Option
	}
	type args struct {
		stock *model.Stock
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nomal text",
			fields: fields{
				Option: Option{
					Format: OutputFormatTypeText,
				},
			},
			args: args{
				stock: &model.Stock{
					Symbol:       "3994.T",
					CurrentPrice: "123.45",
				},
			},
			want: "123.45\t3994.T",
		},
		{
			name: "nomal json",
			fields: fields{
				Option: Option{
					Format: OutputFormatTypeJson,
				},
			},
			args: args{
				stock: &model.Stock{
					Symbol:       "3994.T",
					CurrentPrice: "123.45",
				},
			},
			want: "{\"symbol\":\"3994.T\",\"current_price\":\"123.45\"}",
		},
		{
			name: "nomal csv",
			fields: fields{
				Option: Option{
					Format: OutputFormatTypeCsv,
				},
			},
			args: args{
				stock: &model.Stock{
					Symbol:       "3994.T",
					CurrentPrice: "1234.56",
				},
			},
			want: "symbol,current_price,change,change_pct,open,high,low,volume\n3994.T,1234.56,,,,,,\n",
		},
		{
			name: "detail text",
			fields: fields{
				Option: Option{
					Format:     OutputFormatTypeText,
					ShowDetail: true,
				},
			},
			args: args{
				stock: &model.Stock{
					Symbol:       "3994.T",
					CurrentPrice: "4208",
					Change:       "+120",
					ChangePct:    "+2.93%",
					Open:         "4100",
					High:         "4250",
					Low:          "4080",
					Volume:       "341200",
				},
			},
			want: "3994.T\t4208\t+120\t+2.93%\t4100\t4250\t4080\t341200",
		},
		{
			name: "unknown format returns error",
			fields: fields{
				Option: Option{
					Format: "xml",
				},
			},
			args: args{
				stock: &model.Stock{Symbol: "3994.T", CurrentPrice: "100"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kabuka{
				Option: tt.fields.Option,
			}
			got, err := k.formatOutput(tt.args.stock)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatOutput() got = %v, want %v", got, tt.want)
			}
		})
	}
}
