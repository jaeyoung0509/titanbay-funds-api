# Titanbay Private Markets API

Go REST API for managing private market funds, investors, and investor commitments.

## Stack

- Go 1.26+
- Fiber v3
- PostgreSQL
- PostgreSQL access via `database/sql` using the pgx driver
- Goose migrations
- shopspring/decimal for monetary precision
- zerolog for structured request/response and error logs
- OpenAPI/Swagger UI
- testcontainers-go for integration tests

## Run locally

```bash
docker compose up --build
```

The API starts on `http://localhost:8080`.
Swagger UI is available at `http://localhost:8080/swagger/index.html`.

## Useful commands

```bash
make test
make migrate-up
make seed
```

## API scope

Only the 8 core endpoints are implemented:

- `GET /funds`
- `POST /funds`
- `PUT /funds`
- `GET /funds/{id}`
- `GET /investors`
- `POST /investors`
- `GET /funds/{fund_id}/investments`
- `POST /funds/{fund_id}/investments`

The transaction-related endpoints shown in the HTML spec are intentionally out of scope for this submission.

See [Architecture Decision Records](docs/adr/README.md).

## Design decisions

- PostgreSQL enforces integrity through primary keys, foreign keys, unique constraints, and CHECK constraints.
- Monetary values are stored as `NUMERIC(18,2)` and represented in Go with `decimal.Decimal`-backed custom types.
- The API returns raw JSON objects and arrays, matching the spec examples, rather than wrapping successful responses.
- The code uses a small layered structure: domain types and validation, usecase services for application flow, and HTTP/PostgreSQL adapters for IO.
- Error responses are driven by domain errors and rendered through a central Fiber error handler.
- Request/response logs are emitted as structured zerolog events and are correlated with `X-Request-ID`.

```json
{
  "error": {
    "code": "validation_error",
    "message": "validation failed",
    "fields": {
      "status": "must be one of Fundraising, Investing, Closed"
    }
  }
}
```

## Assumptions

- `vintage_year` is validated to the range `1900-2100`.
- Multiple investment records from the same investor into the same fund are allowed.
- `investment_date` is date-only in JSON.
- Money fields are returned as JSON numbers, not strings.

## AI usage

AI tools were used to accelerate development, especially for scaffolding, README structure, edge-case review, and validation or error-handling design.

The final implementation decisions were reviewed manually, including:

- limiting the scope to the 8 core endpoints
- using PostgreSQL constraints for data integrity
- using decimal-backed money types instead of `float64`
- returning raw success responses to match the provided API spec
- using centralized error handling and request logging

AI-generated suggestions were treated as drafts and adjusted to match the API specification and implementation constraints.
