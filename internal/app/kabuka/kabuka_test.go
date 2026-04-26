package kabuka

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
			want: "symbol,current_price\n3994.T,1234.56\n",
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
