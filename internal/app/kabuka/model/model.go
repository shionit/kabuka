package model

// Stock stock info
type Stock struct {
	Symbol       string `json:"symbol"`
	CurrentPrice string `json:"current_price"`
}
