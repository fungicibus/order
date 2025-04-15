package v1

import "github.com/fungicibus/order/internal/types"

type Storage interface {
	CreateOrder(order types.Order) (string, error)
}
