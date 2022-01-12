package model

// Stock stock info
type Stock struct {
	Symbol       string `json:"symbol" csv:"symbol"`
	CurrentPrice string `json:"current_price" csv:"current_price"`
}
