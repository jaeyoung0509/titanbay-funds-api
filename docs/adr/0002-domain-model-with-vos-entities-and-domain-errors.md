# ADR 0002: Model the domain with value objects, entities, enums, and domain errors

## Status

Accepted

## Context

The earlier implementation mixed transport concerns, persistence concerns, and domain validation. That made the code hard to reason about and encouraged duplication of rules such as date parsing, money formatting, and enum validation.

## Decision

Split the domain into four subpackages:

- `internal/domain/vo` for value objects such as `Money`, `Date`, `Email`, `ID`, and `Timestamp`
- `internal/domain/entity` for aggregate-style business objects such as `Fund`, `Investor`, and `Investment`
- `internal/domain/enum` for fixed business choices such as fund status and investor type
- `internal/domain/error` for domain-level validation, conflict, not-found, and internal errors

Expose money as JSON numbers and dates as `YYYY-MM-DD` strings at the API boundary.

## Consequences

This makes domain rules reusable and explicit, and keeps the HTTP and database layers from owning business validation. The cost is a larger file count, but each file now has a narrow responsibility.
