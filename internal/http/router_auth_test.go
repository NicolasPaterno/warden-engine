package http

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/NicolasPaterno/warden-auth"
	"github.com/NicolasPaterno/warden-auth/authn"
	"github.com/NicolasPaterno/warden-engine/internal/service"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authTestKID      = "engine-test-kid"
	authTestIssuer   = "warden-auth"
	authTestAudience = "warden-engine"
)

type fakeNATS struct{}

func (fakeNATS) IsConnected() bool { return true }

func newAuthedRouter(t *testing.T, key *rsa.PrivateKey) http.Handler {
	t.Helper()
	set := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
		Key:       &key.PublicKey,
		KeyID:     authTestKID,
		Algorithm: "RS256",
		Use:       "sig",
	}}}
	body, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("marshal jwks: %v", err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	verifier := authn.New(srv.URL, authTestIssuer, authTestAudience)
	rules := service.NewRuleService(&fakeRuleRepo{})
	alerts := service.NewAlertService(nil)
	health := NewHealthHandler(nil, fakeNATS{})
	return NewRouter(rules, alerts, verifier, health)
}

func signEngineToken(t *testing.T, key *rsa.PrivateKey, scope string) string {
	t.Helper()
	now := time.Now()
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-1",
			Issuer:    authTestIssuer,
			Audience:  jwt.ClaimStrings{authTestAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
		Scope: scope,
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = authTestKID
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func TestRouterAuth(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	h := newAuthedRouter(t, key)

	get := func(token string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, "/api/rules", nil)
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec
	}

	t.Run("protected route without token returns 401", func(t *testing.T) {
		if rec := get(""); rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})

	t.Run("valid access token passes through to handler", func(t *testing.T) {
		rec := get(signEngineToken(t, key, "access"))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("refresh token is rejected on protected route", func(t *testing.T) {
		if rec := get(signEngineToken(t, key, "refresh")); rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})

	t.Run("health live stays public", func(t *testing.T) {
		rec := doRequest(t, h, http.MethodGet, "/health/live", "")
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}
