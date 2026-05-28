package models

import "time"

// ----------------- 6. CATEGORIES -----------------
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VendorID  uint      `gorm:"not null;uniqueIndex:idx_vendor_name,priority:1;index:idx_categories_vendor_deleted_name,priority:1" json:"vendor_id"` // each vendor can define own categories
	Name      string    `gorm:"size:100;uniqueIndex:idx_vendor_name,priority:2;index:idx_categories_vendor_deleted_name,priority:3" json:"name"`
	IsDeleted bool      `gorm:"default:false;index:idx_categories_vendor_deleted_name,priority:2" json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
