package metric

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct{}

func (h *Handler) Register(r chi.Router) {
	r.Get(URL, h.Heartbeat)
}

// Heartbeat handler
// @Summary Heartbeat metric
// @Tags metric
// @Success 204
// @Failure 400
// @Router /api/heartbeat [get]
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	log.Println("heartbeat works")
	w.WriteHeader(http.StatusNoContent)
}
