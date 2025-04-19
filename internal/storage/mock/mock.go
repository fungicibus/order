package mock

import "github.com/fungicibus/order/internal/types"

type MockStorage struct{}

func (s MockStorage) CreateOrder(order types.Order) (types.Order, error) {
	return types.Order{}, nil
}
