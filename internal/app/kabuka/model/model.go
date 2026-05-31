package model

// Stock stock info
type Stock struct {
	Symbol       string `json:"symbol" csv:"symbol"`
	CurrentPrice string `json:"current_price" csv:"current_price"`
	Change       string `json:"change,omitempty" csv:"change"`
	ChangePct    string `json:"change_pct,omitempty" csv:"change_pct"`
	Open         string `json:"open,omitempty" csv:"open"`
	High         string `json:"high,omitempty" csv:"high"`
	Low          string `json:"low,omitempty" csv:"low"`
	Volume       string `json:"volume,omitempty" csv:"volume"`
}
