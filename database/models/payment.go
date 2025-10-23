package models

import "time"

// ----------------- 9. PAYMENTS -----------------
type Payment struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	InvoiceID     uint64         `gorm:"index;not null"`
	Amount        float64        `gorm:"not null"`
	PaymentMode   string         `gorm:"type:varchar(50)"`
	TransactionID string         `gorm:"unique"`
	Status        *PaymentStatus `gorm:"type:int;not_null;index:idx_payment_status;"`
	PaidAt        *time.Time
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Invoice       Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice"`
}
