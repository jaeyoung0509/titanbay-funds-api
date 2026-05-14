# ADR 0005: Limit the implementation to the eight core endpoints

## Status

Accepted

## Context

The product spec includes additional transaction-related endpoints, but the assignment scope is centered on the core fund, investor, and investment flows.

## Decision

Implement only these eight endpoints:

- `GET /funds`
- `POST /funds`
- `PUT /funds`
- `GET /funds/{id}`
- `GET /investors`
- `POST /investors`
- `GET /funds/{fund_id}/investments`
- `POST /funds/{fund_id}/investments`

Do not expose the transaction-specific endpoints in code or Swagger.

## Consequences

The submission stays aligned with the stated scope and remains simpler to verify. The trade-off is that transaction history is not represented in this version.
