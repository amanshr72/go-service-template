# Go-Service-Boilerplate

A Go microservice demonstrating CRUD, TDD, ports & adapters, GraphQL, JWT auth, and CI — built incrementally to learn backend patterns, not for production use.

## Stack

- **Language:** Go 1.22 (stdlib `net/http`, no framework)
- **DB drivers:** PostgreSQL (`lib/pq`), MongoDB (`mongo-driver`), in-memory storage (`map[int]*User`).
- **API:** REST + GraphQL (`graphql-go`)
- **Auth:** JWT (`golang-jwt/jwt/v5`)
- **Migrations:** Goose, split into DDL and DML
- **Docs:** Swagger (`swaggo`)
- **Testing:** `testing` + `testify`
- **CI:** GitHub Actions

## Architecture

Feature-first folder layout. Each domain folder contains its own model, port (interface), adapters (implementations), service, and handler — not split by technical layer.

```
go-crud2/
├── main.go                        # wires dependencies, starts server
├── cmd/
│   └── migrate/main.go            # single entrypoint for all DB migrations
├── internal/
│   ├── user/
│   │   ├── model.go               # User struct, input DTOs
│   │   ├── port.go                # Repository & Service interfaces (ports)
│   │   ├── repository_postgres.go # Postgres adapter
│   │   ├── repository_mongodb.go  # MongoDB adapter
│   │   ├── repository_inmemory.go # In-memory adapter (no DB needed)
│   │   ├── service.go             # business logic, depends on port not adapter
│   │   ├── handler.go             # REST adapter
│   │   ├── resolver.go            # GraphQL adapter
│   │   ├── routes.go              # route registration for this domain
│   │   └── *_test.go
│   ├── auth/
│   │   ├── model.go               # LoginInput, Claims, TokenResponse
│   │   ├── service.go             # sign/validate JWT
│   │   └── handler.go             # POST /auth/login
│   ├── middleware/
│   │   ├── chain.go               # composes middlewares
│   │   ├── requestid.go           # X-Request-ID injection
│   │   ├── logger.go              # request logging
│   │   ├── recovery.go            # panic recovery
│   │   └── auth.go                # JWT validation middleware
│   └── health/
│       └── handler.go             # /health, /health/ready, /metrics
├── migrations/
│   ├── ddl/                       # schema changes — env-gated
│   └── dml/                       # seed/data changes — safe to auto-run
├── docs/                          # swagger-generated (swag init)
└── .github/workflows/ci.yml       # lint, migrate, test, build
```

### Why this layout

Ports (interfaces) live in `port.go`; adapters (Postgres, MongoDB, in-memory, REST, GraphQL) implement or consume those ports. Swapping a DB or adding a new transport means adding a new adapter file — service and handler logic untouched.

## Setup

```cmd
git clone https://github.com/go-service-template/go-crud2.git
cd go-crud2
go mod download
copy .env.example .env
```

Edit `.env` with real values, then start Postgres (or skip and use `DB_ADAPTER=inmemory`):

```cmd
docker-compose up -d postgres
```

## Migrations

DDL and DML are split and run independently through a single Go entrypoint — never call `goose` directly in this repo.

```cmd
set APP_ENV=dev
go run cmd/migrate/main.go ddl up
go run cmd/migrate/main.go dml up
```

DDL is blocked from this tool when `APP_ENV=prod` — schema changes in prod require a human running migrations manually with elevated, audited DB credentials. Real enforcement is at the Postgres role level (the CI/app DB user has no `CREATE`/`ALTER`/`DROP` grants).

## Running

```cmd
go run main.go
```

Or with a specific adapter:

```cmd
set DB_ADAPTER=inmemory
go run main.go
```

Server starts on `:8080`.

## API

### Auth

```
POST /auth/login
```

```json
{ "email": "admin@test.com", "password": "password" }
```

Returns a JWT. All `/api/*` and `/graphql` routes require `Authorization: Bearer <token>`.

### REST

```
POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

### GraphQL

```
POST /graphql
```

```graphql
{ users { id name email is_active } }
query { user(id: 1) { id name } }
query { activeUsers(active: true) { id name } }
query { userCount }
mutation { createUser(name: "Aman", email: "aman@t.com") { id name } }
mutation { updateUser(id: 1, name: "New") { id name } }
mutation { deleteUser(id: 1) }
```

### Observability

```
GET /health         # liveness
GET /health/ready    # readiness (checks DB)
GET /metrics         # uptime, goroutines, memory, request count
```

### Swagger

```cmd
swag init
go run main.go
```

Visit `http://localhost:8080/swagger/index.html`.

## Testing

Every layer (repository, service, handler, resolver, middleware) is tested in isolation using mock implementations of the relevant interface — no real DB required for the test suite.

```cmd
go test ./...
go test ./... -v -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## CI

GitHub Actions runs on every push/PR to `main`: `go vet` → lint → DDL migration against a throwaway Postgres container → full test suite with coverage → build. See `.github/workflows/ci.yml`.

## Environment Variables

| Variable | Purpose |
|---|---|
| `DATABASE_URL` | Postgres connection string for the app |
| `GOOSE_DRIVER` | Driver for migrations (`postgres`) |
| `GOOSE_DBSTRING` | Connection string for migrations |
| `JWT_SECRET` | Secret used to sign/verify JWTs |
| `APP_ENV` | `dev` / `staging` / `prod` — gates DDL migrations |
| `DB_ADAPTER` | `postgres` / `mongodb` / `inmemory` |

## Notes

This is a learning project, not a production template. No rate limiting, no real password hashing, no production-grade secret management, no graceful shutdown handling. Each of those is a reasonable next exercise.
