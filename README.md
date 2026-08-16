# Go API Boilerplate

[![Go Report Card](https://goreportcard.com/badge/github.com/Joe5451/go-api-boilerplate)](https://goreportcard.com/report/github.com/Joe5451/go-api-boilerplate)

A foundational **Go API project layout** built with **Gin** (REST), **gRPC**, **PostgreSQL**, and **Swagger** documentation. This boilerplate provides a well-structured Clean Architecture foundation to help you start new Go projects without having to rethink the project directory structure from scratch.

## Purpose

This project serves as a **reference implementation** and **starting point** for Go API development. It demonstrates:

- Clean Architecture with clear separation of concerns
- Complete testing setup (unit tests + feature tests)
- Docker configuration for local development
- API documentation with Swagger
- Configuration management
- Error handling patterns

Use this boilerplate as a foundation when starting new Go projects, so you don't have to recreate the project structure and setup from scratch each time.

## Prerequisites

- Go 1.25.1
- Docker
  - For local PostgreSQL via `docker compose`
  - Required for feature tests (uses testcontainers-go)

## Configuration

This project reads configuration from environment variables. A local `.env` file is optional (gitignored) and will be loaded if present.

Create `.env` from the provided template:

```bash
cp .env.example .env
```

`.env.example`:

```dotenv
DEBUG=true

# postgres
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DBNAME=book
POSTGRES_SCHEMA=public

# grpc
GRPC_PORT=50051
```

## Project Layout

```text
.
├── cmd
│   ├── graphql (not yet implemented)
│   ├── grpc
│   └── rest
├── proto/
│   ├── book.proto
│   ├── book.pb.go
│   └── book_grpc.pb.go
├── docs/
├── internal/
│   ├── adapter/
│   │   ├── grpc/
│   │   │   ├── book/
│   │   │   ├── interceptors/
│   │   │   └── grpckit/
│   │   ├── repositories/
│   │   └── rest/
│   │       ├── book/
│   │       ├── middlewares/
│   │       └── restkit/
│   ├── application/
│   │   └── port/
│   │       ├── in/
│   │       └── out/
│   ├── config/
│   ├── constant/
│   ├── domain/
│   └── infra/
├── mocks/
├── test/
│   ├── grpc/
│   └── rest/
├── Dockerfile
├── docker-compose.yml
└── init.sql
```

| Directory/File | Description |
|----------------|-------------|
| `cmd/` | Application entrypoints (REST, gRPC, GraphQL); composition root in each `main.go` |
| `proto/` | Protocol Buffers definitions and generated Go code (`book.proto`, `book.pb.go`, `book_grpc.pb.go`) |
| `docs/` | Generated Swagger artifacts (swaggo) |
| `internal/adapter/` | Implementations of ports (REST/gRPC handlers, repositories) |
| `internal/application/` | Use cases (business flows) |
| `internal/application/port/` | Interfaces (in/out) for dependency inversion |
| `internal/config/` | Config loading and structs |
| `internal/constant/` | Shared error codes/constants |
| `internal/domain/` | Entities + domain rules (no dependencies on other layers) |
| `internal/infra/` | Infrastructure (Postgres pool) |
| `mocks/` | Generated mocks (go.uber.org/mock) |
| `test/` | Feature tests (testcontainers + httptest / bufconn) |
| `Dockerfile` | Container build for the app |
| `docker-compose.yml` | Postgres (and optional app) for local/dev |
| `init.sql` | DB schema init for Postgres containers |

This project follows Clean Architecture. Dependencies point inward:

- **Domain** (`internal/domain`): entities and domain errors (pure Go, no frameworks)
- **Application** (`internal/application` + `internal/application/port`): use cases depend on interfaces, not implementations
- **Adapters** (`internal/adapter`): HTTP/gRPC handlers and repository implementations that satisfy ports
- **Infra** (`internal/infra`): database connection setup (pgxpool)
- **Composition root** (`cmd/`): wires dependencies and starts each transport

This keeps business logic testable and decoupled from transport (HTTP, gRPC) and infrastructure (Postgres).

### Directory and file naming conventions

#### Directory Naming: `lowercase` only

Use lowercase for directory names. Avoid snake_case, CamelCase, and special separators.

#### File naming: `snake_case`

Use `snake_case` for Go source files to improve readability, especially when names combine multiple words.

Single-word files (e.g. `book.go`, `config.go`, `routes.go`) are fine as-is. Test files follow the same rule: `book_service_test.go`, `pg_book_repo_test.go`.

## Run

### 1. Start PostgreSQL (via Docker)

Set env vars (e.g. in `.env`) and start Postgres:

```bash
docker compose up -d postgres
```

The database schema is initialized from `init.sql`.

### 2. Run the REST API server

```bash
go run ./cmd/rest/main.go
```

Server will listen on:
- http://localhost:8080

Swagger UI:
- http://localhost:8080/swagger/index.html

### 3. Run the gRPC server

```bash
go run ./cmd/grpc
```

Server will listen on:
- localhost:50051 (override with `GRPC_PORT`)

When `DEBUG=true`, server reflection is enabled for tools like [grpcurl](https://github.com/fullstorydev/grpcurl):

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"title":"1984","author":"George Orwell"}' \
  localhost:50051 proto.BookService/CreateBook
grpcurl -plaintext -d '{"id":1}' localhost:50051 proto.BookService/GetBook
```

## API Endpoints

Base URL: http://localhost:8080

Books:
- `POST /books`
- `GET /books/:id`
- `GET /books?page=1&per_page=10`
- `PUT /books/:id`
- `DELETE /books/:id`

Example (create a book):

```bash
curl -i -X POST "http://localhost:8080/books" \
  -H "Content-Type: application/json" \
  -d '{"title":"1984","author":"George Orwell"}'
```

## Test

Unit tests (focused on internal layers):

```bash
go test ./internal/...
```

All tests:

```bash
go test ./... -count=1
```

Feature tests (requires Docker; uses [testcontainers](https://github.com/testcontainers/testcontainers-go)):

```bash
go test ./test/rest/... -v
go test ./test/grpc/... -v
```

## Swagger

Install `swag` and regenerate:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/rest/main.go
```

## Protobuf

Regenerate gRPC Go code after changing `proto/book.proto`:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/book.proto
```
