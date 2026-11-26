// Package memory holds outbound adapters backed by in-process data structures,
// useful for local development and tests.
package memory

import (
	"context"

	"github.com/PauloMarcusso/docmind/gateway/internal/app"
	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

// Compile-time check: KeyStore must satisfy the app.KeyStore port.
var _ app.KeyStore = (*KeyStore)(nil)

// KeyStore is an in-memory app.KeyStore. The map is only read after
// construction, so concurrent lookups are safe without a mutex.
type KeyStore struct {
	clients map[string]domain.Client
}

// NewKeyStore returns a KeyStore that maps API keys to clients.
// It copies the map so later changes by the caller do not leak in.
func NewKeyStore(clients map[string]domain.Client) *KeyStore {
	copied := make(map[string]domain.Client, len(clients))
	for key, client := range clients {
		copied[key] = client
	}
	return &KeyStore{clients: copied}
}

// FindClient returns the client that owns key, or domain.ErrKeyNotFound.
func (s *KeyStore) FindClient(ctx context.Context, key string) (domain.Client, error) {
	client, ok := s.clients[key]
	if !ok {
		return domain.Client{}, domain.ErrKeyNotFound
	}
	return client, nil
}
