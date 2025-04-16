package types

import "time"

type Order struct {
	Comment   string
	Timestamp time.Time
	Products  []ProductItem
}

type ProductItem struct {
	ProductId string
	Quantity  int
}
