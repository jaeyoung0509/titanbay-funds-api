# ADR 0001: Use Fiber v3 with structured request logging

## Status

Accepted

## Context

The API needs a modern HTTP stack with a clean middleware model, central error handling, request IDs, and production-friendly logs.

## Decision

Use Fiber v3 as the HTTP framework, `requestid` to propagate a correlation ID, and `contrib/v3/zerolog` for structured request and response logs. Keep the access log focused on metadata such as method, URL, status, latency, and request ID. Do not log request or response bodies by default.

## Consequences

The application gets consistent correlation across logs and easier operational debugging. The downside is one more middleware layer, but the structure is still thin and explicit.
