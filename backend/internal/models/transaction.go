package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID          string          `gorm:"type:varchar(21);primary_key" json:"id"`
	SpaceID     uuid.UUID       `gorm:"type:uuid;not null" json:"space_id"`
	Type        string          `gorm:"type:varchar(20);not null;default:'expense'" json:"type"`
	Title       string          `gorm:"type:varchar(100)" json:"title"`
	Date        time.Time       `gorm:"not null;default:NOW()" json:"date"`
	Currency    string          `gorm:"type:varchar(3);not null;default:'TWD'" json:"currency"`
	TotalAmount decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"total_amount"`
	Note        string          `gorm:"type:text" json:"note"`
	AIStatus    *string         `gorm:"type:varchar(20);column:ai_status" json:"ai_status,omitempty"`
	AIError     string          `gorm:"type:text;column:ai_error" json:"ai_error,omitempty"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	Space   *Space              `gorm:"foreignKey:SpaceID" json:"space,omitempty"`
	Expense *TransactionExpense `gorm:"foreignKey:TransactionID;constraint:OnDelete:CASCADE" json:"expense,omitempty"`
	Debts   []TransactionDebt   `gorm:"foreignKey:TransactionID;constraint:OnDelete:CASCADE" json:"debts,omitempty"`
	Images  []Image             `gorm:"polymorphic:Entity;polymorphicValue:transaction" json:"images,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionExpense struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TransactionID string          `gorm:"type:varchar(21);not null;uniqueIndex" json:"transaction_id"`
	Category      string          `gorm:"type:varchar(50)" json:"category"`
	ExchangeRate  decimal.Decimal `gorm:"type:decimal(12,6);default:1" json:"exchange_rate"`
	BillingAmount decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"billing_amount"`
	HandlingFee   decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"handling_fee"`
	PaymentMethod string          `gorm:"type:varchar(50)" json:"payment_method"`
	LocationURL   string          `gorm:"type:varchar(500);not null;default:''" json:"location_url"`

	// Associations
	Items []TransactionExpenseItem `gorm:"foreignKey:ExpenseID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (TransactionExpense) TableName() string {
	return "transaction_expenses"
}

type TransactionExpenseItem struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ExpenseID uuid.UUID       `gorm:"type:uuid;not null" json:"expense_id"`
	Name      string          `gorm:"type:varchar(100);not null" json:"name"`
	UnitPrice decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"unit_price"`
	Quantity  decimal.Decimal `gorm:"type:decimal(8,2);not null;default:1" json:"quantity"`
	Discount  decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount"`
	Amount    decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"amount"`
}

func (TransactionExpenseItem) TableName() string {
	return "transaction_expense_items"
}
