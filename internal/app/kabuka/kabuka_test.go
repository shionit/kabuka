package kabuka

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
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
				CurrentPrice: "any price",
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
				CurrentPrice: "any price",
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
				CurrentPrice: "any price",
			},
		},
	}
	opts := []cmp.Option{
		cmp.Comparer(func(want *Stock, got *Stock) bool {
			if want.Symbol != got.Symbol {
				return false
			}
			if want.CurrentPrice == "" {
				return true // ignore Price
			}
			// success if any float number
			_, err := strconv.ParseFloat(got.CurrentPrice, 32)
			return err != nil
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
