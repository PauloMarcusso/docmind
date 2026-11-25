package httpapi

import "net/http"

// NewRouter returns the handler that serves all HTTP routes of the gateway.
// It returns the http.Handler interface so middlewares can wrap the mux later.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /version", handleVersion)
	return mux
}
