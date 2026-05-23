//go:build integration

package integration

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const defaultDSN = "postgres://postgres:postgres@postgres:5432/lovelion?sslmode=disable"

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer db.Close()

	// Truncate all application tables so tests start from a clean state.
	// schema_migrations is left intact — migrations are already applied.
	_, err = db.Exec(`TRUNCATE
		users,
		spaces,
		space_members,
		space_invites,
		transactions,
		transaction_expenses,
		transaction_expense_items,
		transaction_debts,
		comparison_stores,
		comparison_products,
		images,
		announcements,
		expense_templates,
		inv_members,
		inv_settlements,
		inv_member_transactions,
		inv_settlement_allocations,
		inv_capital_futures_statements,
		inv_capital_oversea_futures_statements,
		inv_capital_oversea_futures_currencies,
		inv_stock_statements,
		inv_stock_holdings,
		inv_stock_trades
	CASCADE`)
	if err != nil {
		log.Fatalf("truncate tables: %v", err)
	}

	os.Exit(m.Run())
}
