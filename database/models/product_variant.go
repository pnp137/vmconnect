package models

import "time"

type ProductVariant struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	ProductID uint    `gorm:"not null;index;index:idx_product_variants_product_active_price,priority:1" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Name      string  `gorm:"size:100" json:"name"`
	Unit      string  `gorm:"size:50" json:"unit"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `gorm:"index:idx_product_variants_product_active_price,priority:3" json:"price"`
	MRP       float64 `json:"mrp"`
	MOQ       float64 `json:"moq"`
	Stock     int     `json:"stock"`
	IsActive  bool    `gorm:"default:true;index:idx_product_variants_product_active_price,priority:2" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
