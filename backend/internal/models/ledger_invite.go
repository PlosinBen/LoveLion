package models

import (
	"time"

	"github.com/google/uuid"
)

type LedgerInvite struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LedgerID  uuid.UUID  `gorm:"type:uuid;not null" json:"ledger_id"`
	Token     string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"token"`
	IsOneTime bool       `gorm:"type:boolean;not null;default:true" json:"is_one_time"`
	MaxUses   int        `gorm:"type:integer;not null;default:1" json:"max_uses"`
	UseCount  int        `gorm:"type:integer;not null;default:0" json:"use_count"`
	ExpiresAt *time.Time `gorm:"type:timestamptz" json:"expires_at"`
	CreatedBy uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	Ledger  *Ledger `gorm:"foreignKey:LedgerID" json:"ledger,omitempty"`
	Creator *User   `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (LedgerInvite) TableName() string {
	return "ledger_invites"
}
