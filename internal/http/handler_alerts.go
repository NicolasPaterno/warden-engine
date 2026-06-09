package http

import (
	"context"
	"net/http"

	engine "github.com/NicolasPaterno/warden-engine"
	"github.com/go-chi/chi/v5"
)

// handleListAlerts serves GET /alerts — every alert, newest first.
func (router *Router) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alerts, err := router.alerts.GetAll(r.Context())
	writeAlerts(ctx, w, alerts, err)
}

// handleAlertsByRoom serves GET /alerts/room/{room}.
func (router *Router) handleAlertsByRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	room := chi.URLParam(r, "room")
	alerts, err := router.alerts.GetByRoom(r.Context(), room)
	writeAlerts(ctx, w, alerts, err)
}

// handleAlertsByRule serves GET /alerts/rule/{ruleID}.
func (router *Router) handleAlertsByRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ruleID := chi.URLParam(r, "ruleID")
	alerts, err := router.alerts.GetByRule(r.Context(), ruleID)
	writeAlerts(ctx, w, alerts, err)
}

// writeAlerts is the shared tail for the alert read handlers: it turns a
// (alerts, err) pair into an HTTP response.
func writeAlerts(ctx context.Context, w http.ResponseWriter, alerts []engine.Alert, err error) {
	if err != nil {
		writeError(ctx, w, err)
		return
	}
	respondJSON(ctx, w, http.StatusOK, newAlertResponses(alerts))
}
