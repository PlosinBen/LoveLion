---
description: How to run tests for backend and frontend
---

# Running Tests

Uses Docker Compose to run tests in consistent environment.

## Backend Tests (Go)
// turbo
1. Run all Go tests:
```bash
docker compose exec backend go test ./...
```

// turbo
2. Run tests with verbose output:
```bash
docker compose exec backend go test -v ./...
```

// turbo
3. Run tests with coverage:
```bash
docker compose exec backend go test -cover ./...
```

4. Run tests for a specific package:
```bash
docker compose exec backend go test -v ./internal/handlers/...
```

## Frontend Tests (Nuxt/Vitest)
// turbo
5. Run frontend unit tests:
```bash
docker compose exec frontend npm run test
```

// turbo
6. Run frontend tests in watch mode:
```bash
docker compose exec frontend npm run test:watch
```

## Run Tests with Fresh Database
7. Run backend tests with a separate test database:
```bash
docker compose exec backend go test -v ./... -tags=integration
```
