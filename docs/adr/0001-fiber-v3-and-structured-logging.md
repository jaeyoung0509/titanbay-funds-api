# ADR 0001: Use Fiber v3 with structured request logging

## Status

Accepted

## Context

The API needs a modern HTTP stack with a clean middleware model, central error handling, request IDs, and production-friendly logs.

## Decision

Use Fiber v3 as the HTTP framework, `requestid` to propagate a correlation ID, and `contrib/v3/zerolog` for structured request and response logs. Keep the access log focused on metadata such as method, URL, status, latency, and request ID. Do not log request or response bodies by default.

The API process handles `SIGINT` and `SIGTERM` by calling Fiber's graceful shutdown path with a bounded timeout. Fiber is configured with read, write, and idle timeouts so shutdown is not blocked indefinitely by idle connections.

## Consequences

The application gets consistent correlation across logs, easier operational debugging, and predictable behavior when Docker or a local terminal stops the process. The downside is one more middleware layer and a small amount of process lifecycle code, but the structure is still thin and explicit.
