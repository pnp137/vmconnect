package models

import "time"

type MerchantProductVisibility struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	VendorID     uint      `gorm:"not null;index:idx_visibility_vendor_merchant_deleted,priority:1;index:idx_visibility_vendor_merchant_product_deleted,priority:1" json:"vendor_id"`
	MerchantID   uint      `gorm:"not null;index:idx_visibility_vendor_merchant_deleted,priority:2;index:idx_visibility_vendor_merchant_product_deleted,priority:2" json:"merchant_id"`
	ProductID    uint      `gorm:"not null;index:idx_visibility_vendor_merchant_product_deleted,priority:3" json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Vendor       Vendor    `gorm:"foreignKey:VendorID;references:ID" json:"vendor,omitempty"`
	Merchant     Merchant  `gorm:"foreignKey:MerchantID;references:ID" json:"merchant,omitempty"`
	IsWishlisted bool      `gorm:"default:false" json:"is_wishlisted"` // merchant can mark item for wishlist
	IsDeleted    bool      `gorm:"default:false;index:idx_visibility_vendor_merchant_deleted,priority:3;index:idx_visibility_vendor_merchant_product_deleted,priority:4" json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
