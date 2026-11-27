package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// writeJSON writes v as a JSON response with the given status code.
// The status is already sent when encoding fails, so the error can only be logged.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write response failed", "err", err)
	}
}
