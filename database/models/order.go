package models

import "time"

// ----------------- 7. ORDERS -----------------
type Order struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	MerchantID  uint         `gorm:"not null" json:"merchant_id"`
	VendorID    uint         `gorm:"not null" json:"vendor_id"`
	Status      *OrderStatus `gorm:"type:int;not_null;index:idx_order_status;"`
	TotalAmount float64      `json:"total_amount"`
	Notes       string       `gorm:"type:text" json:"notes,omitempty"`

	InvoiceID    string          `json:"invoice_id"`
	IsDeleted    bool            `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	OrderItems   []OrderItem     `gorm:"foreignKey:OrderID" json:"order_items"`
	OrderHistory []OrderActivity `gorm:"foreignKey:OrderID" json:"orderHistory,omitempty"`

	Vendor   Vendor   `gorm:"foreignKey:VendorID" json:"vendor"`
	Merchant Merchant `gorm:"foreignKey:MerchantID" json:"merchant"`
}
