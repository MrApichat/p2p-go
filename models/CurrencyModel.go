package models

type CurrencyModel struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}
