package models

type WalletModel struct {
	Id      int           `json:"id"`
	Total   float64       `json:"total"`
	InOrder float64       `json:"in_order"`
	Coin    CurrencyModel `json:"coin"`
}
