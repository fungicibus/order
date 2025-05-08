package v1

import (
	"fmt"
	"net/http"
)

// TODO: подумать где на лучше всего хранить эту функцию.
// Это может быть не метод api, а метод server или еще выше - main
func (api *API) Ready(w http.ResponseWriter, r *http.Request) {
	if err := api.queue.Ping(r.Context()); err != nil {
		api.WriteError(w, r,
			WithStatusCode(http.StatusInternalServerError),
			WithError(fmt.Errorf("failed ping queue: %w", err)),
		)
		return
	}

	if err := api.storage.Ping(r.Context()); err != nil {
		api.WriteError(w, r,
			WithStatusCode(http.StatusInternalServerError),
			WithError(fmt.Errorf("failed ping storage: %w", err)),
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
