// Package outbound groups the driven adapters: the code the application calls,
// such as Redis, in-memory stores or the rag-api client.
//
// Dependency rule: adapters may import app and domain, and must never import inbound.
package outbound
