package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringArray is a JSON-backed array of strings stored in the DB.
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			*a = []string{}
			return nil
		}
		var arr []string
		if err := json.Unmarshal(v, &arr); err == nil {
			*a = arr
			return nil
		}
		*a = []string{string(v)}
		return nil
	case string:
		if v == "" {
			*a = []string{}
			return nil
		}
		var arr []string
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			*a = arr
			return nil
		}
		*a = []string{v}
		return nil
	default:
		*a = nil
		return nil
	}
}

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// ----------------- 5. PRODUCTS -----------------
type Product struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	VendorID   uint   `gorm:"not null;index;index:idx_products_vendor_category,priority:1;index:idx_products_vendor_active,priority:1;index:idx_products_vendor_featured,priority:1;index:idx_products_vendor_created_at,priority:1;index:idx_products_vendor_updated_at,priority:1;index:idx_products_vendor_name,priority:1" json:"vendor_id"`
	Vendor     Vendor `gorm:"foreignKey:VendorID" json:"vendor"`
	CategoryID uint   `gorm:"not null;index;index:idx_products_vendor_category,priority:2" json:"category_id"`

	Name         string      `gorm:"size:255;not null;index:idx_products_vendor_name,priority:2" json:"name"`
	Description  string      `gorm:"type:text" json:"description"`
	ImageURLs    StringArray `gorm:"type:json" json:"image_urls"`
	ThumbnailURL string      `gorm:"size:500" json:"thumbnail_url,omitempty"`
	IsFeatured   bool        `gorm:"default:false;index:idx_products_vendor_featured,priority:2" json:"is_featured"`
	IsActive     bool        `gorm:"default:true;index:idx_products_vendor_active,priority:2" json:"is_active"`

	Category Category         `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`

	CreatedAt time.Time `gorm:"index:idx_products_vendor_created_at,priority:2" json:"created_at"`
	UpdatedAt time.Time `gorm:"index:idx_products_vendor_updated_at,priority:2" json:"updated_at"`
}
