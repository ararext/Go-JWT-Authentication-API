# Go JWT Authentication API

A production-style JWT authentication REST API built in Go, following Clean Architecture principles. This project is built incrementally with a focus on testability, security, and idiomatic Go — not a copy-paste tutorial.

> **Status: In active development.** This README reflects progress through the foundational layers (config, database, security, validation). Auth endpoints, middleware, and Docker/CI are still being built.

## Tech Stack

| Concern | Choice |
|---|---|
| Language | Go 1.24+ |
| Web framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | MongoDB |
| Auth | JWT (access + refresh tokens) |
| Password hashing | bcrypt |
| Config | godotenv |
| Validation | validator/v10 |
| Logging | zap (structured JSON) |
| Testing | Go testing + testify |

## Architecture

The project follows Clean Architecture — each layer only knows about the layer directly beneath it, and only through interfaces:

```
Handler  →  Service  →  Repository  →  MongoDB
(HTTP)      (business     (interface,
             logic)        no HTTP knowledge)
```

- **Handlers** parse HTTP requests and format responses — no business logic.
- **Services** contain business rules (e.g. "reject signup if email already exists").
- **Repositories** are defined as interfaces and talk to MongoDB — this is what allows services to be unit tested with a fake repository instead of a real database.

```
cmd/server/         entrypoint, server wiring
internal/
  config/            typed environment configuration
  database/          MongoDB connection + index setup
  models/            database-shape structs (never exposed via API)
  dto/               request/response shapes with validation tags
  repository/        UserRepository interface + MongoDB implementation
  service/           business logic (Signup/Login implemented, unit tested)
  handler/           HTTP handlers (signup, login implemented)
  middleware/         JWT auth middleware (implemented)
  utils/             bcrypt hashing, JWT utils, validation helper, response helper
  routes/            route registration (/api/v1/auth/*)
  logger/            zap logger initialization
pkg/
tests/
docs/
```

## Progress

- [x] Project skeleton, typed config loading, structured logging
- [x] MongoDB connection with startup ping + unique email index
- [x] `User` model and `UserRepository` (interface + MongoDB implementation)
- [x] bcrypt password hashing utilities, unit tested
- [x] Request/response DTOs with `validator/v10` tags
- [x] Centralized struct validation helper with readable error messages
- [x] JWT access/refresh token generation and validation, unit tested (round-trip, expiry, tampering, wrong secret)
- [x] Auth service (signup / login business logic), unit tested against a mock repository (no real database needed)
- [x] HTTP handlers and routes (`/api/v1/auth/signup`, `/api/v1/auth/login`) — fully wired end-to-end against real MongoDB
- [x] JWT middleware and protected routes (`/api/v1/users/me`, refresh, logout scaffold)
- [ ] Rate limiting, CORS, secure headers, graceful shutdown
- [ ] Docker + docker-compose
- [ ] Swagger/OpenAPI docs
- [ ] GitHub Actions CI
- [ ] Postman collection

## Setup (current state)

### Prerequisites
- Go 1.24+
- Docker (for running MongoDB locally)

### 1. Clone and install dependencies

```bash
git clone https://github.com/ararext/Go-JWT-Authentication-API.git
cd Go-JWT-Authentication-API
go mod tidy
```

### 2. Start MongoDB

```bash
docker run -d --name mongo-dev --restart unless-stopped -p 27017:27017 mongo:7
```

### 3. Configure environment

Create a `.env` file in the project root:

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=jwt_auth
JWT_SECRET=change-this-to-something-random
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h
```

### 4. Run the server

```bash
go run cmd/server/main.go
```

### 5. Verify

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## API Endpoints

All routes are prefixed with `/api/v1`.

| Method | Endpoint | Description | Auth required |
|---|---|---|---|
| POST | `/auth/signup` | Register a new user, returns access + refresh tokens | No |
| POST | `/auth/login` | Authenticate, returns access + refresh tokens | No |
| POST | `/auth/refresh` | Exchange a valid refresh token for a new access token | No (requires refresh token in body) |
| POST | `/auth/logout` | Client-side logout (stateless for now — see Roadmap) | No |
| GET | `/users/me` | Return the currently authenticated user | Yes (`Authorization: Bearer <accessToken>`) |

**Example — Signup**

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Ararext","email":"ararext@example.com","password":"securepass123"}'
```

**Example — Login**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ararext@example.com","password":"securepass123"}'
```

**Example — Accessing a protected route**

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <accessToken>"
```

All error responses follow a consistent shape:

```json
{ "success": false, "message": "invalid credentials" }
```

## Testing

```bash
go test ./... -v
```

Currently covers:
- Password hashing (bcrypt round-trip, salting verification)
- Struct validation (valid input, weak password, invalid email)
- JWT generation and validation (valid round-trip, expired tokens, malformed tokens, wrong signing secret, tampered signatures)
- AuthService signup/login logic — tested against an in-memory mock repository, no real database required (successful signup, duplicate email rejection, successful login, wrong password, unknown user)
- JWT middleware — protected route rejects missing and malformed tokens (`tests/` package, using `httptest`)

## Environment Variables

| Variable | Description | Example |
|---|---|---|
| `PORT` | Server port | `8080` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DATABASE_NAME` | Database name | `jwt_auth` |
| `JWT_SECRET` | Secret used to sign JWTs — never commit a real value | — |
| `ACCESS_TOKEN_DURATION` | Access token lifetime | `15m` |
| `REFRESH_TOKEN_DURATION` | Refresh token lifetime | `168h` |

## Roadmap

Once the core auth API is complete, planned extensions include: email verification (OTP), password reset, role-based access control, refresh token rotation, Redis-based token blacklist, request rate limiting, Prometheus metrics, and CI/CD.

## License

Not yet decided.