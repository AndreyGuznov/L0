package models

type Payment struct {
	Id          int    `json:"-"`
	Transaction string `json:"transaction"`
	Request     string `json:"request_id"`
	Currency    string `json:"currency"`
	Provider    string `json:"provider"`
	Amount      int    `json:"amount"`
	Paymen      int    `json:"payment_dt"`
	Bank        string `json:"bank"`
	Deliver     int    `json:"delivery_cost"`
	Goods       int    `json:"goods_total"`
	Custom      int    `json:"custom_fee"`
}
