package httpapi

import "net/http"

// NewRouter servers all HTTP gateway routes
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)
	return mux
}