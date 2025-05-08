package v1

import (
	"context"

	"github.com/fungicibus/order/internal/types"
)

type Storage interface {
	CreateOrder(ctx context.Context, order types.Order) error
	GetOrders(ctx context.Context, filters types.GetOrdersFilters) ([]types.Order, error)
	Ping(ctx context.Context) error
}

type Queue interface {
	EnqueueOrder(ctx context.Context, order types.Order) error
	Ping(ctx context.Context) error
}

type UnimplementedQueue struct {
}

func (u UnimplementedQueue) EnqueueOrder(ctx context.Context, order types.Order) error {
	return nil
}
