package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// handleCreateRule serves POST /api/rules.
// The body is the rule to create; ID and CreatedAt are assigned by the service.
func (router *Router) handleCreateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ruleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	rule := req.toEngine()
	if err := rule.Validate(); err != nil {
		writeError(ctx, w, err)
		return
	}
	created, err := router.rules.Create(ctx, rule)
	if err != nil {
		writeError(ctx, w, err)
		return
	}
	respondJSON(ctx, w, http.StatusCreated, newRuleResponse(created))
}

// handleListRules serves GET /api/rules — every rule, newest first.
func (router *Router) handleListRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rules, err := router.rules.GetAll(ctx)
	if err != nil {
		writeError(ctx, w, err)
		return
	}
	respondJSON(ctx, w, http.StatusOK, newRuleResponses(rules))
}

// handleGetRule serves GET /api/rules/{id}.
func (router *Router) handleGetRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rule, err := router.rules.GetByID(ctx, chi.URLParam(r, "id"))
	if err != nil {
		writeError(ctx, w, err)
		return
	}
	respondJSON(ctx, w, http.StatusOK, newRuleResponse(rule))
}

// handleUpdateRule serves PUT /api/rules/{id}.
// The URL is the source of truth for the ID; any ID in the body is ignored.
func (router *Router) handleUpdateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ruleRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	rule := req.toEngine()
	rule.ID = chi.URLParam(r, "id")
	if err := rule.Validate(); err != nil {
		writeError(ctx, w, err)
		return
	}
	updated, err := router.rules.Update(ctx, rule)
	if err != nil {
		writeError(ctx, w, err)
		return
	}
	respondJSON(ctx, w, http.StatusOK, newRuleResponse(updated))
}

// handleDeleteRule serves DELETE /api/rules/{id}.
func (router *Router) handleDeleteRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := router.rules.Delete(ctx, chi.URLParam(r, "id")); err != nil {
		writeError(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
