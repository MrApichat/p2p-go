package models

type MerchantOrderModel struct {
	Id int64 `json:"id"`
	Type MerchantType `json:"type"`
	Fiat CurrencyModel `json:"fiat"`
	Coin CurrencyModel `json:"coin"`
	Merchant DisplayUserModel `json:"merchant"`
	Price float64 `json:"price"`
	AvailableCoin float64 `json:"available_coin"`
	LowerLimit float64 `json:"lower_limit"`
	Status string `json:"status"`
	PaymentMethods []PaymentMethodModel `json:"payment_methods"`
}

type MerchantRequest struct {
	Type MerchantType `form:"type" json:"type"`
	Fiat string `form:"fiat" json:"fiat"`
	Coin string `form:"coin" json:"coin"`
	Price float64 `form:"price" json:"price"`
	AvailableCoin float64 `form:"available_coin" json:"available_coin"`
	LowerLimit float64 `form:"lower_limit" json:"lower_limit"`
	PaymentMethods []string `form:"payment_methods" json:"payment_methods"`
}

type MerchantType string

const (
	Buy MerchantType = "buy"
	Sell MerchantType = "sell"
)