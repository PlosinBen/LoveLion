---
name: testing-guidelines
description: Rules for when and how to run tests during development
---

# Testing Guidelines

Rules that AI should automatically follow during development.

## Backend Testing

When **modifying ANY backend code** (handlers, services, models, utils, etc.):

1. **Rule**: You MUST run unit tests to verify your changes.
2. **Policy**: If *any* test fails, you MUST fix it **immediately** before proceeding. Do not ignore test failures.
3. Ensure corresponding test file exists (`*_test.go`). If not, create one.

### Running Tests
Run all Go tests using Docker:
```bash
docker compose exec backend go test ./...
```


---

## Frontend Browser Testing

When **modifying frontend pages**:

1. Use `browser_subagent` to simulate user behavior on that page
2. **Do NOT write JavaScript test code**
3. Tests should cover the modified functionality flow

### Testing Principles
- Simulate real user actions (click, type, navigate)
- Verify UI state changes
- Verify API call results are reflected in the UI

### Example Test Flow
```
1. Open target page
2. Perform user actions (click buttons, fill forms)
3. Verify page displays correct results
4. Report test results
```

---

## Database Schema Changes

When **modifying database schema** (add/modify migrations):

1. Reset database
2. Run all API tests (which creates test data)

```bash
# 1. Drop and recreate database
docker compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS lovelion"
docker compose exec postgres psql -U postgres -c "CREATE DATABASE lovelion"

# 2. Run migrations
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" up

# 3. Run API tests (creates test data)
docker compose exec backend go test ./...
```

---

## Summary Table

| Scenario | Automatic Action |
|----------|-----------------|
| Add/modify Backend API | Run all Go tests |
| Modify frontend page | Browser test that page |
| Modify database schema | Reset DB + Run all API tests |
