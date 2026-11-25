// Command gateway is the edge service of docmind and the composition root:
// it wires the adapters together and holds no business logic.
package main

import (
	"log"
	"net/http"

	"github.com/PauloMarcusso/docmind/gateway/internal/adapters/inbound/httpapi"
)

func main() {
	router := httpapi.NewRouter()

	addr := ":8080"
	log.Printf("gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
