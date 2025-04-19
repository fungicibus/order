package types

import "time"

type Order struct {
	Id        string
	Comment   string
	Timestamp time.Time
	Products  []ProductItem
	Status    string
}

type ProductItem struct {
	ProductId string
	Quantity  int
}
