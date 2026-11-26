package httpapi

import "net/http"

// NewRouter returns the handler that serves all HTTP routes of the gateway.
// Operational routes are public; /v1/* routes require a valid API key.
func NewRouter(auth Authenticator) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /version", handleVersion)
	mux.Handle("GET /v1/ping", requireAPIKey(auth, http.HandlerFunc(handlePing)))
	return mux
}
