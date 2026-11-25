# 0001. Hexagonal architecture with Clean layering

- Status: accepted
- Date: 2026-11-25

## Context
docmind swaps technologies often (Redis vs DynamoDB, SQS vs Redis Streams, one LLM vs another) and must be testable without real infrastructure.

## Decision
Every service uses Ports & Adapters with Clean layering: `domain`, `app` (use cases and ports), `adapters/inbound`, `adapters/outbound`, and `cmd` as composition root. Dependencies point inward only; ports are declared by the consumer in `app`.

## Consequences
Technologies are replaced by writing a new adapter, and use cases are tested with fakes. The cost is more packages and manual wiring in `main`. Operational endpoints such as `/healthz` stay in the inbound adapter. The rule is enforced by review now and by lint (`depguard`) later.
