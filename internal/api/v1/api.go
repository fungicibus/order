package v1

import (
	"net/http"

	"github.com/fungicibus/order/config"
	"github.com/fungicibus/order/internal/logger"
)

type API struct {
	cfg    *config.Config
	logger *logger.Logger

	storage Storage
	queue   Queue
}

func New(cfg *config.Config, logger *logger.Logger, storage Storage, queue Queue) *API {
	return &API{
		cfg:     cfg,
		logger:  logger,
		storage: storage,
		queue:   queue,
	}
}

func (api *API) GetHandler() http.Handler {
	return Handler(api)
}
