package types

import "time"

type Order struct {
	Id         string
	Comment    string
	Timestamp  time.Time
	UserId     string
	Products   []ProductItem
	Status     string
	OrderTotal float32
}

type ProductItem struct {
	Id       string
	Quantity int
	Price    float32
}

type GetOrdersFilters struct {
	UserId  string
	OrderId string
}
