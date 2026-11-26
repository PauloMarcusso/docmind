package httpapi

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

// apiKeyHeader is the request header that carries the client's API key.
const apiKeyHeader = "X-API-Key"

// Authenticator is the port the HTTP adapter needs from the core.
// It is declared here, by its consumer, and satisfied by *app.Autenticator.
type Authenticator interface {
	Authenticate(ctx context.Context, key string) (domain.Client, error)
}

type errorResponse struct {
	Error string `json: "error"`
}

// requireApiKey lets the request reach next only when the API key is valid.
// Unkown keys get 401; any other failure gest 500 without leaking details
func requireAPIKey(auth Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.Authenticate(r.Context(), r.Header.Get(apiKeyHeader))
		if errors.Is(err, domain.ErrUnauthorized) {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
			return 
		}
		if err != nil {
			log.Printf("authenticate: %v", err)
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		}
		next.ServeHTTP(w,r)
	})
}


