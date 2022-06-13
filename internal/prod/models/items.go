package models

type Items struct {
	Id     int    `json:"-"`
	Chrt   int    `json:"chrt_id"`
	Number string `json:"track_number"`
	Price  int    `json:"price"`
	Rid    string `json:"rid"`
	NameOf string `json:"name"`
	Sale   int    `json:"sale"`
	Size   string `json:"size"`
	Total  int    `json:"total_price"`
	Nm     int    `json:"nm_id"`
	Brand  string `json:"brand"`
	Status int    `json:"status"`
}
