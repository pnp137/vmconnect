package models

import "time"

type Invoice struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OrderIDs    []uint         `gorm:"type:json" json:"order_ids"`
	VendorID    uint           `gorm:"not null" json:"vendor_id"`
	MerchantID  uint           `gorm:"not null" json:"merchant_id"`
	InvoiceNo   string         `gorm:"size:50;uniqueIndex;not null" json:"invoice_no"`
	InvoiceDate time.Time      `gorm:"not null" json:"invoice_date"`
	TotalAmount float64        `gorm:"not null" json:"total_amount"`
	TaxAmount   float64        `gorm:"default:0" json:"tax_amount"`
	PaidAmount  float64        `gorm:"default:0" json:"paid_amount"`
	Status      *InvoiceStatus `gorm:"size:50;default:10" json:"status"`
	IsDeleted   bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
