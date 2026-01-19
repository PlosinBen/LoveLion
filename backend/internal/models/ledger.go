package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Ledger struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Type           string         `gorm:"type:varchar(20);not null;default:'personal'" json:"type"`
	Currencies     datatypes.JSON `gorm:"type:jsonb;default:'[\"TWD\"]'" json:"currencies"`
	Members        datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"members"`
	Categories     datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"categories"`
	PaymentMethods datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"payment_methods"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:LedgerID" json:"transactions,omitempty"`
}

func (Ledger) TableName() string {
	return "ledgers"
}
