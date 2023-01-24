package models

type TradeOrderModel struct {
	Id int64
	User DisplayUserModel
	MerchantOrder MerchantOrderModel
	Amount float64
	PaymentMethod PaymentMethodModel
	TotalPrice float64
	Status TradeStatus
}

type TradeStatus string

const (
	Open TradeStatus = "open" //open order waiting to trade fiat
	Waiting TradeStatus = "waiting" //already sent fiat waiting coin owner confirm
	Complete TradeStatus = "complete" //trader and merchant recieve what they want
	Cancel TradeStatus = "cancel" //you can cancel order when order is create only
)