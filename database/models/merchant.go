package models

import "time"

type Merchant struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OwnerID      uint      `gorm:"column:user_id;not null" json:"owner_id"`
	Owner        User      `gorm:"foreignKey:OwnerID;references:ID" json:"owner"`
	VendorCode   string    `gorm:"size:10;index" json:"vendor_code,omitempty"`
	BusinessName string    `json:"business_name"`
	ShopName     string    `json:"shop_name"`
	Address      string    `json:"address"`
	Pincode      string    `json:"pincode"`
	IsDeleted    bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
