---
description: How to run database migrations
---

# Database Migrations

Uses `golang-migrate` for database schema migrations via Docker Compose.

## Create New Migration
// turbo
1. Create a new migration file pair (run from host machine):
```bash
docker compose exec backend migrate create -ext sql -dir /app/migrations -seq <migration_name>
```

Or if `migrate` is installed locally:
```bash
cd backend
migrate create -ext sql -dir migrations -seq <migration_name>
```

## Run Migrations
// turbo
2. Apply all pending migrations:
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" up
```

## Rollback Last Migration
3. Rollback the last applied migration:
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" down 1
```

## Check Migration Status
// turbo
4. Check current migration version:
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" version
```

## Force Set Version (use with caution)
5. Force set database version (useful for fixing dirty state):
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" force <version>
```

## Reset Database (Drop + Migrate + Seed via Tests)

Use this when database schema changes require a fresh start.

// turbo
6. Drop existing database:
```bash
docker compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS lovelion"
```

// turbo
7. Create fresh database:
```bash
docker compose exec postgres psql -U postgres -c "CREATE DATABASE lovelion"
```

// turbo
8. Run all migrations:
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" up
```

// turbo
9. Run API tests to seed data:
```bash
docker compose exec backend go test ./...
```
