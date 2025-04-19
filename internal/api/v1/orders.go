package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fungicibus/order/internal/types"
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

	products := make([]types.ProductItem, 0, len(request.Data.Products))

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
				WithSourcePointer(fmt.Sprintf("/data/products/%d/product_id", i)),
			)
			return
		}

		products = append(products, types.ProductItem{
			ProductId: product.Id,
			Quantity:  product.Quantity,
		})
	}

	order := types.Order{
		Comment:   comment,
		Products:  products,
		Timestamp: orderTimestamp,
	}
	createdOrder, err := api.storage.CreateOrder(order)
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
			Id:        createdOrder.Id,
			Products:  request.Data.Products,
			Comment:   request.Data.Comment,
			Timestamp: request.Data.Timestamp,
			Status:    OrderItemStatus(createdOrder.Status),
		},
	}

	api.WriteJSON(w, r, response, http.StatusCreated)
}
