package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	EntityID   string    `gorm:"type:varchar(50);not null" json:"entity_id"`
	EntityType string    `gorm:"type:varchar(50);not null" json:"entity_type"`
	FilePath   string    `gorm:"type:text;not null" json:"file_path"`
	BlurHash   string    `gorm:"type:varchar(100)" json:"blur_hash"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Image) TableName() string {
	return "images"
}
