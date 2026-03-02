---
name: testing-guidelines
description: Mandatory testing rules
---

# Testing Mandates

## Backend Testing
- **Trigger**: Modifying ANY backend code (handlers, models, utils).
- **Rule**: Run `docker compose exec backend go test ./...`.
- **Policy**: FIX all failures IMMEDIATELY before proceeding.

## Frontend Testing
- **Trigger**: Modifying frontend pages/components.
- **Rule**: Use `browser_subagent` to verify user flow.
- **Policy**: Simulate actions (click, type, navigate) and verify UI state.

## Database Schema Changes
- **Trigger**: Adding/modifying migrations.
- **Rule**:
  1. Reset DB: `DROP DATABASE IF EXISTS lovelion; CREATE DATABASE lovelion;`
  2. Run Migrations: `docker compose exec backend migrate ... up`
  3. Run API Tests: `docker compose exec backend go test ./...`

## Verification Summary
| Scenario | Action |
|----------|--------|
| API Backend | `go test ./...` |
| UI Frontend | `browser_test` |
| Schema/DB | `reset_db + go test` |