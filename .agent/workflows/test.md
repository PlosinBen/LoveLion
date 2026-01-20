---
description: How to run tests for backend and frontend
---

# Running Tests

Uses Docker Compose for backend tests and browser automation for frontend E2E tests.

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

## Frontend Browser Tests (E2E)

Use `browser_subagent` to test all pages by simulating user behavior. Do NOT write JavaScript test code.

5. Test all pages in order:
   - `/` - Dashboard: verify cards display correctly
   - `/ledger` - Ledger list: verify transactions load
   - `/ledger/add` - Add transaction: fill form and submit
   - `/trips` - Trip list: verify trips display
   - `/trips/new` - Create trip: fill form and submit
   - `/trips/{id}` - Trip detail: verify content loads
   - `/trips/{id}/ledger` - Trip ledger: verify transactions
   - `/trips/{id}/stores` - Price comparison: verify stores list

6. Each page test should:
   - Navigate to the page
   - Verify key elements are visible
   - Test primary user actions
   - Confirm expected results

## Full Test Suite
// turbo
7. Run complete backend tests:
```bash
docker compose exec backend go test -v ./...
```

8. Then run browser tests for all pages (use browser_subagent)
