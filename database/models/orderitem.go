package models

// ----------------- 8. ORDER ITEMS -----------------
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	Order     Order   `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	SubTotal  float64 `json:"subtotal"`
	IsDeleted bool    `gorm:"default:false" json:"is_deleted"`
}
