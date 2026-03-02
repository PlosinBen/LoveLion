# LoveLion Architecture

## Tech Stack
- **Frontend**: Nuxt 4, TailwindCSS, Iconify
- **Backend**: Go (Gin, GORM, golang-migrate)
- **Database**: PostgreSQL

## ID Strategy (MANDATORY)
- **NanoID**: URL-exposed (trips, transactions, stores, invites)
- **UUID**: Internal (users, members, items, ledger_members)

## Data Schema
### Accounting & Sharing
- `ledgers`: Main books. Type: personal/trip.
- `ledger_members`: User-ledger relations. Roles: owner, member. Supports `alias`.
- `ledger_invites`: Shared tokens. One-time or multi-use.
- `transactions`: Records linked to ledgers.
- `transaction_splits/items`: Detailed splits/lines.

### Travel
- `trips`: Trip containers. Linked to a ledger.
- `trip_comparison_*`: Price comparison stores and products.

## Design Mandates
- **Mobile-First**: Bottom navigation.
- **Money**: `DECIMAL(10, 2)` for amounts, `DECIMAL(12, 6)` for rates.
- **Audit**: All tables MUST have `created_at` and `updated_at` (TIMESTAMPTZ).