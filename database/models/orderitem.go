package models

// ----------------- 8. ORDER ITEMS -----------------
type OrderItem struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	OrderID          uint    `gorm:"not null;index" json:"order_id"`
	Order            Order   `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	ProductID        uint    `gorm:"not null;index" json:"product_id"`
	Product          Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	ProductVariantID uint    `gorm:"index" json:"product_variant_id,omitempty"`
	ProductName      string  `gorm:"size:255" json:"product_name,omitempty"`
	VariantName      string  `gorm:"size:255" json:"variant_name,omitempty"`
	Quantity         float64 `json:"quantity"`
	UnitPrice        float64 `json:"unit_price,omitempty"`
	TotalPrice       float64 `json:"total_price,omitempty"`
	Price            float64 `json:"price"`
	SubTotal         float64 `json:"subtotal"`
	IsDeleted        bool    `gorm:"default:false" json:"is_deleted"`
}
