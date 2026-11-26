package domain

import "errors"

// ErrKeyNotFound means the API key does not exist in the key store.
var ErrKeyNotFound = errors.New("api key not found")

// ErrUnauthorized means the caller is not allowed to use the gateway.
var ErrUnauthorized = errors.New("unauthorized")

// Client is a caller of the gateway, identified by its API key.
type Client struct {
	ID string
}
