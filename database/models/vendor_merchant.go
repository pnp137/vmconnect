package models

import "time"

// ----------------- 4. VENDOR-MERCHANT MAPPING -----------------
type VendorMerchantMapping struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VendorID   uint      `gorm:"not null" json:"vendor_id"`
	MerchantID uint      `gorm:"not null" json:"merchant_id"`
	Status     string    `gorm:"type:enum('active','inactive','blocked')" json:"status"`
	IsDeleted  bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
