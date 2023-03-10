package models

type TransferOrderModel struct {
	Id       int64         `json:"id"`
	Coin     CurrencyModel `json:"coin"`
	Sender   UserModel     `json:"sender"`
	Receiver UserModel     `json:"receiver"`
	Amount   float64       `json:"amount"`
	Status   string        `json:"status"`
}

type TransferRequest struct {
	ReceiverEmail string  `form:"receiver_email" json:"receiver_email" validate:"required"`
	Coin          string  `form:"coin" json:"coin" validate:"required"`
	Amount        float64 `form:"amount" json:"amount" validate:"required"`
}

type TransferFilter struct {
	Type string 
	Coin string
	Status string
}
