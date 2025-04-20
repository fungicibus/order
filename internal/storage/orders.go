package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/fungicibus/order/internal/types"
	"github.com/jackc/pgx/v5"
)

func (s *service) CreateOrder(ctx context.Context, order types.Order) error {
	q := `
		insert into orders(order_id,"timestamp",user_id,order_total,"status")
		values($1, $2, $3, $4, $5)
	`
	_, err := s.rwPool.Exec(
		ctx,
		q,
		order.Id,
		order.Timestamp,
		order.UserId,
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

func (s *service) GetOrders(ctx context.Context, filters types.GetOrdersFilters) ([]types.Order, error) {
	q := `
		select 
			o.order_id ,
			o.user_id ,
			o."timestamp" ,
			o.order_total ,
			o.status ,
			oi.product_id ,
			oi.quantity ,
			oi.price 
		from order_items oi 
		inner join orders o on oi.order_id = o.order_id
		where 1=1
	`

	args := make([]interface{}, 0)
	if filters.OrderId != "" {
		q += " and o.order_id = $1"
		args = append(args, filters.OrderId)
	}
	if filters.UserId != "" {
		q += " and o.user_id = $2"
		args = append(args, filters.UserId)
	}

	rows, err := s.roPool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	defer rows.Close()

	ordersMap := make(map[string]types.Order, 0)

	for rows.Next() {
		var (
			order       types.Order
			orderId     string
			userId      string
			timestamp   time.Time
			order_total float32
			status      string
			product     types.ProductItem
		)

		err := rows.Scan(
			&orderId,
			&userId,
			&timestamp,
			&order_total,
			&status,
			&product.Id,
			&product.Quantity,
			&product.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var ok bool
		if order, ok = ordersMap[orderId]; !ok {
			order = types.Order{
				Id:         orderId,
				UserId:     userId,
				Timestamp:  timestamp,
				OrderTotal: order_total,
				Status:     status,
				Products:   make([]types.ProductItem, 0),
			}
		}
		order.Products = append(order.Products, product)
		ordersMap[orderId] = order
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	orders := make([]types.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		orders = append(orders, order)
	}

	return orders, nil
}
