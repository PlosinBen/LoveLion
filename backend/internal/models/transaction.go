package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID            string          `gorm:"type:varchar(21);primary_key" json:"id"`
	LedgerID      uuid.UUID       `gorm:"type:uuid;not null" json:"ledger_id"`
	Title         string          `gorm:"type:varchar(100)" json:"title"`
	Payer         string          `gorm:"type:varchar(50)" json:"payer"`
	Date          time.Time       `gorm:"not null;default:NOW()" json:"date"`
	Currency      string          `gorm:"type:varchar(3);not null;default:'TWD'" json:"currency"`
	TotalAmount   decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"total_amount"`
	ExchangeRate  decimal.Decimal `gorm:"type:decimal(12,6);default:1" json:"exchange_rate"`
	BillingAmount decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"billing_amount"`
	HandlingFee   decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"handling_fee"`
	Category      string          `gorm:"type:varchar(50)" json:"category"`
	PaymentMethod string          `gorm:"type:varchar(50)" json:"payment_method"`
	Note          string          `gorm:"type:text" json:"note"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	Ledger *Ledger            `gorm:"foreignKey:LedgerID" json:"ledger,omitempty"`
	Items  []TransactionItem  `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	Splits []TransactionSplit `gorm:"foreignKey:TransactionID" json:"splits,omitempty"`
	Images []Image            `gorm:"polymorphic:Entity;polymorphicValue:transaction" json:"images,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionItem struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TransactionID string          `gorm:"type:varchar(21);not null" json:"transaction_id"`
	Name          string          `gorm:"type:varchar(100);not null" json:"name"`
	UnitPrice     decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"unit_price"`
	Quantity      decimal.Decimal `gorm:"type:decimal(8,2);not null;default:1" json:"quantity"`
	Discount      decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount"`
	Amount        decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"amount"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (TransactionItem) TableName() string {
	return "transaction_items"
}
