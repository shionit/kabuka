package kabuka

import (
	"strconv"
	"testing"

	_ "github.com/shionit/kabuka/internal/app/kabuka/fetcher/jp"
	_ "github.com/shionit/kabuka/internal/app/kabuka/fetcher/us"

	"github.com/shionit/kabuka/internal/app/kabuka/model"

	"github.com/google/go-cmp/cmp"
)

const (
	anyPrice = "any price"
)

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
			name: "if kabuka 3994.T then Money Forward, Inc.",
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
			name: "if kabuka 4373 then Simplex Holdings, Inc.",
			fields: fields{
				Option: Option{
					Symbol: "4373",
				},
			},
			want: &model.Stock{
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
