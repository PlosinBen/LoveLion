---
description: How to run database migrations
---

# Database Migrations

Uses `golang-migrate` for database schema migrations via Docker Compose.

## Reset Database (Drop + Migrate + Seed)

Use this for development when database schema changes.

// turbo
1. Stop backend (to release DB connections):
```bash
docker compose stop backend
```

// turbo
2. Drop existing database:
```bash
docker compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS lovelion"
```

// turbo
3. Create fresh database:
```bash
docker compose exec postgres psql -U postgres -c "CREATE DATABASE lovelion"
```

// turbo
4. Start backend:
```bash
docker compose start backend
```

// turbo
5. Run all migrations:
```bash
docker compose exec backend go run cmd/migrate/main.go
```

// turbo
6. Seed test data:
```bash
docker compose exec backend go run cmd/seed/main.go
```

Test account after seeding:
- **Username**: `dev`
- **Password**: `dev123`

---

## Create New Migration
// turbo
5. Create a new migration file pair:
```bash
docker compose exec backend migrate create -ext sql -dir /app/migrations -seq <migration_name>
```

## Run Migrations Only
// turbo
6. Apply all pending migrations (without reset):
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" up
```

## Rollback Last Migration
7. Rollback the last applied migration:
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" down 1
```

## Check Migration Status
// turbo
8. Check current migration version:
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" version
```

## Force Set Version (use with caution)
9. Force set database version (useful for fixing dirty state):
```bash
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" force <version>
```
