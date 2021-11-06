package kabuka

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	anyPrice = "any price"
)

func TestKabuka_Fetch(t *testing.T) {
	type fields struct {
		Option Option
	}
	type testStock struct {
		*Stock
		isWant bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Stock
		wantErr bool
	}{
		{
			name: "if kabuka 3994.T then Money Forward, Inc.",
			fields: fields{
				Option: Option{
					Symbol: "3994.T",
				},
			},
			want: &Stock{
				Symbol:       "3994.T",
				CurrentPrice: anyPrice,
			},
		},
		{
			name: "if kabuka 4373 then Simplex Holdings, Inc.",
			fields: fields{
				Option: Option{
					Symbol: "4373",
				},
			},
			want: &Stock{
				Symbol:       "4373.T",
				CurrentPrice: anyPrice,
			},
		},
		{
			name: "if kabuka AAPL then Apple Inc.",
			fields: fields{
				Option: Option{
					Symbol: "AAPL",
				},
			},
			want: &Stock{
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
			got, err := k.Fetch()
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
