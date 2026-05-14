# Architecture Decisions

This folder records the main implementation decisions for the Titanbay Funds API.
Each ADR is short and focused on one design choice so the reasoning is easy to review later.

## Documents

- [ADR 0001: Fiber v3 and structured logging](./0001-fiber-v3-and-structured-logging.md) - why the API uses Fiber v3, request IDs, and zerolog-based access logging.
- [ADR 0002: Domain model with VOs, entities, enums, and domain errors](./0002-domain-model-with-vos-entities-and-domain-errors.md) - how the domain is split to keep money, date, enum, and error rules explicit.
- [ADR 0003: Port and adapter boundaries](./0003-port-and-adapter-boundaries.md) - why handlers call use cases through ports instead of talking to SQL directly.
- [ADR 0004: PostgreSQL, Goose, and sqlc-shaped queries](./0004-postgres-goose-and-sqlc.md) - how schema constraints, migrations, and query access are handled.
- [ADR 0005: Core eight-endpoints scope](./0005-core-eight-endpoints-scope.md) - why only the eight required endpoints are implemented in this submission.

## Reading Order

If you want the shortest path through the design, read the documents in this order:

1. [ADR 0005](./0005-core-eight-endpoints-scope.md)
2. [ADR 0002](./0002-domain-model-with-vos-entities-and-domain-errors.md)
3. [ADR 0003](./0003-port-and-adapter-boundaries.md)
4. [ADR 0004](./0004-postgres-goose-and-sqlc.md)
5. [ADR 0001](./0001-fiber-v3-and-structured-logging.md)
