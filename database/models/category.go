package models

import "time"

// ----------------- 6. CATEGORIES -----------------
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VendorID  uint      `gorm:"not null" json:"vendor_id"` // each vendor can define own categories
	Name      string    `gorm:"size:100;uniqueIndex:idx_vendor_name,priority:2" json:"name"`
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
