//go:build integration

package kabuka

import (
	"strconv"
	"testing"

	_ "github.com/shionit/kabuka/internal/app/kabuka/fetcher/jp"
	_ "github.com/shionit/kabuka/internal/app/kabuka/fetcher/us"

	"github.com/shionit/kabuka/internal/app/kabuka/model"

	"github.com/google/go-cmp/cmp"
)

// anyPrice is a sentinel used in integration tests to indicate "any non-empty price is acceptable".
const anyPrice = "any price"

func TestKabuka_fetch(t *testing.T) {
	type fields struct {
		Option Option
	}
	type testStock struct {
		*model.Stock
		isWant bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *model.Stock
		wantErr bool
	}{
		{
			name: "if kabuka 3994.T @東証PRM then Money Forward, Inc.",
			fields: fields{
				Option: Option{
					Symbol: "3994.T",
				},
			},
			want: &model.Stock{
				Symbol:       "3994.T",
				CurrentPrice: anyPrice,
			},
		},
		{
			name: "if kabuka 4412 @東証GRT then Science Arts, Inc.",
			fields: fields{
				Option: Option{
					Symbol: "4412",
				},
			},
			want: &model.Stock{
				Symbol:       "4412.T",
				CurrentPrice: anyPrice,
			},
		},
		{
			name: "if kabuka 8604.T @東証PRM then Nomura Holdings, Inc. (IP for multi market)",
			fields: fields{
				Option: Option{
					Symbol: "8604.T",
				},
			},
			want: &model.Stock{
				Symbol:       "8604.T",
				CurrentPrice: anyPrice,
			},
		},
		{
			name: "if kabuka NMR @NYSE then Nomura Holdings, Inc. (IP for multi market)",
			fields: fields{
				Option: Option{
					Symbol: "NMR",
				},
			},
			want: &model.Stock{
				Symbol:       "NMR",
				CurrentPrice: anyPrice,
			},
		},
		{
			name: "if kabuka AAPL @NASDAQ then Apple Inc.",
			fields: fields{
				Option: Option{
					Symbol: "AAPL",
				},
			},
			want: &model.Stock{
				Symbol:       "AAPL",
				CurrentPrice: anyPrice,
			},
		},
		{
			name: "if kabuka WRONG999 then error.",
			fields: fields{
				Option: Option{
					Symbol: "WRONG999",
				},
			},
			wantErr: true,
		},
	}
	opts := []cmp.Option{
		cmp.Comparer(func(a *testStock, b *testStock) bool {
			if a.Symbol != b.Symbol {
				return false
			}
			var want, got *testStock
			if a.isWant {
				want, got = a, b
			} else {
				want, got = b, a
			}
			if want.CurrentPrice == anyPrice {
				if got.CurrentPrice == "---" {
					return true // when market is closed
				}
				// success if any float number
				_, err := strconv.ParseFloat(got.CurrentPrice, 32)
				return err == nil
			}
			return want.CurrentPrice == got.CurrentPrice
		}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kabuka{
				Option: tt.fields.Option,
			}
			got, err := k.fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			want := &testStock{
				tt.want,
				true,
			}
			if diff := cmp.Diff(want, &testStock{got, false}, opts...); diff != "" {
				t.Errorf("Fetch() diff(-want +got): %v", diff)
			}
		})
	}
}
