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

	if request.Data.ProductId == "" {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("product_id must not be empty"),
			WithSourcePointer("/data/product_id"),
		)
		return
	}

	if request.Data.Quantity <= 0 {
		api.WriteError(w, r,
			WithStatusCode(http.StatusUnprocessableEntity),
			WithDetail("quantity must be greater than zero"),
			WithSourcePointer("/data/quantity"),
		)
		return
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

	order := types.Order{
		Comment:   comment,
		ProductId: request.Data.ProductId,
		Quantity:  request.Data.Quantity,
		Timestamp: orderTimestamp,
	}
	createdId, err := api.storage.CreateOrder(order)
	if err != nil {
		api.WriteError(w, r,
			WithStatusCode(http.StatusInternalServerError),
			WithError(fmt.Errorf("failed to create order: %w", err)),
		)
		return
	}

	response := CreateOrderResponse{
		Data: CreatedOrder{
			Id:        createdId,
			ProductId: request.Data.ProductId,
			Quantity:  request.Data.Quantity,
			Comment:   request.Data.Comment,
			Timestamp: request.Data.Timestamp,
		},
	}

	api.WriteJSON(w, r, response, http.StatusCreated)
}
