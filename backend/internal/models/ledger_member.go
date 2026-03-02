package models

import (
	"time"

	"github.com/google/uuid"
)

type LedgerMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LedgerID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_ledger_user" json:"ledger_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_ledger_user" json:"user_id"`
	Role      string    `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
	Alias     string    `gorm:"type:varchar(50)" json:"alias"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Ledger *Ledger `gorm:"foreignKey:LedgerID" json:"ledger,omitempty"`
}

func (LedgerMember) TableName() string {
	return "ledger_members"
}
