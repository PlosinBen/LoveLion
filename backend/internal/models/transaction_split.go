package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionSplit struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TransactionID string          `gorm:"type:varchar(21);not null" json:"transaction_id"`
	MemberID      uuid.UUID       `gorm:"type:uuid;not null" json:"member_id"`
	Amount        decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"amount"`
	IsPayer       bool            `gorm:"default:false" json:"is_payer"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`

	// Associations
	Transaction *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	Member      *TripMember  `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

func (TransactionSplit) TableName() string {
	return "transaction_splits"
}
