package main

import (
	"log"
	"net/http"

	"github.com/PauloMarcusso/gateway/internal/adapters/inbound/httpapi"
)

func main() {

	router := httpapi.NewRouter()

	addr := ":8080"
	log.Printf("gateway listening in %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
