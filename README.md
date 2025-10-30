# Todo List API

A modular RESTful API for managing todo lists, built with Go, chi, GORM, Zap logging, and JWT authentication.

## Features
- CRUD endpoints for todo lists with request validation via `go-playground/validator`.
- PostgreSQL persistence using GORM with automatic migrations.
- Structured logging powered by Uber's Zap.
- JWT (HS256) authentication middleware for API routes.
- Clean architecture layering (handler → service → repository).
- Comprehensive unit tests using `stretchr/testify` and `DATA-DOG/go-sqlmock`.

## Getting Started

### Prerequisites
- Go 1.22+
- PostgreSQL 14+
- `golangci-lint` (optional, for linting)

### Environment Setup
1. Copy `.env.example` to `.env` and adjust values as needed.
2. Ensure PostgreSQL is running and accessible with the configured credentials.
3. Run database migrations automatically on application start (`gorm.AutoMigrate`).

### Useful Commands
```bash
make tidy   # sync dependencies
make fmt    # format code
make vet    # go vet
make lint   # golangci-lint
make build  # build binary
make run    # run the app locally
make test   # run tests with race detector and coverage
```

### Running
```bash
make run
```
The server listens on `http://localhost:8080`. Health endpoint: `GET /health`.

### Authentication
All `/api/v1` routes require a valid JWT signed with HS256 using the configured secret and issuer.

### API Overview
See [`openapi.yaml`](openapi.yaml) for the full API specification.

## Project Structure
```
├── cmd/app/main.go
├── database/
├── handler/
├── middleware/
├── model/
├── package/
│   ├── config/
│   ├── logger/
│   └── response/
├── repository/
├── router/
├── service/
└── test/
```

## Testing
Tests are located alongside their respective packages and under `test/` for shared fixtures. Run `make test` to execute with the race detector and coverage enabled.
