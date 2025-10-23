package models

import (
	"time"
)

// ----------------- 5. PRODUCTS -----------------
type Product struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	VendorID      uint       `gorm:"not null" json:"vendor_id"`
	Vendor        Vendor     `gorm:"foreignKey:VendorID" json:"vendor"`
	CategoryID    uint       `gorm:"not null" json:"category_id"` // ✅ linked to Categories table
	Name          string     `gorm:"size:255" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	Category      string     `gorm:"size:100" json:"category"`
	SKU           int        `gorm:"uniqueIndex" json:"sku"`
	Unit          string     `gorm:"size:50" json:"unit"` // e.g. "box", "kg", "piece"
	Price         float64    `json:"price"`
	Stock         int        `json:"stock"`
	ImageURL      string     `gorm:"size:500" json:"image_url"` // S3 URL
	ThumbnailURL  string     `gorm:"size:500" json:"thumbnail"` // optional
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	IsFeatured    bool       `gorm:"default:false" json:"is_featured"`
	IsDeleted     bool       `gorm:"default:false" json:"is_deleted"`
	LastRestocked *time.Time `json:"last_restocked"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
