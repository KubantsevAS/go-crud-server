# go-crud-server

A small HTTP CRUD service written in Go. The main goal of this project is to demonstrate a **clean, layered architecture with explicit repositories**, simple dependency injection, and clear separation of concerns.

## High‑level architecture

### Request flow

```text
HTTP Request
    ↓
Handler (internal/<feature>/handler.go)
    ↓
Payload / DTO (internal/<feature>/payload.go)
    ↓
Service (internal/<feature>/service.go) [optional]
    ↓
Repository (internal/<feature>/repository.go)
    ↓
Database (pkg/db)
    ↓
Response builder (pkg/response)
```

---

- **`cmd/main.go`**: Application composition and HTTP server startup.
  - Wires **config**, **DB**, **event bus**, **repositories**, **services**, **handlers**, and **middlewares**.
- **`internal/*`**: Feature modules (`auth`, `link`, `stat`, `user`). Each feature keeps its own
  handler, payloads/DTOs, models and repositories.
- **`pkg/*`**: Cross‑cutting packages: database access, DI interfaces, JWT, middleware, request/response helpers, event bus.

## Repository‑centric design

The core of this project is the **repository layer**, which hides all persistence logic behind simple Go interfaces/structs. Business code never talks to GORM directly – it depends on repositories.

### Repositories per feature

- **Links** – `internal/link/repository.go`
  - `LinkRepository` owns all CRUD operations on `Link` entities (create, get by hash/id, update, delete, list with pagination, count).
  - Uses a `*db.Db` (GORM wrapper) injected at construction time: `NewLinkRepository(database *db.Db) *LinkRepository`.
- **Statistics** – `internal/stat/repository.go`
  - `StatRepository` encapsulates click aggregation logic.
  - Exposes behaviours like `AddClick(linkId uint)` and `GetAll(by string, from, to time.Time)` that return aggregated stats instead of raw rows.
- **Users** – `internal/user/repository.go`
  - `UserRepository` encapsulates user creation and lookup: `Create(*User)`, `GetByEmail(email string)`.

Each repository **only knows about the DB wrapper** (`pkg/db`) and feature models. Higher layers see repositories as simple Go types/interfaces and don’t need to know GORM details.

### Dependency inversion via interfaces

To keep services independent from concrete repository implementations, the project defines **small interfaces** in `pkg/di` that describe just the behaviour needed:

- **`pkg/di/interfaces.go`**
  - `IUserRepository` – used by the auth service for registration and login.
  - `IStatRepository` – used by the statistics service for click tracking.

Services depend on these interfaces, not on concrete structs.

Example: **Auth service** (`internal/auth/service.go`)

- `AuthService` has a field `UserRepository di.IUserRepository`.
- At runtime, `cmd/main.go` passes a concrete `*user.UserRepository` which implements `IUserRepository`.
- This allows testing `AuthService` with **mocks/fakes** instead of the real DB.

## Composition root (`cmd/main.go`)

`main` acts as the **composition root** where all dependencies are wired:

- Load configuration: `configs.LoadConfig()`.
- Create database: `db.NewDb(conf)`.
- Init infrastructure: HTTP router (`http.NewServeMux()`), `eventBus` (`pkg/event`).
- **Create repositories**: `NewLinkRepository`, `NewUserRepository`, `NewStatRepository`.
- **Create services**:
  - `AuthService` with `IUserRepository`.
  - `StatService` with `IStatRepository` and `EventBus`.
- **Register handlers**: `NewAuthHandler`, `NewLinkHandler`, `NewStatHandler` – each receives only the dependencies it needs (repositories, services, config, event bus).
- **Wrap with middleware**: CORS, logging, and common middleware via `pkg/middleware.Chain`.

This makes it easy to see the full object graph and to change implementations in a single place.

## Feature modules

Each feature module under `internal/` follows the same structure:

- **`handler.go`**: HTTP endpoints (routing, reading request, invoking services/repos, writing response).
- **`payload.go`**: Request DTOs and validation helpers (via `pkg/request`).
- **`model.go`**: GORM models used by the repositories.
- **`repository.go`**: DB interactions only – no HTTP or business logic.
- **`service.go`** (where present): Business rules and orchestration that do not belong in HTTP handlers.

This keeps the structure predictable and makes each layer easy to navigate.

## Supporting infrastructure

- **Database (`pkg/db`)**
  - GORM wrapper responsible for opening the DB connection using values from `configs`.
  - Shared between all repositories.
- **Event bus (`pkg/event`)**
  - Simple event dispatcher that lets services/handlers emit domain events (e.g. clicks) decoupled from consumers.
- **Middleware (`pkg/middleware`)**
  - Common HTTP middleware (CORS, logging, auth, common concerns) that can be combined with `Chain`.
- **Request/response helpers**
  - `pkg/request`: decoding, validation and generic request handling helpers.
  - `pkg/response`: utilities for shaping uniform JSON responses and HTTP status codes.
- **JWT (`pkg/jwt`)**
  - Token generation and validation for auth flows.

## How to run

- **Requirements**: Go (matching version in `go.mod`), running database compatible with the GORM config, Docker (optional).
- **Local run**:
  - Configure DB connection and JWT/settings in `configs` or env vars.
  - Run database migrations in `migrations/` (or via `docker-compose.yml`).
  - Start server:

    ```bash
    go run ./cmd
    ```

  - The server listens on port `8081` by default (see `cmd/main.go`).

## Testing

- Unit tests for auth and JWT:
  - `cmd/auth_test.go`
  - `internal/auth/handler_test.go`
  - `internal/auth/service_test.go`
  - `pkg/jwt/jwt_test.go`

Run all tests:

```bash
go test ./...
```

## When to use this template

This project is a good starting point if you want:

- A **simple but realistic HTTP service** with separate layers (handler → service → repository → DB).
- A clear example of **repository pattern** over GORM.
- Lightweight **dependency inversion** using small interfaces.
- A structure that is **easy to test and evolve** as the codebase grows.
