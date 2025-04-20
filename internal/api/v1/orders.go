package v1

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/fungicibus/order/internal/types"
	"github.com/google/uuid"
)

// Create order
// (POST /orders)
func (api *API) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request CreateOrderRequest
	if err := api.ReadJSON(w, r, &request); err != nil {
		api.WriteError(w, r,
			WithStatusCode(http.StatusBadRequest),
			WithError(fmt.Errorf("failed to read request body: %w", err)),
		)
		return
	}

	if request.Data.UserId == "" {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("user_id must not be empty"),
			WithSourcePointer("/data/user_id"),
		)
		return
	}

	comment := ""
	if request.Data.Comment != nil {
		comment = *request.Data.Comment
	}

	orderTimestamp, err := time.Parse(time.RFC3339, request.Data.Timestamp)
	if err != nil {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("timestamp must be in RFC3339 format"),
			WithSourcePointer("/data/timestamp"),
		)
		return
	}
	if orderTimestamp.After(time.Now()) {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("timestamp must be in the past"),
			WithSourcePointer("/data/timestamp"),
		)
		return
	}

	if len(request.Data.Products) == 0 {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("products must not be empty"),
			WithSourcePointer("/data/products"),
		)
		return
	}

	products := make([]types.ProductItem, 0, len(request.Data.Products))
	var orderTotal float32
	for i, product := range request.Data.Products {
		if product.Id == "" {
			api.WriteError(w, r,
				WithStatusCode(http.StatusUnprocessableEntity),
				WithDetail("product_id must not be empty"),
				WithSourcePointer(fmt.Sprintf("/data/products/%d/product_id", i)),
			)
			return
		}

		if product.Quantity <= 0 {
			api.WriteError(w, r,
				WithStatusCode(http.StatusUnprocessableEntity),
				WithDetail("quantity must be greater than zero"),
				WithSourcePointer(fmt.Sprintf("/data/products/%d/quantity", i)),
			)
			return
		}

		// NOTE: allow price=0 for special items
		if product.Price < 0 {
			api.WriteError(w, r,
				WithStatusCode(http.StatusUnprocessableEntity),
				WithDetail("price must be zero or positive"),
				WithSourcePointer(fmt.Sprintf("/data/products/%d/price", i)),
			)
			return
		}

		orderTotal += product.Price

		products = append(products, types.ProductItem{
			Id:       product.Id,
			Quantity: product.Quantity,
			Price:    product.Price,
		})
	}
	orderTotal = float32(math.Round(float64(orderTotal)*100) / 100)

	// NOTE: allow order_total=0 if order contains only free items
	if request.Data.OrderTotal < 0 {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("order_total must be zero or positive"),
			WithSourcePointer("/data/order_total"),
		)
		return
	}
	if request.Data.OrderTotal != orderTotal {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("order_total must be equal to sum of all products prices"),
			WithSourcePointer("/data/order_total"),
		)
		return
	}

	orderID := uuid.NewString()

	order := types.Order{
		Id:         orderID,
		Products:   products,
		Comment:    comment,
		Timestamp:  orderTimestamp,
		Status:     string(Pending),
		OrderTotal: request.Data.OrderTotal,
		UserId:     request.Data.UserId,
	}

	err = api.storage.CreateOrder(r.Context(), order)
	if err != nil {
		api.WriteError(w, r,
			WithStatusCode(http.StatusInternalServerError),
			WithError(fmt.Errorf("failed to create order: %w", err)),
		)
		return
	}

	// send to kafka

	// then asynchronously consume kafka message with confirmed/rejected
	// when receive message, do notification to user

	response := CreateOrderResponse{
		Data: OrderItem{
			Id:        orderID,
			Products:  request.Data.Products,
			Comment:   request.Data.Comment,
			Timestamp: request.Data.Timestamp,
			Status:    Pending,
		},
	}

	api.WriteJSON(w, r, response, http.StatusCreated)
}

// Get all orders
// (GET /admin/orders)
func (api *API) AdminGetOrders(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth != api.cfg.Security.AdminKey {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnauthorized),
			WithDetail("Authorization token does not match admin's"),
		)
	}

	orders, err := api.storage.GetOrders(r.Context(), types.GetOrdersFilters{})
	if err != nil {
		api.WriteError(w, r,
			WithStatusCode(http.StatusInternalServerError),
			WithError(fmt.Errorf("failed to get orders: %w", err)),
		)
		return
	}

	responseData := make([]OrderItem, 0, len(orders))
	for _, order := range orders {
		products := make([]ProductItem, 0, len(order.Products))
		for _, product := range order.Products {
			products = append(products, ProductItem{
				Id:       product.Id,
				Quantity: product.Quantity,
				Price:    product.Price,
			})
		}

		responseData = append(responseData, OrderItem{
			Id:         order.Id,
			Comment:    &order.Comment,
			Timestamp:  order.Timestamp.Format(time.RFC3339),
			Status:     OrderItemStatus(order.Status),
			OrderTotal: order.OrderTotal,
			UserId:     order.UserId,
			Products:   products,
		})
	}

	response := AdminGetOrdersResponse{
		Data: responseData,
	}

	api.WriteJSON(w, r, response, http.StatusOK)
}
