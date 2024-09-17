package api

import (
	"github.com/pervukhinpm/link-shortener.git/internal/service"
	"net/http"
)

type DatabaseHealthHandler struct {
	ping *service.PingService
}

func NewDatabaseHealthHandler(
	ping *service.PingService,
) *DatabaseHealthHandler {
	return &DatabaseHealthHandler{
		ping: ping,
	}
}

func (h *DatabaseHealthHandler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	err := h.ping.PingDB(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
