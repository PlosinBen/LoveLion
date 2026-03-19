package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Space struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Type           string         `gorm:"type:varchar(20);not null;default:'personal'" json:"type"` // personal, trip, group
	BaseCurrency   string         `gorm:"type:varchar(3);default:'TWD'" json:"base_currency"`
	Currencies     datatypes.JSON `gorm:"type:jsonb;default:'[\"TWD\"]'" json:"currencies"`
	SplitMembers   datatypes.JSON `gorm:"type:jsonb;default:'[]';column:split_members" json:"split_members"`
	Categories     datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"categories"`
	PaymentMethods datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"payment_methods"`
	StartDate      *time.Time     `gorm:"type:date" json:"start_date"`
	EndDate        *time.Time     `gorm:"type:date" json:"end_date"`
	CoverImage     string         `gorm:"-" json:"cover_image"`
	IsPinned       bool           `gorm:"default:false" json:"is_pinned"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:SpaceID" json:"transactions,omitempty"`
	Members      []SpaceMember `gorm:"foreignKey:SpaceID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
	Invites      []SpaceInvite `gorm:"foreignKey:SpaceID;constraint:OnDelete:CASCADE" json:"invites,omitempty"`
	Images       []Image       `gorm:"polymorphic:Entity;polymorphicValue:space" json:"images,omitempty"`
}

// PopulateCoverImage sets CoverImage from the first associated image.
func (s *Space) PopulateCoverImage() {
	if len(s.Images) > 0 {
		s.CoverImage = s.Images[0].FilePath
	}
}

func (Space) TableName() string {
	return "spaces"
}

type SpaceMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SpaceID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_space_user" json:"space_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_space_user" json:"user_id"`
	Role      string    `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
	Alias     string    `gorm:"type:varchar(50)" json:"alias"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Space *Space `gorm:"foreignKey:SpaceID" json:"space,omitempty"`
}

func (SpaceMember) TableName() string {
	return "space_members"
}

type SpaceInvite struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SpaceID   uuid.UUID  `gorm:"type:uuid;not null" json:"space_id"`
	Token     string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"token"`
	IsOneTime bool       `gorm:"type:boolean;not null;default:true" json:"is_one_time"`
	MaxUses   int        `gorm:"type:integer;not null;default:1" json:"max_uses"`
	UseCount  int        `gorm:"type:integer;not null;default:0" json:"use_count"`
	ExpiresAt *time.Time `gorm:"type:timestamptz" json:"expires_at"`
	CreatedBy uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	Space   *Space `gorm:"foreignKey:SpaceID" json:"space,omitempty"`
	Creator *User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (SpaceInvite) TableName() string {
	return "space_invites"
}
