package models

import "time"

type Vendor struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID;references:ID" json:"user"` // Ensure User is defined in your models
	CompanyName string    `json:"company_name"`
	GSTNumber   string    `json:"gst_number"`
	Address     string    `json:"address"`
	LogoURL     string    `json:"logo_url"` // S3 link for vendor logo
	IsDeleted   bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
