package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Ledger represents a "Space" in the application (e.g., Personal, Trip, Group)
type Ledger struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Type           string         `gorm:"type:varchar(20);not null;default:'personal'" json:"type"` // personal, trip, group
	BaseCurrency   string         `gorm:"type:varchar(3);default:'TWD'" json:"base_currency"`
	Currencies     datatypes.JSON `gorm:"type:jsonb;default:'[\"TWD\"]'" json:"currencies"`
	MemberNames    datatypes.JSON `gorm:"type:jsonb;default:'[]';column:members" json:"member_names"` // Deprecated but kept for compatibility
	Categories     datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"categories"`
	PaymentMethods datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"payment_methods"`
	StartDate      *time.Time     `gorm:"type:date" json:"start_date"`
	EndDate        *time.Time     `gorm:"type:date" json:"end_date"`
	CoverImage     string         `gorm:"-" json:"cover_image"`
	IsPinned       bool           `gorm:"default:false" json:"is_pinned"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction  `gorm:"foreignKey:LedgerID" json:"transactions,omitempty"`
	Members      []LedgerMember `gorm:"foreignKey:LedgerID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
	Invites      []LedgerInvite `gorm:"foreignKey:LedgerID;constraint:OnDelete:CASCADE" json:"invites,omitempty"`
	Images       []Image        `gorm:"polymorphic:Entity;polymorphicValue:space" json:"images,omitempty"`
}

// PopulateCoverImage sets CoverImage from the first associated image.
func (l *Ledger) PopulateCoverImage() {
	if len(l.Images) > 0 {
		l.CoverImage = l.Images[0].FilePath
	}
}

func (Ledger) TableName() string {
	return "ledgers"
}

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