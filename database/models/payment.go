package models

import "time"

// Payment represents a payment transaction for an order
// Supports flexible payment tracking - merchants can pay anytime, any mode, any amount
type Payment struct {
	ID      uint64 `gorm:"primaryKey" json:"id"`
	OrderID uint   `gorm:"index;not null" json:"order_id"`

	// Who & When
	SubmittedBy uint       `gorm:"not null;comment:'User ID who submitted payment'" json:"submitted_by"`
	ActorRole   *ActorRole `gorm:"type:int;not null" json:"actor_role"` // merchant or vendor

	// Payment Details
	Amount      float64      `gorm:"not null" json:"amount"`
	PaymentMode *PaymentMode `gorm:"type:int;not null" json:"payment_mode"` // Cash, Card, UPI, NetBanking, etc.
	PaymentType string       `gorm:"size:50" json:"payment_type,omitempty"` // "token", "partial", "full", "cod_collection", "settlement"

	// Transaction Info (optional, depends on payment mode)
	UTR             string `gorm:"size:100" json:"utr,omitempty"`              // For UPI/Bank transfers
	ReferenceNumber string `gorm:"size:100" json:"reference_number,omitempty"` // For cheque/other modes
	TransactionID   string `gorm:"size:100" json:"transaction_id,omitempty"`   // For online payments

	// Verification
	Status     *PaymentStatus `gorm:"type:int;not null" json:"status"` // Pending, Processing, Completed, Failed, etc.
	Verified   bool           `gorm:"default:false" json:"verified"`   // Vendor confirmation flag
	VerifiedBy uint           `json:"verified_by,omitempty"`           // Vendor UserID who verified
	VerifiedAt *time.Time     `json:"verified_at,omitempty"`

	// Metadata
	Notes    string  `gorm:"type:text" json:"notes,omitempty"`
	Metadata JSONMap `gorm:"type:json" json:"metadata,omitempty"`

	// Timestamps
	PaidAt    time.Time `json:"paid_at"`    // When payment was made
	CreatedAt time.Time `json:"created_at"` // When record was created
	UpdatedAt time.Time `json:"updated_at"` // When record was last updated
	IsDeleted bool      `gorm:"default:false;index" json:"is_deleted"`

	// Relationships
	Order           Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	SubmittedByUser User  `gorm:"foreignKey:SubmittedBy" json:"submitted_by_user,omitempty"`
}

// TableName specifies the table name for Payment model
func (Payment) TableName() string {
	return "payments"
}
