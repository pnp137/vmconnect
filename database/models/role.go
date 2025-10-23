package models

import "time"

// ----------------- ROLES -----------------
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;unique;not null" json:"name"` // e.g. "vendor", "merchant", "admin"
	Description string    `gorm:"size:255" json:"description"`
	IsDeleted   bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
