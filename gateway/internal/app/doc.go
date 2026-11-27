// Package app holds the use cases and the ports (interfaces) they consume.
//
// Dependency rule: it may import domain, and must never import adapters.
// Ports are declared here, by the consumer, and implemented by outbound adapters.
package app
