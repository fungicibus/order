package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

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
	router := chi.NewRouter()
	router.Get("/ready", api.Ready)
	return HandlerFromMux(api, router)
}
