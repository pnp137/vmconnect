package models

import "time"

type ProductVariant struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Name      string  `gorm:"size:100" json:"name"`
	Unit      string  `gorm:"size:50" json:"unit"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	MRP       float64 `json:"mrp"`
	MOQ       float64 `json:"moq"`
	Stock     int     `json:"stock"`
	IsActive  bool    `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
