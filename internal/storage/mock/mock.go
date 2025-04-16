package mock

import "github.com/fungicibus/order/internal/types"

type MockStorage struct{}

func (s MockStorage) CreateOrder(order types.Order) (string, error) {
	return "mock-order-id", nil
}
