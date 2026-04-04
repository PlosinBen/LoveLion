package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func d(v string) decimal.Decimal {
	d, _ := decimal.NewFromString(v)
	return d
}

func TestBuildExpenseItems_SingleItem(t *testing.T) {
	expenseID := uuid.New()
	items, total := buildExpenseItems(expenseID, []ExpenseItemInput{
		{Name: "Ramen", UnitPrice: d("150"), Quantity: d("1")},
	})

	assert.Len(t, items, 1)
	assert.Equal(t, "Ramen", items[0].Name)
	assert.True(t, d("150").Equal(items[0].Amount))
	assert.True(t, d("150").Equal(total))
}

func TestBuildExpenseItems_MultipleItems(t *testing.T) {
	expenseID := uuid.New()
	items, total := buildExpenseItems(expenseID, []ExpenseItemInput{
		{Name: "A", UnitPrice: d("100"), Quantity: d("2")},
		{Name: "B", UnitPrice: d("50"), Quantity: d("1")},
	})

	assert.Len(t, items, 2)
	assert.True(t, d("200").Equal(items[0].Amount)) // 100 * 2
	assert.True(t, d("50").Equal(items[1].Amount))   // 50 * 1
	assert.True(t, d("250").Equal(total))
}

func TestBuildExpenseItems_ZeroQuantityDefaultsToOne(t *testing.T) {
	expenseID := uuid.New()
	items, total := buildExpenseItems(expenseID, []ExpenseItemInput{
		{Name: "Item", UnitPrice: d("300"), Quantity: d("0")},
	})

	assert.Len(t, items, 1)
	assert.True(t, d("1").Equal(items[0].Quantity))
	assert.True(t, d("300").Equal(items[0].Amount))
	assert.True(t, d("300").Equal(total))
}

func TestBuildExpenseItems_WithDiscount(t *testing.T) {
	expenseID := uuid.New()
	items, total := buildExpenseItems(expenseID, []ExpenseItemInput{
		{Name: "Item", UnitPrice: d("100"), Quantity: d("3"), Discount: d("10")},
	})

	// amount = (100 - 10) * 3 = 270
	assert.True(t, d("270").Equal(items[0].Amount))
	assert.True(t, d("270").Equal(total))
}

func TestBuildExpenseItems_Empty(t *testing.T) {
	items, total := buildExpenseItems(uuid.New(), nil)
	assert.Nil(t, items)
	assert.True(t, decimal.Zero.Equal(total))
}

func TestBuildExpenseItems_SetsExpenseID(t *testing.T) {
	expenseID := uuid.New()
	items, _ := buildExpenseItems(expenseID, []ExpenseItemInput{
		{Name: "A", UnitPrice: d("10"), Quantity: d("1")},
		{Name: "B", UnitPrice: d("20"), Quantity: d("1")},
	})

	for _, item := range items {
		assert.Equal(t, expenseID, item.ExpenseID)
		assert.NotEqual(t, uuid.Nil, item.ID)
	}
}

// --- calcSettledAmount tests ---

func TestCalcSettledAmount_SpotPaid(t *testing.T) {
	debt := DebtInput{Amount: d("500"), IsSpotPaid: true}
	expense := ExpenseInput{}
	result := calcSettledAmount(debt, d("1000"), expense, "TWD")
	assert.True(t, decimal.Zero.Equal(result))
}

func TestCalcSettledAmount_BaseCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		billing  string
		want     string
	}{
		{"TWD currency", "TWD", "0", "300"},
		{"zero billing implies base", "JPY", "0", "300"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debt := DebtInput{Amount: d("300"), IsSpotPaid: false}
			expense := ExpenseInput{BillingAmount: d(tt.billing)}
			result := calcSettledAmount(debt, d("1000"), expense, tt.currency)
			assert.True(t, d(tt.want).Equal(result), "got %s", result)
		})
	}
}

func TestCalcSettledAmount_ForeignCurrencyWithBilling(t *testing.T) {
	tests := []struct {
		name        string
		debtAmount  string
		totalAmount string
		billing     string
		want        string
	}{
		{
			name:        "proportional ceiling: 300/1000 * 4500 = 1350",
			debtAmount:  "300",
			totalAmount: "1000",
			billing:     "4500",
			want:        "1350",
		},
		{
			name:        "ceiling rounds up: 333/1000 * 4500 = 1498.5 -> 1499",
			debtAmount:  "333",
			totalAmount: "1000",
			billing:     "4500",
			want:        "1499",
		},
		{
			name:        "exact division: 500/1000 * 4000 = 2000",
			debtAmount:  "500",
			totalAmount: "1000",
			billing:     "4000",
			want:        "2000",
		},
		{
			name:        "small fraction: 1/3 * 300 = 100",
			debtAmount:  "1",
			totalAmount: "3",
			billing:     "300",
			want:        "100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debt := DebtInput{Amount: d(tt.debtAmount), IsSpotPaid: false}
			expense := ExpenseInput{BillingAmount: d(tt.billing)}
			result := calcSettledAmount(debt, d(tt.totalAmount), expense, "JPY")
			assert.True(t, d(tt.want).Equal(result), "got %s, want %s", result, tt.want)
		})
	}
}

func TestCalcSettledAmount_ZeroTotal(t *testing.T) {
	debt := DebtInput{Amount: d("100"), IsSpotPaid: false}
	expense := ExpenseInput{BillingAmount: d("500")}
	// totalAmount is zero — should return zero (division guard)
	result := calcSettledAmount(debt, d("0"), expense, "JPY")
	assert.True(t, decimal.Zero.Equal(result))
}

// --- buildDebts tests ---

func TestBuildDebts_BaseCurrency(t *testing.T) {
	expense := ExpenseInput{}
	debts := buildDebts("txn-1", []DebtInput{
		{PayerName: "Bob", PayeeName: "Alice", Amount: d("500")},
		{PayerName: "Carol", PayeeName: "Alice", Amount: d("300")},
	}, d("800"), &expense, "TWD")

	assert.Len(t, debts, 2)
	assert.Equal(t, "Bob", debts[0].PayerName)
	assert.Equal(t, "Alice", debts[0].PayeeName)
	assert.True(t, d("500").Equal(debts[0].Amount))
	assert.True(t, d("500").Equal(debts[0].SettledAmount)) // base currency: settled = amount
	assert.Equal(t, "txn-1", debts[0].TransactionID)

	assert.True(t, d("300").Equal(debts[1].Amount))
	assert.True(t, d("300").Equal(debts[1].SettledAmount))
}

func TestBuildDebts_SpotPaid(t *testing.T) {
	expense := ExpenseInput{}
	debts := buildDebts("txn-1", []DebtInput{
		{PayerName: "Bob", PayeeName: "Alice", Amount: d("500"), IsSpotPaid: true},
	}, d("500"), &expense, "TWD")

	assert.Len(t, debts, 1)
	assert.True(t, debts[0].IsSpotPaid)
	assert.True(t, decimal.Zero.Equal(debts[0].SettledAmount))
}

func TestBuildDebts_ForeignCurrency(t *testing.T) {
	expense := ExpenseInput{BillingAmount: d("4500")}
	debts := buildDebts("txn-1", []DebtInput{
		{PayerName: "Bob", PayeeName: "Alice", Amount: d("300")},
		{PayerName: "Carol", PayeeName: "Alice", Amount: d("700")},
	}, d("1000"), &expense, "JPY")

	assert.Len(t, debts, 2)
	// Bob: 300/1000 * 4500 = 1350
	assert.True(t, d("1350").Equal(debts[0].SettledAmount), "got %s", debts[0].SettledAmount)
	// Carol: 700/1000 * 4500 = 3150
	assert.True(t, d("3150").Equal(debts[1].SettledAmount), "got %s", debts[1].SettledAmount)
}

func TestBuildDebts_PaymentNoExpense(t *testing.T) {
	// For payments, expense is nil — settled = amount
	debts := buildDebts("txn-1", []DebtInput{
		{PayerName: "Bob", PayeeName: "Alice", Amount: d("1000")},
	}, d("1000"), nil, "TWD")

	assert.Len(t, debts, 1)
	assert.True(t, d("1000").Equal(debts[0].Amount))
	assert.True(t, d("1000").Equal(debts[0].SettledAmount))
}

func TestBuildDebts_Empty(t *testing.T) {
	expense := ExpenseInput{}
	debts := buildDebts("txn-1", nil, d("100"), &expense, "TWD")
	assert.Nil(t, debts)
}

func TestBuildDebts_UniqueIDs(t *testing.T) {
	expense := ExpenseInput{}
	debts := buildDebts("txn-1", []DebtInput{
		{PayerName: "A", PayeeName: "B", Amount: d("100")},
		{PayerName: "C", PayeeName: "B", Amount: d("200")},
	}, d("300"), &expense, "TWD")

	ids := map[uuid.UUID]bool{}
	for _, debt := range debts {
		assert.NotEqual(t, uuid.Nil, debt.ID)
		assert.False(t, ids[debt.ID], "duplicate debt ID")
		ids[debt.ID] = true
	}
}
