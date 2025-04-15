package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fungicibus/order/internal/middleware"
)

func (api *API) ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Set max size to prevent exceptionally large requests
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create a decoder with strict settings
	dec := json.NewDecoder(r.Body)

	// Decode the request body into the destination struct
	err := dec.Decode(dst)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Check for additional JSON data
	if dec.More() {
		return errors.New("request body must only contain a single JSON object")
	}

	return nil
}

func (api *API) WriteJSON(w http.ResponseWriter, r *http.Request, response any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		api.WriteError(w, r, WithError(fmt.Errorf("failed to encode response body: %w", err)))
	}
}

func (api *API) WriteError(w http.ResponseWriter, r *http.Request, opts ...ErrorOption) {
	// Apply configs
	config := &ErrorConfig{
		StatusCode: http.StatusInternalServerError,
		Detail:     "Server Error",
	}
	for _, opt := range opts {
		opt(config)
	}

	// Log
	requestId := middleware.GetRequestID(r.Context())
	logEvent := api.logger.Error().
		Int("statusCode", config.StatusCode).
		Str("detail", config.Detail).
		Str("id", requestId)

	if config.Error != nil {
		logEvent.Err(config.Error)
	}

	logEvent.Send()

	// Send http response
	now := time.Now().Format(time.RFC3339)
	errItem := ErrorItem{
		Detail: config.Detail,
		Id:     requestId,
		Status: strconv.Itoa(config.StatusCode),
		Meta: ErrorMeta{
			Timestamp: &now,
		},
	}
	if config.ErrorSourcePounter != "" {
		errItem.Source = &ErrorSource{
			Pointer: config.ErrorSourcePounter,
		}
	}

	response := ErrorResponse{
		Errors: []ErrorItem{errItem},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(config.StatusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.logger.Error().Err(err).Msg("failed to encode error response body")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type ErrorOption func(*ErrorConfig)
type ErrorConfig struct {
	Error              error
	StatusCode         int
	Detail             string
	ErrorSourcePounter string
}

func WithError(err error) ErrorOption {
	return func(c *ErrorConfig) {
		c.Error = err
	}
}

func WithStatusCode(statusCode int) ErrorOption {
	return func(c *ErrorConfig) {
		c.StatusCode = statusCode
		c.Detail = http.StatusText(statusCode)
	}
}

func WithDetail(detail string) ErrorOption {
	return func(c *ErrorConfig) {
		c.Detail = detail
	}
}

func WithSourcePointer(errorSourcePounter string) ErrorOption {
	return func(c *ErrorConfig) {
		c.ErrorSourcePounter = errorSourcePounter
	}
}
