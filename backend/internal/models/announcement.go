package models

import "time"

type Announcement struct {
	ID             string     `gorm:"type:varchar(20);primaryKey" json:"id"`
	Title          string     `gorm:"type:varchar(255);not null" json:"title"`
	Content        string     `gorm:"type:text;not null;default:''" json:"content"`
	Status         string     `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	BroadcastStart *time.Time `gorm:"type:timestamptz" json:"broadcast_start"`
	BroadcastEnd   *time.Time `gorm:"type:timestamptz" json:"broadcast_end"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
