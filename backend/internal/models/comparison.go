package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ComparisonStore struct {
	ID           string    `gorm:"type:varchar(21);primary_key" json:"id"`
	TripID       string    `gorm:"type:varchar(21);not null" json:"trip_id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	GoogleMapURL string    `gorm:"type:text" json:"google_map_url"`
	Location     string    `gorm:"type:text" json:"location"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	Trip     *Trip               `gorm:"foreignKey:TripID" json:"trip,omitempty"`
	Products []ComparisonProduct `gorm:"foreignKey:StoreID" json:"products,omitempty"`
}

func (ComparisonStore) TableName() string {
	return "trip_comparison_stores"
}

type ComparisonProduct struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	StoreID   string          `gorm:"type:varchar(21);not null" json:"store_id"`
	Name      string          `gorm:"type:varchar(100);not null" json:"name"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`
	Currency  string          `gorm:"type:varchar(3);default:'TWD'" json:"currency"`
	Unit      string          `gorm:"type:varchar(20)" json:"unit"`
	Note      string          `gorm:"type:text" json:"note"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	Store *ComparisonStore `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (ComparisonProduct) TableName() string {
	return "trip_comparison_products"
}
