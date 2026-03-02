---
name: api-conventions
description: API design conventions
---

# API Conventions Mandate

## URL & Methods
- **Resource Names**: Plural nouns (e.g., `/trips`, `/ledgers`).
- **Style**: kebab-case for multi-word (e.g., `/trip-members`).
- **Nesting**: Parent -> Child (e.g., `/ledgers/{id}/transactions`).
- **GET**: 200 OK.
- **POST**: 201 Created. Return resource directly.
- **PUT/PATCH/DELETE**: 200 OK.

## Request/Response Format
- **Success**: Return resource directly (No wrapper).
- **Error**: MUST be `{"error": "message"}`.
- **IDs**: 
  - NanoID: Publicly visible (trips, ledgers, transactions, invites).
  - UUID: Internal linkage (users, ledger_members, items).

## Handler Structure
- **Transactions**: Mandatory for multi-table ops.
- **Validation**: MUST use Gin binding (`required`, `min=1`).
- **Verification**: ALWAYS verify ownership/access BEFORE processing.
- **ID Gen**: Use `utils.NewShortID` for public entities.