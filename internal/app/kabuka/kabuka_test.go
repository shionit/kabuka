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
	}
	opts := []cmp.Option{
		cmp.Comparer(func(a *Stock, b *Stock) bool {
			if a.Symbol != b.Symbol {
				return false
			}
			if a.CurrentPrice == anyPrice {
				// success if any float number
				_, err := strconv.ParseFloat(b.CurrentPrice, 32)
				return err == nil
			} else if b.CurrentPrice == anyPrice {
				// success if any float number
				_, err := strconv.ParseFloat(a.CurrentPrice, 32)
				return err == nil
			}
			return false
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
			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Errorf("Fetch() diff(-want +got): %v", diff)
			}
		})
	}
}
