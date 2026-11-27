package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

// KeyStore finds the client that owns an API key.
// It returns domain.ErrKeyNotFound when the key does not exist.
type KeyStore interface {
	FindClient(ctx context.Context, key string) (domain.Client, error)
}

// Authenticator is the use case that identifies a client by its API key.
type Authenticator struct {
	store KeyStore
}

// NewAuthenticator returns an Authenticator backed by the given store.
func NewAuthenticator(store KeyStore) *Authenticator {
	return &Authenticator{store: store}
}

// Authenticate returns the client that owns key.
// It returns domain.ErrUnauthorized for an empty or unknown key, and a wrapped
// error when the store itself fails, so infrastructure failures are not reported as 401.
func (a *Authenticator) Authenticate(ctx context.Context, key string) (domain.Client, error) {
	if key == "" {
		return domain.Client{}, domain.ErrUnauthorized
	}

	client, err := a.store.FindClient(ctx, key)
	if errors.Is(err, domain.ErrKeyNotFound) {
		return domain.Client{}, domain.ErrUnauthorized
	}
	if err != nil {
		return domain.Client{}, fmt.Errorf("find client: %w", err)
	}

	return client, nil
}
