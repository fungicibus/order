package types

import "time"

type Order struct {
	Id         string        `json:"id"`
	Comment    string        `json:"comment"`
	Timestamp  time.Time     `json:"timestamp"`
	UserId     string        `json:"userId"`
	Products   []ProductItem `json:"products"`
	Status     string        `json:"status"`
	OrderTotal float32       `json:"orderTotal"`
}

type ProductItem struct {
	Id       string  `json:"id"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

type GetOrdersFilters struct {
	UserId  string `json:"userId"`
	OrderId string `json:"orderId"`
}
