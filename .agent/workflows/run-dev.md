---
description: How to run the development environment locally
---

# Run Development Environment

Uses Docker Compose with volume mounts. No build required for development.

## Start All Services
// turbo
1. Start the entire development stack:
```bash
docker compose up -d
```

First run will take longer as it downloads dependencies.

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

## Restart Backend (after code changes)
// turbo
4. Restart backend to apply Go code changes:
```bash
docker compose restart backend
```

## Stop Services
// turbo
5. Stop all services:
```bash
docker compose down
```

## Access Database
// turbo
6. Connect to PostgreSQL shell:
```bash
docker compose exec postgres psql -U postgres -d lovelion
```
