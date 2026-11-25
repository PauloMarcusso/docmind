package httpapi

import "net/http"

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Service: "gateway",
	})
}
