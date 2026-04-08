package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseTemplateData struct {
	Title         string                `json:"title"`
	Category      string                `json:"category"`
	Currency      string                `json:"currency"`
	PaymentMethod string                `json:"payment_method"`
	LocationURL   string                `json:"location_url"`
	Note          string                `json:"note"`
	TotalAmount   decimal.Decimal       `json:"total_amount"`
	Items         []ExpenseTemplateItem `json:"items"`
	Debts         []ExpenseTemplateDebt `json:"debts"`
}

type ExpenseTemplateItem struct {
	Name      string          `json:"name"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Discount  decimal.Decimal `json:"discount"`
}

type ExpenseTemplateDebt struct {
	PayerName  string          `json:"payer_name"`
	PayeeName  string          `json:"payee_name"`
	Amount     decimal.Decimal `json:"amount"`
	IsSpotPaid bool            `json:"is_spot_paid"`
}

type ExpenseTemplate struct {
	ID        uuid.UUID           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SpaceID   uuid.UUID           `gorm:"type:uuid;not null" json:"space_id"`
	Name      string              `gorm:"type:varchar(100);not null" json:"name"`
	Data      ExpenseTemplateData `gorm:"type:jsonb;serializer:json" json:"data"`
	CreatedAt time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ExpenseTemplate) TableName() string {
	return "expense_templates"
}
