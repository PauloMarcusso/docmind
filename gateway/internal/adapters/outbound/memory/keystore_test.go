package memory_test

import (
	"errors"
	"testing"

	"github.com/PauloMarcusso/docmind/gateway/internal/adapters/outbound/memory"
	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

func TestKeyStoreFindClient(t *testing.T) {
	store := memory.NewKeyStore(map[string]domain.Client{
		"k": {ID: "dev"},
	})

	tests := []struct {
		name    string
		key     string
		wantID  string
		wantErr error
	}{
		{name: "existing key returns its client", key: "k", wantID: "dev"},
		{name: "missing key returns ErrKeyNotFound", key: "nope", wantErr: domain.ErrKeyNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := store.FindClient(t.Context(), tc.key)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error: got %v, want %v", err, tc.wantErr)
			}
			if client.ID != tc.wantID {
				t.Errorf("client ID: got %q, want %q", client.ID, tc.wantID)
			}
		})
	}
}

// TestNewKeyStoreCopiesMap proves that changing the caller's map after
// construction does not leak into the store.
func TestNewKeyStoreCopiesMap(t *testing.T) {
	clients := map[string]domain.Client{"k": {ID: "dev"}}
	store := memory.NewKeyStore(clients)

	clients["late"] = domain.Client{ID: "intruder"}

	_, err := store.FindClient(t.Context(), "late")
	if !errors.Is(err, domain.ErrKeyNotFound) {
		t.Fatalf("error: got %v, want %v", err, domain.ErrKeyNotFound)
	}
}
