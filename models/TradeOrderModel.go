package models

type TradeOrderModel struct {
	Id            int64              `json:"id"`
	User          DisplayUserModel   `json:"user"`
	MerchantOrder MerchantOrderModel `json:"merchant_order"`
	Amount        float64            `json:"amount"`
	PaymentMethod PaymentMethodModel `json:"payment_method"`
	TotalPrice    float64            `json:"total_price"`
	Status        TradeStatus        `json:"status"`
}

type TradeStatus string

const (
	Open     TradeStatus = "open"     //open order waiting to trade fiat
	Waiting  TradeStatus = "waiting"  //already sent fiat waiting coin owner confirm
	Complete TradeStatus = "complete" //trader and merchant recieve what they want
	Cancel   TradeStatus = "cancel"   //you can cancel order when order is create only
)
