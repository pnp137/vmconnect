package models

import "time"

type Merchant struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	VendorCode string    `gorm:"size:10;index" json:"vendor_code,omitempty"`
	ShopName   string    `json:"shop_name"`
	Address    string    `json:"address"`
	Pincode    string    `json:"pincode"`
	IsDeleted  bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
