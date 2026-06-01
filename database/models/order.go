package models

import "time"

// ----------------- 7. ORDERS -----------------
type Order struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	OrderNo    string `gorm:"uniqueIndex;size:100" json:"order_no"` // Unique order number
	MerchantID *uint  `gorm:"index" json:"merchant_id,omitempty"`
	VendorID   uint   `gorm:"not null;index" json:"vendor_id"`

	OrderFor       string   `gorm:"size:20;index" json:"order_for,omitempty"`
	CustomerName   string   `gorm:"size:255" json:"customer_name,omitempty"`
	CustomerMobile string   `gorm:"size:20" json:"customer_mobile,omitempty"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	GoogleMapsURL  string   `gorm:"size:500" json:"google_maps_url,omitempty"`

	BusinessName    string `gorm:"size:255" json:"business_name,omitempty"`
	DeliveryAddress string `gorm:"type:text" json:"delivery_address,omitempty"`

	// Order State
	Status          *OrderStatus `gorm:"type:int;not null;index" json:"status"`
	StatusUpdatedAt *time.Time   `json:"status_updated_at,omitempty"`
	StatusUpdatedBy uint         `json:"status_updated_by"` // User ID who last updated status

	// Financial Tracking
	TotalAmount       float64 `json:"total_amount"`
	PaidAmount        float64 `gorm:"default:0" json:"paid_amount"`        // Auto-calculated from payments
	OutstandingAmount float64 `gorm:"default:0" json:"outstanding_amount"` // TotalAmount - PaidAmount

	// Payment Tracking (Independent of order status)
	OrderPaid   bool       `gorm:"default:false;index" json:"order_paid"` // Vendor confirms received money
	OrderPaidAt *time.Time `json:"order_paid_at,omitempty"`
	OrderPaidBy uint       `json:"order_paid_by,omitempty"` // Vendor UserID who marked as paid

	// Invoice Tracking
	InvoiceGenerated   bool       `gorm:"default:false" json:"invoice_generated"`
	InvoiceGeneratedAt *time.Time `json:"invoice_generated_at,omitempty"`
	InvoiceNumber      string     `gorm:"size:100" json:"invoice_number,omitempty"`
	InvoiceID          string     `json:"invoice_id,omitempty"` // Legacy field

	// Shipping Tracking (Simple - no logistics)
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`   // When vendor sent goods
	DeliveredAt *time.Time `json:"delivered_at,omitempty"` // When merchant received goods

	// Metadata
	ItemCount          int     `gorm:"default:0" json:"item_count"`
	Notes              string  `gorm:"type:text" json:"notes,omitempty"`
	CancellationReason string  `gorm:"type:text" json:"cancellation_reason,omitempty"`
	Metadata           JSONMap `gorm:"type:json" json:"metadata,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Vendor        Vendor          `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	Merchant      Merchant        `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	OrderItems    []OrderItem     `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	Payments      []Payment       `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
	OrderActivity []OrderActivity `gorm:"foreignKey:OrderID" json:"order_activity,omitempty"`
}
