---
description: How to run database migrations
---

# Database Migrations

## Reset Database (Drop + Migrate + Seed)

// turbo
```bash
./bin/refresh-database
```

Test account: `dev` / `dev123`

---

## Run Migrations Only

// turbo
```bash
./bin/migrate
```

---

## Create New Migration

// turbo
```bash
docker compose exec backend migrate create -ext sql -dir /app/migrations -seq <migration_name>
```

Creates `backend/migrations/XXXXXX_<name>.up.sql` and `.down.sql`
