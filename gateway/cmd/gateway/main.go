// Command gateway is the edge service of docmind and the composition root:
// it wires the adapters together and holds no business logic.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PauloMarcusso/docmind/gateway/internal/adapters/inbound/httpapi"
	"github.com/PauloMarcusso/docmind/gateway/internal/adapters/outbound/memory"
	"github.com/PauloMarcusso/docmind/gateway/internal/app"
	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

func main() {
	apiKey := os.Getenv("GATEWAY_API_KEY")
	if apiKey == "" {
		log.Fatal("GATEWAY_API_KEY is required")
	}

	// Wiring goes from the outside in: outbound adapter -> use case -> inbound adapter.
	store := memory.NewKeyStore(map[string]domain.Client{
		apiKey: {ID: "dev"},
	})
	auth := app.NewAuthenticator(store)
	router := httpapi.NewRouter(auth)

	addr := ":8080"
	log.Printf("gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
