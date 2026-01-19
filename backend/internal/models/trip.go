package models

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID           string     `gorm:"type:varchar(21);primary_key" json:"id"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	Description  string     `gorm:"type:text" json:"description"`
	StartDate    *time.Time `gorm:"type:date" json:"start_date"`
	EndDate      *time.Time `gorm:"type:date" json:"end_date"`
	BaseCurrency string     `gorm:"type:varchar(3);default:'TWD'" json:"base_currency"`
	CreatedBy    uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	LedgerID     *uuid.UUID `gorm:"type:uuid" json:"ledger_id"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	Creator *User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Ledger  *Ledger      `gorm:"foreignKey:LedgerID" json:"ledger,omitempty"`
	Members []TripMember `gorm:"foreignKey:TripID" json:"members,omitempty"`
}

func (Trip) TableName() string {
	return "trips"
}

type TripMember struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TripID    string     `gorm:"type:varchar(21);not null" json:"trip_id"`
	UserID    *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Name      string     `gorm:"type:varchar(50);not null" json:"name"`
	IsOwner   bool       `gorm:"default:false" json:"is_owner"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Associations
	Trip *Trip `gorm:"foreignKey:TripID" json:"trip,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (TripMember) TableName() string {
	return "trip_members"
}
