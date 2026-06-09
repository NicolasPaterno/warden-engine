package http

import (
	"net/http"

	"github.com/NicolasPaterno/warden-engine/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	rules  *service.RuleService
	alerts *service.AlertService
}

func NewRouter(rules *service.RuleService, alerts *service.AlertService, healthHandler *HealthHandler) http.Handler {
	router := Router{rules: rules, alerts: alerts}
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestSize(1 << 20)) // cap request bodies at 1 MiB

	//routes
	r.Get("/health/live", healthHandler.Live)
	r.Get("/health/ready", healthHandler.Ready)

	r.Route("/api/rules", func(r chi.Router) {
		r.Post("/", router.handleCreateRule)
		r.Get("/", router.handleListRules)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", router.handleGetRule)
			r.Put("/", router.handleUpdateRule)
			r.Delete("/", router.handleDeleteRule)
		})
	})
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", router.handleListAlerts)
		r.Get("/room/{room}", router.handleAlertsByRoom)
		r.Get("/rule/{ruleID}", router.handleAlertsByRule)
	})
	return r
}
