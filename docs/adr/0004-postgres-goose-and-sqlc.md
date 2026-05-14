# ADR 0004: Use PostgreSQL constraints, Goose migrations, and sqlc-shaped queries

## Status

Accepted

## Context

The API handles financial data and relational ownership, so data integrity belongs in the database as well as in Go code.

## Decision

Use PostgreSQL as the source of truth for constraints, Goose for schema migrations, and sqlc-shaped query wrappers for explicit SQL access. Enforce UUID primary keys, unique email addresses, foreign keys, positive money checks, enum checks, and bounded vintage year checks in the schema.

## Consequences

The database rejects invalid state even if application code misses a rule. Query behavior stays visible and explicit instead of being hidden behind a large ORM abstraction.
