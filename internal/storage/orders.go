package storage

import (
	"context"
	"fmt"

	"github.com/fungicibus/order/internal/types"
	"github.com/jackc/pgx/v5"
)

func (s *service) CreateOrder(ctx context.Context, order types.Order) error {
	q := `
		insert into orders(order_id,"timestamp",order_total,"status")
		values($1, $2, $3, $4)
	`
	_, err := s.rwPool.Exec(
		ctx,
		q,
		order.Id,
		order.Timestamp,
		order.OrderTotal,
		order.Status,
	)
	if err != nil {
		return fmt.Errorf("faield to insert into orders: %w", err)
	}

	q = `
		insert into order_items(order_id, product_id, quantity, price)
		values($1, $2, $3, $4)
	`
	batch := &pgx.Batch{}

	for _, product := range order.Products {
		batch.Queue(q, order.Id, product.Id, product.Quantity, product.Price)
	}
	results := s.rwPool.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert into order_items: batch %d: %w", i, err)
		}
	}
	return nil
}
