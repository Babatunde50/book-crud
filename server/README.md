# Server

A Go service that provides:
- A CRUD **Books API** backed by Postgres
- A **URL Cleanup/Redirection** endpoint

## Prerequisites
- Go 1.24+ (or your installed version)
- Docker (for local Postgres + API tests)
- `make`
- (Optional) [`swag`](https://github.com/swaggo/swag) for Swagger generation

## Run (dev)

```bash
# Spins up postgres:16 and runs the API with automigrate
make run/dev

# Or run with your own database:
make run ARGS='-http-port=4748 -db-dsn="postgres:postgres@localhost:5432/byfood?sslmode=disable" -db-automigrate=true'
```

## Configuration (flags)

- `-http-port` (default 4748)
- `-db-dsn` (e.g. postgres://postgres:postgres@localhost:5432/byfood?sslmode=disable)
- `-db-automigrate` (true|false)
- `-base-url` (optional; defaults to http://localhost:4748)

## Swagger

Swagger UI is served at:

```bash
http://localhost:4748/swagger/index.html
```

Regenerate docs after handler comment changes:

```bash
swag init --generalInfo cmd/api/main.go --output cmd/api/docs
```

## Migrations

```bash
# Create a new migration
make migrations/new name=create_books_table

# Apply all up migrations
make migrations/up

# Rollback all
make migrations/down

# Go to a specific version
make migrations/goto version=000001

# Force a version (recover from a bad state)
make migrations/force version=000001

# Show current version
make migrations/version
```

## Testing

```bash
# Run everything:
make test

# Business layer only:
make test/business

# API handlers only (spins a Postgres container via internal/docker in tests):
make test/api


# Coverage (API package – function summary):
make cover/api

# Open HTML coverage for API:
make cover/api/html

# Coverage for all packages:
make cover/all
```

## API

Base URL: http://localhost:4748

### Books

- `GET /books` — List all books
- `POST /books` — Create a book
- `GET /books/{id}` — Get a book by ID
- `PUT /books/{id}` — Update a book by ID
- `DELETE /books/{id}` — Delete a book by ID

#### Create

```bash
curl -X POST http://localhost:4748/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"Clean Code","author":"Robert C. Martin","year":2008}'
```

#### Get by ID

```bash
curl http://localhost:4748/books/<uuid>
```

#### Update (partial allowed)

```bash
curl -X PUT http://localhost:4748/books/<uuid> \
  -H 'Content-Type: application/json' \
  -d '{"author":"New Author"}'
```

#### Delete

```bash
curl -X DELETE http://localhost:4748/books/<uuid>
```

### URL Processor

- `POST /url/process` — Process a URL with operation in ["canonical","redirection","all"]

#### Example (all)

```bash
curl -X POST http://localhost:4748/url/process \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://BYFOOD.com/food-EXPeriences?query=abc/","operation":"all"}'
# => {"processed_url":"https://www.byfood.com/food-experiences"}
```

#### Example (canonical)

```bash
curl -X POST http://localhost:4748/url/process \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://BYFOOD.com/food-EXPeriences?query=abc/","operation":"canonical"}'
# => {"processed_url":"https://BYFOOD.com/food-EXPeriences"}
```

## Error Format

### 400/404/405/500

```json
{ "Error": "Human readable message" }
```

### 422 validation

```json
{
  "FieldErrors": { "title": "title is required" },
  "Errors": []
}
```

These shapes are used by the client to display field-level and form-level errors.

## Project Structure

```bash
cmd/api/                  # HTTP server, routes, handlers, swagger docs
  docs/                   # (generated) swagger artifacts
business/book/            # Book core (domain), interfaces, errors
business/book/bookdb/     # SQLX store implementation for Book
business/urlprocessor/    # Canonical/redirection logic
internal/database/        # DB connect + migrations (iofs)
internal/docker/          # Test helper to spin containers
internal/request/         # JSON decode helpers
internal/response/        # JSON encode + metrics response writer
internal/validator/       # validation struct and helpers
assets/migrations/        # SQL migrations
```

## Commit Convention

This repo uses Conventional Commits:

- `feat`: new feature
- `fix`: bug fix
- `docs`: documentation only
- `test`: tests only
- `refactor`: code change without behavior change
- `chore`: tooling, deps, CI, etc.

### Examples

```
feat(api): add url processing endpoint
fix(bookdb): wrap sql.ErrNoRows with domain ErrNotFound
docs(swagger): annotate create/update book handlers
test(api): add table tests for /books and /url/process
```