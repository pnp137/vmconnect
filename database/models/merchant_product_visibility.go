package models

import "time"

type MerchantProductVisibility struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	VendorID     uint      `gorm:"not null" json:"vendor_id"`
	MerchantID   uint      `gorm:"not null" json:"merchant_id"`
	ProductID    uint      `gorm:"not null" json:"product_id"`
	IsWishlisted bool      `gorm:"default:false" json:"is_wishlisted"` // merchant can mark item for wishlist
	IsDeleted    bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
