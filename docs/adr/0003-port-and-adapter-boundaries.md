# ADR 0003: Use a port and adapter boundary with use cases

## Status

Accepted

## Context

Directly coupling handlers to SQL queries makes the application harder to test and more brittle when persistence changes.

## Decision

Introduce `port` interfaces for repository and input contracts, implement orchestration in `usecase` services, and keep the HTTP layer thin. The PostgreSQL implementation lives in `internal/adapter/postgres` and translates SQL rows and database errors into domain objects and domain errors.

## Consequences

Handlers no longer know about SQL or `sqlc` details. The application can be tested at the use case layer without HTTP, and the persistence adapter can evolve independently.
