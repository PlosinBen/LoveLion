---
trigger: always_on
---

# Security & Access Control Mandates

## Data Integrity
- **Access Control**: Every mutation (PUT/POST/DELETE) MUST verify `LedgerMember` or `TripMember` access.
- **Ownership**: Only `owner` can modify ledger settings, manage members, or revoke invites.
- **Leaking**: NEVER return `PasswordHash` or secrets in JSON responses.
- **Transactions**: Multi-table updates MUST use database transactions.

## Verification Patterns
- **Ledger Access**: `h.verifyLedgerAccess(ledgerID, userID)`
- **Trip Access**: `h.verifyTripAccess(tripID, userID)`
- **Invites**: Token MUST be valid (not expired, within use limits).

## ID Privacy
- **Public IDs**: USE NanoID for URL-exposed entities to prevent ID enumeration.
- **Internal IDs**: USE UUID for internal relationships.
