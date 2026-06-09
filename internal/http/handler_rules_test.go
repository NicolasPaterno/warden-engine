package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	engine "github.com/NicolasPaterno/warden-engine"
	"github.com/NicolasPaterno/warden-engine/internal/service"
	"github.com/go-chi/chi/v5"
)

// fakeRuleRepo is an in-memory engine.RuleRepository for handler tests. Each
// field lets a test force the outcome of one method.
type fakeRuleRepo struct {
	saveErr   error
	byID      engine.Rule
	byIDErr   error
	all       []engine.Rule
	allErr    error
	deleteErr error
}

func (f *fakeRuleRepo) Save(context.Context, engine.Rule) error { return f.saveErr }
func (f *fakeRuleRepo) GetByID(context.Context, string) (engine.Rule, error) {
	return f.byID, f.byIDErr
}
func (f *fakeRuleRepo) GetAll(context.Context) ([]engine.Rule, error)     { return f.all, f.allErr }
func (f *fakeRuleRepo) GetEnabled(context.Context) ([]engine.Rule, error) { return f.all, f.allErr }
func (f *fakeRuleRepo) Update(_ context.Context, r engine.Rule) (engine.Rule, error) {
	return r, f.byIDErr
}
func (f *fakeRuleRepo) Delete(context.Context, string) error { return f.deleteErr }

// newTestRouter wires the real services on top of a fake repo and returns a
// chi router with only the rule routes mounted, so handler tests exercise the
// same routing the app uses without needing a database or health handler.
func newTestRouter(repo *fakeRuleRepo) http.Handler {
	router := &Router{rules: service.NewRuleService(repo)}
	r := chi.NewRouter()
	r.Route("/api/rules", func(r chi.Router) {
		r.Post("/", router.handleCreateRule)
		r.Get("/", router.handleListRules)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", router.handleGetRule)
			r.Put("/", router.handleUpdateRule)
			r.Delete("/", router.handleDeleteRule)
		})
	})
	return r
}

func doRequest(t *testing.T, h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, target, nil)
	} else {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

const validRuleBody = `{
	"name": "bedroom hot",
	"room": "bedroom",
	"condition": {"sensor_type": "temperature", "operator": "gt", "threshold": 30},
	"action": {"type": "alert", "payload": "{\"message\":\"hot\",\"severity\":\"warning\"}"},
	"enabled": true
}`

func TestHandleCreateRule(t *testing.T) {
	t.Run("valid body returns 201 with a generated id", func(t *testing.T) {
		rec := doRequest(t, newTestRouter(&fakeRuleRepo{}), http.MethodPost, "/api/rules/", validRuleBody)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body: %s", rec.Code, rec.Body.String())
		}
		var got ruleResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("response is not valid ruleResponse JSON: %v", err)
		}
		if got.ID == "" {
			t.Error("expected the server to assign an id")
		}
		if got.Name != "bedroom hot" {
			t.Errorf("name = %q, want 'bedroom hot'", got.Name)
		}
	})

	t.Run("invalid sensor type returns 400", func(t *testing.T) {
		body := strings.Replace(validRuleBody, `"temperature"`, `"pressure"`, 1)
		rec := doRequest(t, newTestRouter(&fakeRuleRepo{}), http.MethodPost, "/api/rules/", body)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		rec := doRequest(t, newTestRouter(&fakeRuleRepo{}), http.MethodPost, "/api/rules/", "{not json")

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body: %s", rec.Code, rec.Body.String())
		}
	})
}

func TestHandleListRules(t *testing.T) {
	repo := &fakeRuleRepo{all: []engine.Rule{{ID: "r1", Name: "a"}, {ID: "r2", Name: "b"}}}
	rec := doRequest(t, newTestRouter(repo), http.MethodGet, "/api/rules/", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got []ruleResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not a JSON array of rules: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d rules, want 2", len(got))
	}
}

func TestHandleGetRuleNotFound(t *testing.T) {
	repo := &fakeRuleRepo{byIDErr: engine.ErrNotFound}
	rec := doRequest(t, newTestRouter(repo), http.MethodGet, "/api/rules/missing", "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandleDeleteRule(t *testing.T) {
	t.Run("existing rule returns 204", func(t *testing.T) {
		rec := doRequest(t, newTestRouter(&fakeRuleRepo{}), http.MethodDelete, "/api/rules/r1", "")

		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204", rec.Code)
		}
	})

	t.Run("missing rule returns 404", func(t *testing.T) {
		repo := &fakeRuleRepo{deleteErr: engine.ErrNotFound}
		rec := doRequest(t, newTestRouter(repo), http.MethodDelete, "/api/rules/missing", "")

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
	})
}
