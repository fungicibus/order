package v1

import "net/http"

// Create order
// (POST /orders)
func (api *API) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request CreateOrderRequest
	if err := api.ReadJSON(w, r, &request); err != nil {
		api.WriteError(w, r, WithStatusCode(http.StatusBadRequest), WithError(err))
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
