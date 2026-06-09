package http

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NATSChecker interface {
	IsConnected() bool
}

type HealthHandler struct {
	db   *pgxpool.Pool
	nats NATSChecker
}

func NewHealthHandler(db *pgxpool.Pool, nats NATSChecker) *HealthHandler {
	return &HealthHandler{
		db:   db,
		nats: nats,
	}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, "ok")
	if err != nil {
		return
	}
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(r.Context()); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}
	if !h.nats.IsConnected() {
		http.Error(w, "nats unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "ok")
}
