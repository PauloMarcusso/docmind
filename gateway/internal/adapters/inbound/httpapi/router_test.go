package httpapi_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PauloMarcusso/docmind/gateway/internal/adapters/inbound/httpapi"
	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

// fakeAuth is a hand-written httpapi.Authenticator. It ignores the key and
// returns err, so each case decides the auth outcome without the real core.
type fakeAuth struct {
	err error
}

// Authenticate returns the configured error, or a fixed client when err is nil.
func (f fakeAuth) Authenticate(ctx context.Context, key string) (domain.Client, error) {
	if f.err != nil {
		return domain.Client{}, f.err
	}
	return domain.Client{ID: "dev"}, nil
}

// TestRouter checks status, body and content type of every route,
// including the 500 path that guards against a missing return in the middleware.
func TestRouter(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		authErr    error
		wantStatus int
		wantBody   string // empty means: do not check body or content type
	}{
		{
			name:       "healthz is public",
			method:     http.MethodGet,
			path:       "/healthz",
			authErr:    domain.ErrUnauthorized,
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok","service":"gateway"}`,
		},
		{
			name:       "version is public",
			method:     http.MethodGet,
			path:       "/version",
			authErr:    domain.ErrUnauthorized,
			wantStatus: http.StatusOK,
			wantBody:   `{"version":"0.1.0"}`,
		},
		{
			name:       "ping without a valid key is 401",
			method:     http.MethodGet,
			path:       "/v1/ping",
			authErr:    domain.ErrUnauthorized,
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"error":"unauthorized"}`,
		},
		{
			name:       "ping with a valid key is 200",
			method:     http.MethodGet,
			path:       "/v1/ping",
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:       "auth failure is 500 and does not reach the handler",
			method:     http.MethodGet,
			path:       "/v1/ping",
			authErr:    errors.New("auth backend down"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"internal error"}`,
		},
		{
			name:       "wrong method is 405",
			method:     http.MethodPost,
			path:       "/v1/ping",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "unknown route is 404",
			method:     http.MethodGet,
			path:       "/nao-existe",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := httpapi.NewRouter(fakeAuth{err: tc.authErr})
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.wantBody == "" {
				return
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tc.wantBody {
				t.Errorf("body: got %s, want %s", got, tc.wantBody)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", got, "application/json")
			}
		})
	}
}
