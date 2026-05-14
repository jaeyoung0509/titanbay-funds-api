# AGENTS.md

Guidelines for AI coding agents working in this repository.

These rules bias toward small, verifiable changes. For trivial tasks, use judgment, but do not skip verification when behavior changes.

## 1. Project Goal

This repository implements the Titanbay take-home API for managing private market funds, investors, and investments.

The submission should stay focused on the eight required endpoints:

- `GET /funds`
- `POST /funds`
- `PUT /funds`
- `GET /funds/{id}`
- `GET /investors`
- `POST /investors`
- `GET /funds/{fund_id}/investments`
- `POST /funds/{fund_id}/investments`

Do not add transaction endpoints or unrelated product features unless explicitly requested.

## 2. Think Before Coding

Before implementing:

- State assumptions when the request is ambiguous.
- Ask when a missing decision can change the API contract, schema, or architecture.
- Surface tradeoffs briefly when there are multiple reasonable paths.
- Prefer the simplest approach that satisfies the current requirement.

Do not hide uncertainty behind code. If the task is unclear, name what is unclear.

## 3. Simplicity First

Keep changes small and directly tied to the task.

- No speculative features.
- No new abstractions for single-use code.
- No broad refactors unless the task requires them.
- No new dependencies unless they clearly reduce complexity or match an existing decision.
- If a change becomes much larger than expected, pause and explain why.

The target is code that a reviewer can understand quickly.

## 4. Surgical Changes

Touch only the files needed for the task.

- Match the existing style and package boundaries.
- Clean up imports, variables, and functions made unused by your own changes.
- Do not remove unrelated dead code unless asked.
- Do not rewrite docs, formatting, or tests outside the scope of the request.

Every changed line should have a clear reason.

## 5. Architecture Boundaries

Keep the current layered structure intact:

- `internal/domain/vo`: value objects such as `Money`, `Date`, `Email`, `ID`, and `Timestamp`
- `internal/domain/enum`: business enums such as fund status and investor type
- `internal/domain/entity`: domain entities such as `Fund`, `Investor`, and `Investment`
- `internal/domain/error`: domain errors and error kinds
- `internal/port`: repository and application input contracts
- `internal/usecase`: application flow and orchestration
- `internal/adapter/postgres`: PostgreSQL persistence adapter
- `internal/handler`, `internal/request`, `internal/response`: HTTP adapter concerns
- `internal/app`: Fiber app wiring and middleware setup

Rules:

- Domain code must not import HTTP, Fiber, SQL, or sqlc packages.
- Use cases should depend on ports, not concrete persistence adapters.
- Handlers should stay thin: parse input, call use cases, map responses.
- PostgreSQL errors should be translated into domain errors in the adapter.
- Error responses should flow through domain errors and the central Fiber error handler.

## 6. API Contract Rules

Follow the provided API spec and current README.

- Successful responses are raw JSON objects or arrays, not envelopes.
- Money fields must serialize as JSON numbers, not strings.
- `investment_date` must serialize as `YYYY-MM-DD`.
- `created_at` must serialize as an RFC3339 date-time.
- Invalid enum values should return `400 validation_error` with a field-level error.
- Missing resources should return `404 not_found`.
- Duplicate investor email should return `409 conflict`.

## 7. Database Rules

PostgreSQL is part of the correctness model.

- Put schema changes in Goose migrations.
- Keep UUID primary keys, foreign keys, uniqueness, CHECK constraints, and useful indexes explicit.
- Keep money as `NUMERIC(18,2)`.
- Keep date-only values as `DATE`.
- Do not bypass database constraints with application-only assumptions.

## 8. Verification

Define success criteria before changing behavior.

Examples:

- Validation change: add or update a failing validation test, then make it pass.
- API behavior change: cover the endpoint in integration tests.
- Schema change: run migrations through tests or Docker.
- Documentation-only change: inspect rendered Markdown-sensitive sections such as tables, links, and Mermaid blocks.

Default verification:

```bash
go test ./...
```

For local runtime checks:

```bash
docker compose up --build
curl -i http://localhost:8080/health
curl -i http://localhost:8080/swagger
```

## 9. Documentation

Keep documentation aligned with the implementation.

- Update `README.md` when setup, API scope, assumptions, or operational behavior changes.
- Update `docs/swagger/openapi.yaml` when an API contract changes.
- Add or update ADRs when an architectural decision changes.
- Use relative links in Markdown. Do not commit local absolute paths.

## 10. Git Hygiene

- Use conventional commits.
- Prefer small commits grouped by intent.
- Do not commit unrelated untracked files.
- Do not amend or rewrite history unless explicitly requested.
- Before finalizing, check:

```bash
git status --short
go test ./...
```

## 11. Agent Operating Style

Work like a careful senior engineer:

- Read the relevant code before editing.
- Make the smallest useful change.
- Verify the result.
- Report what changed, how it was verified, and any remaining risk.

These guidelines are working if diffs are smaller, tests are more focused, and surprises are surfaced before they become implementation mistakes.
