응 맞아. PostgreSQL은 과제에서 명시 요구사항이야. “Database: PostgreSQL with proper schema design and relationships”라고 되어 있으니까 DB 선택지는 사실상 Postgres 고정이고, 네가 말한 방향이 딱 좋아.

그리고 API spec은 공식 문서 기준으로 Funds 4개, Investors 2개, Investments 2개가 핵심 8개 endpoint야. 문서에는 transaction endpoint도 추가로 있지만, 과제 본문이 말한 “all 8 endpoints”는 Funds/Investors/Investments로 보는 게 맞음. Funds/Investors/Investments는 각각 fund, investor, investment CRUD성 API이고, 모델 필드도 target_size_usd, amount_usd 같은 decimal money field를 포함함. ￼

아래는 네가 말한 방향까지 반영한 완성형 PRD + ADR 구성안이야.

⸻

PRD: Titanbay Private Markets API

1. Objective

Build a backend REST API for managing private market funds, investors, and investor commitments using Go and PostgreSQL.

The service should be:

- Spec-compliant with the provided Titanbay API document
- Easy to run locally with Docker
- Backed by PostgreSQL
- Safe for monetary values
- Well-tested with unit and integration tests
- Documented with Swagger/OpenAPI and ADRs
- Clear about AI-assisted development

Primary goal:

Implement the 8 core endpoints for Funds, Investors, and Investments.

Non-goal:

Do not implement the additional Transactions endpoints unless explicitly required later.

Reason: the task brief says “all 8 endpoints,” and the spec’s first three sections define exactly 8 endpoints: 4 Funds, 2 Investors, 2 Investments. The transaction section is extra and contains additional endpoints beyond those 8. ￼

⸻

2. Tech Stack

Chosen stack

Language: Go
HTTP framework: Fiber
Database: PostgreSQL
SQL access: sqlc + pgx
Migration: goose
Money/decimal: shopspring/decimal
API docs: Swagger/OpenAPI
Testing:

- Unit tests
- Integration tests
- testcontainers-go with PostgreSQL
  Local runtime:
- Docker Compose

Why this stack

Go:
Matches the company’s backend direction and is suitable for small, reliable APIs.
Fiber:
Lightweight, fast, and simple REST routing.
PostgreSQL:
Required by the task and appropriate for relational fund/investor/investment data.
sqlc:
Keeps SQL explicit while providing generated type-safe Go code.
goose:
Simple migration workflow with plain SQL migrations.
shopspring/decimal:
Avoids float64 precision issues for fund sizes and investment amounts.
Swagger/OpenAPI:
Makes the API contract easy to inspect and test locally.
testcontainers-go:
Allows integration tests against a real PostgreSQL database instead of mocks.

⸻

3. Scope

3.1 In scope

Funds

GET /funds
POST /funds
PUT /funds
GET /funds/{id}

The spec defines a Fund with:

id
name
vintage_year
target_size_usd
status
created_at

Valid statuses:

Fundraising
Investing
Closed

The spec shows GET /funds, POST /funds, PUT /funds, and GET /funds/{id} with raw JSON object/array responses. ￼

⸻

Investors

GET /investors
POST /investors

The spec defines an Investor with:

id
name
investor_type
email
created_at

Valid investor types:

Individual
Institution
Family Office

The spec shows GET /investors and POST /investors. ￼

⸻

Investments

GET /funds/{fund_id}/investments
POST /funds/{fund_id}/investments

The spec defines an Investment with:

id
investor_id
fund_id
amount_usd
investment_date

fund_id comes from the path, while investor_id, amount_usd, and investment_date come from the request body for creation. ￼

⸻

3.2 Out of scope

The spec page also includes transaction-related endpoints:

GET /transactions
POST /transactions/process
PUT /transactions/{transaction_id}/reverse
GET /funds/{fund_id}/total-value
POST /admin/recalculate-fees

These are not part of the core 8 endpoints and should be documented as out of scope for the initial submission. ￼

README wording:

The API specification page includes additional transaction-related endpoints. I treated them as out of scope because the task brief asks for all 8 endpoints, which correspond to Funds, Investors, and Investments.

⸻

4. Functional Requirements

4.1 Funds

FR-F1: List funds

Endpoint:

GET /funds

Acceptance criteria:

- Returns 200 OK
- Returns raw JSON array
- Empty database returns []
- Each fund includes:
  id, name, vintage_year, target_size_usd, status, created_at
- Results are ordered by created_at DESC

Response shape:

[
{
"id": "550e8400-e29b-41d4-a716-446655440000",
"name": "Titanbay Growth Fund I",
"vintage_year": 2024,
"target_size_usd": 250000000.00,
"status": "Fundraising",
"created_at": "2024-01-15T10:30:00Z"
}
]

⸻

FR-F2: Create fund

Endpoint:

POST /funds

Request:

{
"name": "Titanbay Growth Fund II",
"vintage_year": 2025,
"target_size_usd": 500000000.00,
"status": "Fundraising"
}

Acceptance criteria:

- Returns 201 Created
- Generates UUID id
- Automatically sets created_at
- Validates required fields
- Validates target_size_usd as positive decimal
- Validates status enum
- Persists amount as PostgreSQL NUMERIC(18,2)
- Returns created fund in spec shape

⸻

FR-F3: Update fund

Endpoint:

PUT /funds

Request:

{
"id": "550e8400-e29b-41d4-a716-446655440000",
"name": "Titanbay Growth Fund I",
"vintage_year": 2024,
"target_size_usd": 300000000.00,
"status": "Investing"
}

Acceptance criteria:

- Returns 200 OK
- Updates existing fund
- Returns 400 for invalid UUID/body
- Returns 404 if fund does not exist
- Validates all mutable fields
- Preserves created_at

Design note:

Although REST APIs often use PUT /funds/{id}, this implementation follows the provided spec, which defines PUT /funds with id in the request body.

⸻

FR-F4: Get fund by ID

Endpoint:

GET /funds/{id}

Acceptance criteria:

- Returns 200 OK for existing fund
- Returns 400 for invalid UUID
- Returns 404 for unknown UUID
- Returns raw fund object

⸻

4.2 Investors

FR-I1: List investors

Endpoint:

GET /investors

Acceptance criteria:

- Returns 200 OK
- Returns raw JSON array
- Empty database returns []
- Each investor includes:
  id, name, investor_type, email, created_at
- Results are ordered by created_at DESC

⸻

FR-I2: Create investor

Endpoint:

POST /investors

Request:

{
"name": "CalPERS",
"investor_type": "Institution",
"email": "privateequity@calpers.ca.gov"
}

Acceptance criteria:

- Returns 201 Created
- Generates UUID id
- Automatically sets created_at
- Validates required fields
- Validates investor_type enum
- Validates email format
- Enforces unique email
- Duplicate email returns 409 Conflict

⸻

4.3 Investments

FR-IV1: List investments for fund

Endpoint:

GET /funds/{fund_id}/investments

Acceptance criteria:

- Returns 200 OK
- Returns raw JSON array
- Returns 400 for invalid fund_id UUID
- Returns 404 if fund does not exist
- Empty fund returns []
- Results are ordered by investment_date DESC, id DESC

⸻

FR-IV2: Create investment for fund

Endpoint:

POST /funds/{fund_id}/investments

Request:

{
"investor_id": "880e8400-e29b-41d4-a716-446655440003",
"amount_usd": 75000000.00,
"investment_date": "2024-09-22"
}

Acceptance criteria:

- Returns 201 Created
- fund_id is read from path
- investor_id is read from body
- Validates both IDs as UUIDs
- Returns 404 if fund does not exist
- Returns 404 if investor does not exist
- Validates amount_usd as positive decimal
- Validates investment_date as YYYY-MM-DD
- Persists amount_usd as PostgreSQL NUMERIC(18,2)

⸻

5. Non-functional Requirements

NFR-1: Local development

The project must run locally with minimal setup.

Required:

docker compose up --build

Expected services:

api
postgres

Optional but recommended:

make migrate-up
make seed
make test

⸻

NFR-2: API documentation

Add Swagger/OpenAPI.

Recommended approach:

Use swaggo/swag with Fiber-compatible Swagger handler.

Endpoints:

GET /swagger/\*

Acceptance criteria:

- Swagger UI is available locally
- All 8 endpoints are documented
- Request/response DTOs are documented
- Common error response is documented

README should include:

Swagger UI is available at http://localhost:8080/swagger/index.html

⸻

NFR-3: Data integrity

Database must enforce integrity, not only application code.

Required constraints:

- UUID primary keys
- NOT NULL constraints
- Foreign keys from investments to funds and investors
- CHECK constraints for fund status
- CHECK constraints for investor type
- CHECK constraints for positive amounts
- UNIQUE constraint for investor email

⸻

NFR-4: Monetary precision

No persistent financial value should be represented as float64.

Required:

- PostgreSQL NUMERIC(18,2)
- decimal.Decimal in Go DTO/query layer where practical
- Validation for positive amounts
- Validation for max 2 decimal places

README should state:

Financial amounts are stored as PostgreSQL NUMERIC(18,2) and represented in Go using decimal.Decimal to avoid floating-point precision issues.

⸻

NFR-5: Error handling

Use common domain error infrastructure.

Required:

- domain error type
- Central Fiber error middleware
- Consistent JSON error responses
- Correct HTTP status mapping

Error response:

{
"error": {
"code": "validation_error",
"message": "validation failed",
"fields": {
"status": "must be one of Fundraising, Investing, Closed"
}
}
}

Error codes:

validation_error
not_found
conflict
internal_error

⸻

NFR-6: Testing

Implement both unit and integration tests.

Required testing layers:

Unit tests:

- validation functions
- money/decimal helpers
- domain error mapping
  Integration tests:
- HTTP endpoint tests
- real PostgreSQL via testcontainers-go
- migrations applied before tests

Minimum integration test coverage:

- POST /funds success
- POST /funds invalid status
- GET /funds/{id} not found
- POST /investors success
- POST /investors duplicate email conflict
- POST /funds/{fund_id}/investments success
- POST /funds/{fund_id}/investments unknown investor
- GET /funds/{fund_id}/investments success

⸻

6. Data Model

6.1 Migration: funds

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE funds (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name TEXT NOT NULL CHECK (length(trim(name)) > 0),
vintage_year INT NOT NULL CHECK (vintage_year BETWEEN 1900 AND 2100),
target_size_usd NUMERIC(18, 2) NOT NULL CHECK (target_size_usd > 0),
status TEXT NOT NULL CHECK (status IN ('Fundraising', 'Investing', 'Closed')),
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

6.2 Migration: investors

CREATE TABLE investors (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name TEXT NOT NULL CHECK (length(trim(name)) > 0),
investor_type TEXT NOT NULL CHECK (investor_type IN ('Individual', 'Institution', 'Family Office')),
email TEXT NOT NULL UNIQUE,
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

6.3 Migration: investments

CREATE TABLE investments (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
investor_id UUID NOT NULL REFERENCES investors(id) ON DELETE RESTRICT,
fund_id UUID NOT NULL REFERENCES funds(id) ON DELETE RESTRICT,
amount_usd NUMERIC(18, 2) NOT NULL CHECK (amount_usd > 0),
investment_date DATE NOT NULL
);

6.4 Indexes

CREATE INDEX idx_funds_created_at ON funds (created_at DESC);
CREATE INDEX idx_investors_created_at ON investors (created_at DESC);
CREATE INDEX idx_investments_fund_id_date ON investments (fund_id, investment_date DESC);
CREATE INDEX idx_investments_investor_id ON investments (investor_id);

6.5 No uniqueness on fund/investor investment pair

Do not add:

UNIQUE (fund_id, investor_id)

Reason:

The specification does not state that an investor can only invest once into a fund.
Allowing multiple investment records gives a more flexible commitment history model.

Document this in ADR.

⸻

7. Project Structure

Recommended final structure:

.
├── cmd/
│ └── api/
│ └── main.go
├── docs/
│ ├── adr/
│ │ ├── 0001-use-go-fiber-sqlc-goose.md
│ │ ├── 0002-use-postgres-numeric-decimal.md
│ │ ├── 0003-use-central-error-middleware.md
│ │ ├── 0004-use-testcontainers-for-integration-tests.md
│ │ └── 0005-scope-core-eight-endpoints.md
│ └── swagger/
├── internal/
│ ├── app/
│ │ └── app.go
│ ├── db/
│ │ ├── query/
│ │ │ ├── funds.sql
│ │ │ ├── investors.sql
│ │ │ └── investments.sql
│ │ └── sqlc/
│ ├── domain/
│ │ ├── entity/
│ │ ├── enum/
│ │ ├── error/
│ │ └── vo/
│ ├── port/
│ ├── usecase/
│ ├── adapter/
│ │ └── postgres/
│ ├── handler/
│ │ ├── funds.go
│ │ ├── investors.go
│ │ └── investments.go
│ ├── middleware/
│ │ ├── error.go
│ │ └── logger.go
│ ├── request/
│ │ └── bind.go
│ └── response/
│   └── models.go
├── migrations/
│ └── 001_init.sql
├── seed/
│ └── seed.sql
├── tests/
│ └── integration/
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── sqlc.yaml
├── go.mod
└── README.md

⸻

8. Common API Infrastructure

8.1 Generic request binding

Implement:

func BindAndValidate[T any](c *fiber.Ctx) (*T, error)

Responsibilities:

- Parse JSON request body
- Validate required fields
- Return domain error on invalid body

Use case:

req, err := request.BindAndValidate[CreateFundRequest](c)
if err != nil {
return err
}

⸻

8.2 Domain error

Implement:

type Kind string
const (
KindValidation Kind = "validation_error"
KindNotFound Kind = "not_found"
KindConflict Kind = "conflict"
KindInternal Kind = "internal_error"
)
type Error struct {
Kind Kind
Message string
Fields map[string]string
Err error
}

Constructors:

Validation(message string, fields map[string]string) *Error
NotFound(resource string) *Error
Conflict(message string) *Error
Internal(err error) *Error

⸻

8.3 Error middleware

Fiber config:

app := fiber.New(fiber.Config{
ErrorHandler: middleware.ErrorHandler,
})

Middleware behavior:

- domain error -> status + JSON error
- Unknown error -> 500 internal_error
- Hide internal DB details from client
- Log internal error server-side

⸻

8.4 Success response policy

Do not use generic success envelope.

Use raw spec response shapes:

[
{
"id": "...",
"name": "..."
}
]

not:

{
"data": []
}

Reason:

The API spec examples return raw objects and raw arrays.

⸻

9. sqlc Plan

9.1 sqlc.yaml

version: "2"
sql:

- engine: "postgresql"
  queries: "internal/db/query"
  schema: "migrations"
  gen:
  go:
  package: "db"
  out: "internal/db/sqlc"
  sql_package: "pgx/v5"
  emit_json_tags: true
  overrides: - db_type: "pg_catalog.numeric"
  go_type:
  import: "github.com/shopspring/decimal"
  package: "decimal"
  type: "Decimal"

Fallback if pgx scanning becomes annoying:

Use pgtype.Numeric in generated sqlc models and map to decimal.Decimal in DTOs.

⸻

9.2 funds.sql

-- name: ListFunds :many
SELECT id, name, vintage_year, target_size_usd, status, created_at
FROM funds
ORDER BY created_at DESC;
-- name: CreateFund :one
INSERT INTO funds (name, vintage_year, target_size_usd, status)
VALUES ($1, $2, $3, $4)
RETURNING id, name, vintage_year, target_size_usd, status, created_at;
-- name: GetFund :one
SELECT id, name, vintage_year, target_size_usd, status, created_at
FROM funds
WHERE id = $1;
-- name: UpdateFund :one
UPDATE funds
SET
name = $2,
vintage_year = $3,
target_size_usd = $4,
status = $5
WHERE id = $1
RETURNING id, name, vintage_year, target_size_usd, status, created_at;
-- name: FundExists :one
SELECT EXISTS (
SELECT 1 FROM funds WHERE id = $1
);

⸻

9.3 investors.sql

-- name: ListInvestors :many
SELECT id, name, investor_type, email, created_at
FROM investors
ORDER BY created_at DESC;
-- name: CreateInvestor :one
INSERT INTO investors (name, investor_type, email)
VALUES ($1, $2, $3)
RETURNING id, name, investor_type, email, created_at;
-- name: GetInvestor :one
SELECT id, name, investor_type, email, created_at
FROM investors
WHERE id = $1;
-- name: InvestorExists :one
SELECT EXISTS (
SELECT 1 FROM investors WHERE id = $1
);

⸻

9.4 investments.sql

-- name: ListInvestmentsByFund :many
SELECT id, investor_id, fund_id, amount_usd, investment_date
FROM investments
WHERE fund_id = $1
ORDER BY investment_date DESC, id DESC;
-- name: CreateInvestment :one
INSERT INTO investments (investor_id, fund_id, amount_usd, investment_date)
VALUES ($1, $2, $3, $4)
RETURNING id, investor_id, fund_id, amount_usd, investment_date;

⸻

10. Testing Strategy

10.1 Unit tests

Target packages:

internal/domain/entity
internal/domain/enum
internal/domain/error
internal/domain/vo
internal/port
internal/usecase
internal/request

Test examples:

domain/money_test.go

- accepts positive decimal
- rejects zero
- rejects negative
- rejects more than 2 decimal places
  domain/fund_test.go
- validates Fundraising
- validates Investing
- validates Closed
- rejects invalid status
  domain/investor_test.go
- validates Individual
- validates Institution
- validates Family Office
- rejects invalid type
  domain/error/error_test.go
- validation error maps to 400
- not found maps to 404
- conflict maps to 409

⸻

10.2 Integration tests with testcontainers

Use:

go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres

Test setup:

1. Start PostgreSQL container
2. Run goose migrations
3. Start Fiber app with test DB URL
4. Send HTTP requests using httptest or app.Test()
5. Assert status code and JSON body
6. Terminate container

Recommended integration test file:

tests/integration/api_test.go

Minimum cases:

POST /funds

- 201 success
- 400 invalid status
- 400 invalid amount
  GET /funds/{id}
- 200 success
- 404 not found
- 400 invalid UUID
  POST /investors
- 201 success
- 409 duplicate email
  POST /funds/{fund_id}/investments
- 201 success
- 404 unknown fund
- 404 unknown investor
- 400 invalid amount
  GET /funds/{fund_id}/investments
- 200 success
- 404 unknown fund

⸻

11. Swagger/OpenAPI Plan

11.1 Tooling

Use:

go get github.com/gofiber/swagger
go get github.com/swaggo/swag/cmd/swag

Generate docs:

swag init -g cmd/api/main.go -o docs/swagger

Expose:

app.Get("/swagger/\*", swagger.HandlerDefault)

11.2 Swagger requirements

Document:

- All 8 endpoints
- Request DTOs
- Response DTOs
- Error response DTO
- UUID path params
- Example values

Swagger title:

Titanbay Private Markets API

Swagger version:

1.0.0

⸻

12. Docker Plan

12.1 Dockerfile

Multi-stage build:

builder:
golang image
download deps
build binary
runtime:
distroless or alpine
copy binary
run API

For take-home simplicity, Alpine is acceptable.

12.2 docker-compose.yml

Services:

services:
postgres:
image: postgres:16
environment:
POSTGRES_USER: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_DB: titanbay
ports: - "5432:5432"
api:
build: .
environment:
DATABASE_URL: postgres://postgres:postgres@postgres:5432/titanbay?sslmode=disable
PORT: 8080
ports: - "8080:8080"
depends_on: - postgres

Optional:

Run migrations at API startup or via Makefile.

Preferred for take-home:

Run migrations at API startup only if simple and reliable.
Otherwise document `make migrate-up`.

For easiest reviewer experience:

API startup can run goose migrations automatically.

Document this decision in ADR.

⸻

13. ADR Plan

Put ADRs under:

docs/adr/

Use simple format:

# ADR-000X: Title

## Status

Accepted

## Context

## Decision

## Consequences

⸻

ADR-0001: Use Go, Fiber, sqlc, pgx, and goose

# ADR-0001: Use Go, Fiber, sqlc, pgx, and goose

## Status

Accepted

## Context

The task asks for a backend REST API backed by PostgreSQL. Titanbay's backend stack is Go-based, and the service is small enough that explicit SQL and a lightweight HTTP framework are sufficient.

## Decision

Use Go with Fiber for HTTP routing, pgx for PostgreSQL access, sqlc for type-safe SQL generation, and goose for database migrations.

## Consequences

This keeps the implementation simple and explicit. SQL remains visible in the repository, while sqlc provides compile-time safety for query inputs and outputs. Fiber provides fast and concise routing, though it does not use net/http directly, so request cancellation needs deliberate handling.

⸻

ADR-0002: Use PostgreSQL NUMERIC and decimal.Decimal for money

# ADR-0002: Use PostgreSQL NUMERIC and decimal.Decimal for monetary values

## Status

Accepted

## Context

The API deals with fund target sizes and investment amounts. These are financial values and should not be persisted or manipulated using binary floating-point types.

## Decision

Store monetary values as PostgreSQL NUMERIC(18,2). Represent them in Go using decimal.Decimal where practical. Validate that amounts are positive and have at most two decimal places.

## Consequences

This avoids float64 precision issues and makes the financial intent explicit. It adds a small amount of conversion complexity at the API/database boundary, but that trade-off is appropriate for financial data.

⸻

ADR-0003: Use domain errors and Fiber error middleware

# ADR-0003: Use domain errors and Fiber error middleware

## Status

Accepted

## Context

The API needs consistent validation, not found, conflict, and internal error responses. Handling these ad hoc in every handler would create duplication and inconsistent response shapes.

## Decision

Introduce a domain error type with kind, message, optional field errors, and wrapped internal error. Configure Fiber with a central ErrorHandler that maps domain errors to JSON responses.

## Consequences

Handlers stay focused on request handling and business logic. Error responses are consistent. Internal error details are hidden from clients while still allowing server-side logging.

⸻

ADR-0004: Use testcontainers-go for integration tests

# ADR-0004: Use testcontainers-go for PostgreSQL integration tests

## Status

Accepted

## Context

The API relies on PostgreSQL constraints, foreign keys, NUMERIC fields, and migrations. Mocking the database would not verify the most important behavior.

## Decision

Use testcontainers-go to run PostgreSQL for integration tests. Apply goose migrations before running API tests.

## Consequences

Integration tests verify the real database schema and constraints. Tests are slower than pure unit tests and require Docker, but they provide stronger confidence for a database-backed take-home task.

⸻

ADR-0005: Implement only the core 8 endpoints

# ADR-0005: Implement the core 8 endpoints from the task brief

## Status

Accepted

## Context

The task brief asks for all 8 endpoints. The provided API page contains Funds, Investors, and Investments sections that together define 8 endpoints. The same page also contains additional Transactions endpoints.

## Decision

Implement Funds, Investors, and Investments. Treat Transactions as out of scope for the initial submission.

## Consequences

The implementation stays aligned with the stated 2-3 hour time expectation. The README explicitly documents this scope decision so reviewers can see that the transaction endpoints were considered.

⸻

ADR-0006: Use raw success responses and structured error responses

# ADR-0006: Use spec-compatible success responses and structured error responses

## Status

Accepted

## Context

The API spec examples return raw objects and arrays for successful responses. The spec does not define a detailed error response format.

## Decision

Return success responses in the exact raw shape shown by the spec. Use a consistent structured JSON format for errors.

## Consequences

The API remains compatible with the provided contract while still providing predictable error handling for clients.

⸻

14. README Outline

# Titanbay Private Markets API

## Overview

A Go REST API for managing private market funds, investors, and investments.

## Tech Stack

- Go
- Fiber
- PostgreSQL
- sqlc
- pgx
- goose
- shopspring/decimal
- Swagger/OpenAPI
- testcontainers-go

## Running locally

````bash
docker compose up --build

API:

http://localhost:8080

Swagger:

http://localhost:8080/swagger/index.html

Running tests

go test ./...

API endpoints

List the 8 endpoints.

Design decisions

Summarise ADRs.

Assumptions

* Transaction endpoints are out of scope.
* Multiple investments from the same investor into the same fund are allowed.
* Monetary values are stored as NUMERIC(18,2).

AI-assisted development

Explain how AI was used.

---
# 15. Implementation Sequence
Recommended order:
```text
1. Repo scaffold
2. Docker Compose + Postgres
3. Go module + Fiber app
4. Config loading
5. goose migration
6. sqlc config + SQL queries
7. DB connection with pgx
8. Domain error + error middleware
9. Generic request binding
10. Domain validators
11. Funds handlers
12. Investors handlers
13. Investments handlers
14. Swagger annotations
15. Seed SQL
16. Unit tests
17. Integration tests with testcontainers
18. README
19. ADR docs

Timebox priority:

P0:
  - 8 endpoints working
  - Postgres schema
  - Docker local run
  - README
P1:
  - Decimal handling
  - Error middleware
  - Swagger
  - Unit tests
P2:
  - testcontainers integration tests
  - ADR polish
  - seed data

⸻

16. Final Definition of Done

The submission is done when:

- `docker compose up --build` starts the API and Postgres
- Swagger UI is available
- Goose migrations create the schema
- All 8 core endpoints work
- Success responses match the spec shape
- Error responses are consistent
- Monetary values use NUMERIC/decimal, not float64 persistence
- Unit tests pass
- Integration tests run against PostgreSQL via testcontainers
- README explains setup, API, design decisions, assumptions, and AI usage
- ADRs document key architecture decisions

⸻

17. 최종 전략

네 방향 그대로 가면 됨.

Go/Fiber + Postgres + sqlc + goose
+ decimal money handling
+ Swagger
+ Docker Compose
+ unit tests
+ testcontainers integration tests
+ ADR

이렇게 제출하면 단순 CRUD 과제여도 “아 이 사람은 API contract, DB integrity, financial precision, testability, documentation까지 보는구나”라는 인상이 남음.
특히 Titanbay 같은 private markets infra 쪽에는 decimal 처리 + relational integrity + ADR 조합이 꽤 잘 맞아.
````
