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
