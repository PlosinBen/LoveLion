---
description: How to build production artifacts
---

# Build for Production

Uses Docker for production builds.

## Build All Production Images
// turbo
1. Build production Docker images:
```bash
docker compose -f docker-compose.prod.yml build
```

## Build Backend Only
// turbo
2. Build Go production binary:
```bash
docker compose exec backend go build -o bin/server main.go
```

Or build production backend image:
```bash
docker build -t lovelion-backend:latest ./backend
```

## Build Frontend Only
// turbo
3. Build Nuxt production bundle:
```bash
docker compose exec frontend npm run build
```

Or build production frontend image:
```bash
docker build -t lovelion-frontend:latest ./frontend
```

## Preview Frontend Production Build
// turbo
4. Preview the production build locally:
```bash
docker compose exec frontend npm run preview
```

## Export Static Frontend (if applicable)
5. Generate static site:
```bash
docker compose exec frontend npm run generate
```
