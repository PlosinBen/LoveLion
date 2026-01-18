# LoveLion Architecture Reference

## Tech Stack
- **Frontend**: Nuxt 4 (Vue.js) + Iconify
- **Backend**: Go (Gin + GORM + golang-migrate)
- **Database**: PostgreSQL

## ID Strategy (Hybrid)
- **NanoID (VARCHAR)**: URL-exposed entities (`trips`, `transactions`, `stores`)
- **UUID**: Internal entities (`users`, `members`, `items`)

## Database Tables

### Core
| Table | ID Type | Description |
|-------|---------|-------------|
| `users` | UUID | User accounts |
| `trips` | NanoID | Travel trips |
| `trip_members` | UUID | Trip-to-user relationships |

### Accounting
| Table | ID Type | Description |
|-------|---------|-------------|
| `ledgers` | UUID | Account books (personal or trip) |
| `transactions` | NanoID | Expense records |
| `transaction_splits` | UUID | Bill splitting details |
| `transaction_items` | UUID | Line items per transaction |

### Travel Features
| Table | ID Type | Description |
|-------|---------|-------------|
| `trip_comparison_stores` | NanoID | Stores for price comparison |
| `trip_comparison_products` | UUID | Products with prices |

## API Endpoints

### Authentication
- `POST /api/users/register`
- `POST /api/users/login`
- `GET /api/users/me`

### Trips
- `GET/POST /api/trips`
- `GET/PUT/DELETE /api/trips/{id}`

### Ledgers & Transactions
- `GET/POST /api/ledgers`
- `GET/PUT/DELETE /api/ledgers/{id}`
- `GET/POST /api/ledgers/{id}/transactions`
- `GET/PUT/DELETE /api/ledgers/{ledger_id}/transactions/{id}`

### Price Comparison
- `GET/POST /api/trips/{trip_id}/stores`
- `GET /api/trips/{trip_id}/products`
- `GET/POST /api/trips/{trip_id}/stores/{store_id}/products`
- `PUT/DELETE /api/trips/{trip_id}/stores/{store_id}/products/{id}`

## Frontend Routes

### Dashboard
- `/` - Main dashboard with Personal Bookkeeping & My Trips cards

### Bookkeeping
- `/ledger` - View personal ledger
- `/ledger/add` - Add transaction
- `/ledger/{transaction_id}` - Transaction detail
- `/ledger/{transaction_id}/edit` - Edit transaction
- `/ledger/settings` - Ledger configuration

### Trips
- `/trips` - Trip list
- `/trips/{id}` - Trip dashboard
- `/trips/{id}/ledger` - Trip's ledger
- `/trips/{id}/ledger/add` - Add trip transaction
- `/trips/{id}/ledger/{transaction_id}` - Trip transaction detail
- `/trips/{id}/ledger/{transaction_id}/edit` - Edit trip transaction
- `/trips/{id}/stores` - Price comparison
- `/trips/{id}/stores/{store_id}` - Store products

## Design Notes
- **Mobile-First**: Bottom navigation on mobile, sidebar on desktop
- **Money Format**: `DECIMAL(10, 2)` for amounts, `DECIMAL(12, 6)` for exchange rates
- **Timestamps**: All tables have `created_at` and `updated_at` (TIMESTAMPTZ)
