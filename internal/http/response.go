package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	engine "github.com/NicolasPaterno/warden-engine"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// respondJSON marshals v to a buffer first, so an encoding failure can still be
// turned into a 500 — once WriteHeader is called the status is locked in.
func respondJSON(ctx context.Context, w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		writeError(ctx, w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// writeError maps a domain error to an HTTP status. Client errors (4xx) echo a
// safe message; server faults (5xx) are logged and attached to the trace span,
// and the client only ever sees a generic message — never internal details.
func writeError(ctx context.Context, w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, engine.ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, engine.ErrInvalid):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		slog.ErrorContext(ctx, "internal server error", "error", err)
		span := trace.SpanFromContext(ctx)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// decodeJSON reads the request body into dst. It returns false (after writing
// the response) on failure: 413 if the body exceeded the size limit set by the
// RequestSize middleware, 400 for malformed JSON.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return false
	}
	return true
}
