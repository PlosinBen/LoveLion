package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionDebt struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TransactionID string          `gorm:"type:varchar(21);not null" json:"transaction_id"`
	PayerName     string          `gorm:"type:varchar(50);not null" json:"payer_name"`
	PayeeName     string          `gorm:"type:varchar(50);not null" json:"payee_name"`
	Amount        decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"amount"`
	SettledAmount decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"settled_amount"`
	IsSpotPaid    bool            `gorm:"not null;default:false" json:"is_spot_paid"`
}

func (TransactionDebt) TableName() string {
	return "transaction_debts"
}
