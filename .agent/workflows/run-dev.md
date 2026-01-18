---
description: How to run the development environment locally
---

# Run Development Environment

Uses Docker Compose for consistent development setup.

## Start All Services
// turbo
1. Start the entire development stack (backend, frontend, database):
```bash
docker compose up -d
```

## View Logs
// turbo
2. Follow logs for all services:
```bash
docker compose logs -f
```

3. Follow logs for a specific service:
```bash
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

## Stop Services
// turbo
4. Stop all services:
```bash
docker compose down
```

## Rebuild After Code Changes
5. Rebuild containers after dependency changes:
```bash
docker compose build --no-cache
docker compose up -d
```

## Access Database
// turbo
6. Connect to PostgreSQL shell:
```bash
docker compose exec postgres psql -U postgres -d lovelion
```
