# docmind: working agreement

RAG assistant over technical AI papers: upload PDFs, ask questions, get answers
with citations. It is also a **learning project** built to cover a job
description for a Python + Go AI/backend role (RAG, microservices, containers,
cloud, SQL + NoSQL, observability, CI/CD).

## Services (monorepo)

| Dir              | Language | Role                                      |
|------------------|----------|-------------------------------------------|
| `gateway/`       | Go       | Edge: API key auth, rate limit, reverse proxy, SSE streaming |
| `rag-api/`       | Python   | FastAPI + LangGraph: retrieval, reranking, generation        |
| `ingest-worker/` | Python   | Async PDF extraction, chunking, embeddings                   |

Data: Postgres + pgvector (documents, vectors, hybrid search), MongoDB or
DynamoDB (chat history). Queue: Redis Streams locally, SQS on AWS.
Deploy target: ECS Fargate (no EKS: cost). Kubernetes manifests run on kind/k3d.

## Architecture: Clean Architecture + Hexagonal (Ports & Adapters)

Applies to every service. The dependency rule: **dependencies point inward only.**

| Package                      | May import        | Must never import                |
|------------------------------|-------------------|----------------------------------|
| `domain`                     | nothing in module | frameworks, `app`, `adapters`    |
| `app` (use cases + ports)    | `domain`          | `adapters`                       |
| `adapters/inbound/*`         | `app`, `domain`   | other adapters                   |
| `adapters/outbound/*`        | `app`, `domain`   | `inbound`                        |
| `cmd/*` (composition root)   | everything        | holds no business logic          |

Go-specific rules:
- Interfaces (ports) are **declared by the consumer** in `app`, not by the implementer.
- Keep interfaces small (1-3 methods). No `implements`: satisfaction is implicit.
- Do not create empty layers "for the future". A folder appears when it has code.
- Not everything goes through the core: operational endpoints such as `/healthz` live directly in the inbound adapter.
- Manual dependency injection in `cmd/*/main.go`; no DI framework.
- Standard library first; add a dependency only when it clearly pays off.
- Never name a package like a standard-library one (`httpapi`, not `http`).
- Wrap errors with `%w`; never ignore an error silently; `log.Fatal` only in `main`.
- Propagate `context.Context` as the first parameter of anything doing I/O.

## Teaching mode (important)

The user is learning, mainly Go. Unless they say otherwise:
- **The user writes the code.** Do not hand over finished feature code. Explain the concept, give the contract (names, signatures, expected behavior), let them implement, then review what they send. Scaffolding (directories, `go.mod`, config) is fine to create.
- Explain in **Brazilian Portuguese**, short but very didactic, like a class: goal -> concept -> code -> why each part exists -> best practice -> exercise -> how to test.
- Always say *why* a piece of code or a practice exists, and compare with Python when it helps.
- Code, comments, commit messages, ADRs and READMEs are written in **English**.
- Teach each best practice at the moment it is first used.
- When an approach is not idiomatic Go, say so and show the idiomatic one.

## Workflow

- Conventional Commits (`feat:`, `fix:`, `refactor:`, `docs:`, `chore:`, `test:`).
- Run `make check` (gofmt + vet + tests with `-race`) before committing.
- Architectural decisions are recorded as ADRs in `docs/adr/` (copy `0000-template.md`).
- Never commit secrets. `.env` is ignored; only `.env.example` is versioned.
- Do not push to `main` directly; work on a branch.

## Roadmap

1. Gateway in Go (hexagonal): HTTP server, config + `slog`, graceful shutdown, tests, middlewares, reverse proxy with timeouts, API key auth + rate limit, SSE streaming, multi-stage Dockerfile.
2. RAG core in Python: FastAPI, pgvector, Docker Compose.
3. Distributed: queue, ingest worker, chat history in NoSQL.
4. Quality: LangGraph agent, RAG evaluation, intent classifier (scikit-learn), A/B test of retrieval strategies, OpenTelemetry + Langfuse.
5. Production: Kubernetes manifests, GitHub Actions CI/CD, AWS deploy (optionally Terraform).

## Resuming a session

Course progress, pending exercises and lesson notes live in `.notes/` (local,
git-ignored, so it may not exist in a fresh clone).

1. If `.notes/PROGRESS.md` exists, read it first and continue from "Next step".
2. If it does not exist, inspect the code and `docs/adr/` to infer where we are,
   then **ask the user** which lesson they are on before teaching anything.
