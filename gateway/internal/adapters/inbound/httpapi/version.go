package httpapi

import "net/http"

// version is the gateway release. It is an operational detail, so it lives
// in the inbound adapter and does not go through the core.
const version = "0.1.0"

type versionResponse struct {
	Version string `json:"version"`
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, versionResponse{Version: version})
}
