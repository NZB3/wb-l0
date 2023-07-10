package models

import "github.com/Rhymond/go-money"

type Payment struct {
	Transaction  string         `json:"transaction"`
	RequestID    string         `json:"request_id"`
	Currency     money.Currency `json:"currency"`
	Provider     string         `json:"provider"`
	Amount       money.Money    `json:"amount"`
	PaymentDT    int            `json:"payment_dt"`
	Bank         string         `json:"bank"`
	DeliveryCost money.Money    `json:"delivery_cost"`
	GoodsTotal   int            `json:"goods_total"`
	CustomFee    money.Money    `json:"custom_fee"`
}
