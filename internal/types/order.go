package types

import "time"

type Order struct {
	Comment   string
	ProductId string
	Quantity  int
	Timestamp time.Time
}
