## Project Overview
- Project Name: Golang Small App (RESTful API or CLI)
- Tech Stack:
    - Language: Go 1.22+
    - Framework: go-chi/chi — lightweight HTTP router and middleware
    - Database: PostgreSQL 14+ using gorm.io/gorm
    - Migrations: Auto migration via GORM model tags
    - Validation: go-playground/validator/v10
    - Logging: uber-go/zap
    - Authentication: JWT (HS256)
    - Linting: golangci-lint
    - Testing: built-in testing, httptest, stretchr/testify and DATA-DOG/go-sqlmock
- Goal: 

A small, modular RESTful API built with idiomatic Go, clean architecture principles, structured logging, and best security practices.

## Recommended Folder Structure
```
├── main.go                 # main entry point (application startup)
├── .env.example            # environment variable example file
├── database/               # database connection and auto-migration initialization
├── handler/                # HTTP route handlers (controllers)
├── middleware/             # custom middlewares (auth, logger, recovery, etc.)
├── model/                  # data models (GORM structs + DTOs)
├── repository/             # data access layer (PostgreSQL queries via GORM)
├── router/                 # route initialization using go-chi/chi
├── service/                # business logic (use cases, domain rules)
├── package/                # reusable packages (utilities, helpers)
├── test/                   # test helpers and fixtures
├── Makefile                # build/test automation
├── go.mod / go.sum
└── README.md / AGENTS.md
```
Note:
- Use dependency injection between layers (handler → service → repository).
- Database migrations are handled automatically via GORM’s AutoMigrate() using model tags.
- Keep business logic and I/O concerns separated for maintainability and testability.

## Build and Test Commands
Run all commands from the repository root.

### Without Makefile
``` bash
go fmt ./... && go vet ./...
go mod tidy
go build -trimpath -ldflags="-s -w" -o bin/app ./cmd/app
go run ./cmd/app
go test ./... -race -shuffle=on -cover
```

### With Makefile (recommended)
```make
.PHONY: tidy fmt vet lint build run test cover

tidy: ; go mod tidy
fmt: ; go fmt ./...
vet: ; go vet ./...
lint: ; golangci-lint run
build: ; go build -trimpath -ldflags="-s -w" -o bin/app ./cmd/app
run: ; go run ./cmd/app
test: ; go test ./... -race -shuffle=on -cover
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
```

## Code Style Guidelines
- Formatting: must pass go fmt.
- Static Analysis: run go vet and golangci-lint.
- Error Handling:
    - Always handle and wrap errors with context `(fmt Errorf("doing X: %w", err))`.
    - Avoid `panic` except for fatal startup errors.
- Logging:
    - Use `zap.Logger` for structured logging.
    - Never log sensitive data (tokens, passwords, etc.).
- Context Usage:
    - Pass `context.Context` through all layers (`handler`, `service`, `repository`).
- Dependency Boundaries:
    - `handler` → `service` → `repository` → `database`.
    - Avoid circular dependencies.
- Naming:
    - Use descriptive function and variable names.
    - Keep files modular (< ~500 LOC).

## Testing Instructions
- Testing Levels:
    - Unit Tests:
        - Implemented in:
            - `service/` → test core business logic (use cases).
            - `repository/` → test data operations using github.com/DATA-DOG/go-sqlmock
            - `handler/` → test HTTP request/response using httptest.
            
            **Example: repository test with go-sqlmock**
            ```go
            db, mock, _ := sqlmock.New()
            defer db.Close()

            gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
            repo := NewUserRepository(gormDB)

            mock.ExpectQuery(`SELECT \* FROM "users"`).
                WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Budi"))

            users, err := repo.FindAll(context.Background())
            require.NoError(t, err)
            require.Len(t, users, 1)
            ```
    - Integration Tests:    
        - Run actual DB and API interactions, verifying full flow.
    - E2E/API Tests:
        - Simulate real HTTP calls using `httptest` or Postman automation.
- Testing Conventions:
    - File name: `*_test.go`
    - Function name: `TestFunction_Scenario_Expected`
    - Use `t.Run()` for subtests.
- Assertions:
    - Use `stretchr/testify` for assertions and mocks.
- Coverage Goal:
    - Minimum 70% for `service`, `repository`, and `handler`.
- Race Detection:
    - Always run tests with the `-race` flag.
- Fixtures:
    - Store sample data in `test/fixtures` to ensure reproducible results.

## Security Considerations
- Environment & Secrets
    - Use `.env` or environment variables for configuration.
    - Never commit secrets to version control.
- Authentication
    - Use JWT (HS256) and validate claims, issuer, and expiry.
- Input Validation
    - Use `validator/v10` for request validation.
- Database
    - Use GORM’s ORM layer to prevent SQL injection.
    - Use `AutoMigrate()` for schema sync on startup.
- Timeouts
    - Define HTTP and DB timeouts and use `context.WithTimeout`.
- Logging
    - Avoid logging PII, passwords, or tokens.
- Transport Security
    - Enforce HTTPS/TLS in production.
- Dependencies
    - Run `govulncheck ./...` regularly for vulnerability checks.


## Environment Variables

Environment variables control the behavior of the app and database connections.
Below is a standard `.env` configuration example:
```bash
# Application
APP_NAME=app
PORT=8080
DEBUG=false

# Database Configuration
DB_NAME=database_test
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_TIMEZONE="Asia/Jakarta"
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=1
DB_MAX_IDLE_TIME=10
DB_MAX_LIFE_TIME=50
```

### Description
| Variable              | Description                                        |
| --------------------- | -------------------------------------------------- |
| **APP_NAME**          | Application name identifier.                       |
| **PORT**              | HTTP server port.                                  |
| **DEBUG**             | Enables detailed logging when set to `true`.       |
| **DB_NAME**           | PostgreSQL database name.                          |
| **DB_USERNAME**       | Database username.                                 |
| **DB_PASSWORD**       | Database password (keep secure).                   |
| **DB_HOST**           | Database host or IP address.                       |
| **DB_TIMEZONE**       | Database timezone (recommended: `"Asia/Jakarta"`). |
| **DB_MAX_IDLE_CONNS** | Max number of idle DB connections.                 |
| **DB_MAX_OPEN_CONNS** | Max number of open DB connections.                 |
| **DB_MAX_IDLE_TIME**  | Max time (seconds) a connection can remain idle.   |
| **DB_MAX_LIFE_TIME**  | Max lifetime (seconds) of a single connection.     |

⚠️ Security Note:
Do not commit .env files to Git. Use .env.example for safe templates and load environment variables with github.com/joho/godotenv.

## Commit Messages or Pull Request Guidelines
### Conventional Commits (required)

**Format:**
`<type>(scope): <subject>`

**Common types:**
`feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `ci`, `build`

**Examples:**
- `feat(api): add CRUD for users`
- `fix(repo): handle not found error`
- `refactor(service): improve transaction handling`
- `test(handler): add tests for /auth/login`
- `chore: update zap dependency to v1.27.0`

#### Pull Request Checklist
- Passes fmt, vet, lint, and test -race
- No secrets or credentials in diffs
- Documentation updated (if behavior/config changed)
- Includes test coverage for new features
- Includes motivation, design, and impact summary

### Security Gotchas (Common Pitfalls)
- SQL Injection → Always use GORM’s query builder.
- No Timeout → Always define context deadlines and server timeouts.
- JWT Misuse → Reject invalid/expired tokens; verify alg.
- Path Traversal → Validate file paths before file access.
- Insecure TLS → Never use InsecureSkipVerify: true.
- Error Leaks → Do not expose internal error messages to users.
- Goroutine Leaks → Ensure all goroutines cancel with context.
- Dependency Vulnerabilities → Regularly audit and patch dependencies.

### Recommended Tooling
- golangci-lint
```bash 
golangci-lint run
```
- govulncheck
```bash
govulncheck ./...
```
- Pre-commit Hook (optional)
```bash
cat > .git/hooks/pre-commit <<'SH'
#!/usr/bin/env bash
set -e
go fmt ./...
go vet ./...
golangci-lint run
go test ./... -race
SH
chmod +x .git/hooks/pre-commit
```

### Environment & Local Setup
1. Copy configs/.env.example → .env and fill local values.
2. Start PostgreSQL (via Docker or local instance).
3. Auto-migrate models using GORM’s AutoMigrate() on app startup.
4. Run the app:
    ``` bash
    make run
    ```
5. Access API:
    ``` bash
    http://localhost:8080/health
    ```

## Contact & Ownership
- Tech Owner: Project Maintainer
- Security Contact: Security/DevOps Team Email
- Report vulnerabilities privately through email instead of public issues.