package httpapi

import "net/http"

type pingResponse struct {
	Status string `json:"status"`
}

// handlePing is a temporary protected route used to check API key auth.
// It goes away once the reverse proxy to the rag-api exists.
func handlePing(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, pingResponse{Status: "ok"})
}
