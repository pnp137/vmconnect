package models

import "time"

// ----------------- 4. VENDOR-MERCHANT MAPPING -----------------
type VendorMerchantMapping struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VendorID   uint      `gorm:"not null;index:idx_vendor_merchant_active,priority:1;index:idx_vendor_merchant_lookup,priority:1" json:"vendor_id"`
	MerchantID uint      `gorm:"not null;index:idx_vendor_merchant_active,priority:3;index:idx_vendor_merchant_lookup,priority:2" json:"merchant_id"`
	Status     string    `gorm:"type:enum('active','inactive','blocked')" json:"status"`
	IsDeleted  bool      `gorm:"default:false;index:idx_vendor_merchant_active,priority:2;index:idx_vendor_merchant_lookup,priority:3" json:"is_deleted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
