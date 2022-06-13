package models

type Product struct {
	Product_id   int    `json:"-"`
	Uid          string `json:"order_uid"`
	TrackNum     string `json:"track_number"`
	Entry        string `json:"entry"`
	Locale       string `json:"locale"`
	Delivery     Delivery
	Payment      Payment
	Items        []Items
	Signature    string `json:"internal_signature"`
	Customer     string `json:"customer_id"`
	DeliveryServ string `json:"delivery_service"`
	Shardkey     string `json:"shardkey"`
	SmId         int    `json:"sm_id"`
	DateOf       string `json:"date_created"`
	OofShard     string `json:"oof_shard"`
}
