package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PauloMarcusso/docmind/gateway/internal/app"
	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

// errStoreDown simulates an infrastructure failure, such as Redis being unreachable.
var errStoreDown = errors.New("store down")

// fakeStore is a hand-written app.KeyStore whose answer each test case controls.
// It records whether it was called, so a test can prove the use case skipped it.
type fakeStore struct {
	client domain.Client
	err    error
	called bool
}

// FindClient returns the configured client and error and marks the store as called.
// The pointer receiver makes the called flag survive after the method returns.
func (f *fakeStore) FindClient(ctx context.Context, key string) (domain.Client, error) {
	f.called = true
	return f.client, f.err
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		store      *fakeStore
		wantErr    error
		wantID     string
		wantCalled bool
	}{
		{
			name:       "empty key is rejected without calling the store",
			key:        "",
			store:      &fakeStore{},
			wantErr:    domain.ErrUnauthorized,
			wantCalled: false,
		},
		{
			name:       "unknown key is unauthorized",
			key:        "nope",
			store:      &fakeStore{err: domain.ErrKeyNotFound},
			wantErr:    domain.ErrUnauthorized,
			wantCalled: true,
		},
		{
			name:       "store failure is wrapped, not unauthorized",
			key:        "k",
			store:      &fakeStore{err: errStoreDown},
			wantErr:    errStoreDown,
			wantCalled: true,
		},
		{
			name:       "valid key returns the client",
			key:        "k",
			store:      &fakeStore{client: domain.Client{ID: "dev"}},
			wantID:     "dev",
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auth := app.NewAuthenticator(tc.store)

			client, err := auth.Authenticate(t.Context(), tc.key)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error: got %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != domain.ErrUnauthorized && errors.Is(err, domain.ErrUnauthorized) {
				t.Fatalf("error: got %v, must not be ErrUnauthorized", err)
			}
			if client.ID != tc.wantID {
				t.Errorf("client ID: got %q, want %q", client.ID, tc.wantID)
			}
			if tc.store.called != tc.wantCalled {
				t.Errorf("store called: got %v, want %v", tc.store.called, tc.wantCalled)
			}
		})
	}
}
