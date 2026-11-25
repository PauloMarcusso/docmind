// Package httpapi is the inbound HTTP adapter: routes, handlers and middlewares.
//
// Dependency rule: it may import app and domain, and must never import other adapters.
// Operational endpoints such as /healthz live here and do not go through the core.
package httpapi
