package model

import "time"

type Message struct {
	ID        string    `gorm:"not null;unique index;primaryKey" json:"id"`
	To        string    `gorm:"size:255;" json:"to,omitempty"`
	Content   string    `gorm:"size:255;" json:"content,omitempty"`
	Status    string    `gorm:"size:50;" json:"status,omitempty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
