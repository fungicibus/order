package v1

import (
	"context"

	"github.com/fungicibus/order/internal/types"
)

type Storage interface {
	CreateOrder(ctx context.Context, order types.Order) error
}
