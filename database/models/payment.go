package models

import "time"

// Payment represents a payment transaction for an order.
// This model supports ultra-flexible payment tracking for unorganized markets where:
// - Merchants can pay anytime, any mode, any amount (token/partial/full/credit/COD)
// - Multiple payment modes can be used for a single order
// - Payment verification is independent of order status
// - System tracks complete payment ledger with audit trail
//
// Key Design Principles:
// 1. Payment ≠ Order Status: Payments are tracked independently from order state
// 2. Ledger-Based: All payments (even failed ones) are recorded for complete history
// 3. Vendor Verification: Vendor confirms payment receipt separately from merchant submission
// 4. Flexible Modes: Cash, UPI, Card, Cheque, Bank Transfer, COD - all supported
type Payment struct {
	ID      uint64 `gorm:"primaryKey" json:"id"`
	OrderID uint   `gorm:"index;not null;index:idx_payments_order_deleted_created,priority:1;index:idx_payments_order_deleted_status,priority:1" json:"order_id"` // Foreign key to orders table

	// Actor Information: Who submitted the payment and when
	SubmittedBy uint       `gorm:"not null;comment:'User ID who submitted payment'" json:"submitted_by"` // References users.id
	ActorRole   *ActorRole `gorm:"type:int;not null" json:"actor_role"`                                  // merchant(3) or vendor(2) - who initiated payment

	// Payment Details
	Amount      float64      `gorm:"not null" json:"amount"`                // Amount in this payment transaction
	PaymentMode *PaymentMode `gorm:"type:int;not null" json:"payment_mode"` // 10=Cash, 20=Card, 30=UPI, 40=NetBanking, etc. (see payment_mode_enum.go)
	PaymentType string       `gorm:"size:50" json:"payment_type,omitempty"` // Descriptive: "token", "partial", "full", "cod_collection", "settlement"

	// Transaction Identifiers (Optional - depends on payment mode)
	UTR             string `gorm:"size:100" json:"utr,omitempty"`              // Unique Transaction Reference for UPI/Bank transfers
	ReferenceNumber string `gorm:"size:100" json:"reference_number,omitempty"` // Cheque number, PO number, or other reference
	TransactionID   string `gorm:"size:100" json:"transaction_id,omitempty"`   // Online payment gateway transaction ID

	// Verification Status
	Status     *PaymentStatus `gorm:"type:int;not null;index:idx_payments_order_deleted_status,priority:3" json:"status"` // 10=Pending, 20=Processing, 30=Completed, 40=Failed, 50=Cancelled (see payment_status_enum.go)
	Verified   bool           `gorm:"default:false" json:"verified"`                                                      // Vendor confirmation flag: true = vendor confirmed receipt
	VerifiedBy uint           `json:"verified_by,omitempty"`                                                              // User ID of vendor who verified payment
	VerifiedAt *time.Time     `json:"verified_at,omitempty"`                                                              // Timestamp when vendor verified

	// Additional Information
	Notes    string  `gorm:"type:text" json:"notes,omitempty"`    // Free-text notes about payment (e.g., "Partial payment for goods received")
	Metadata JSONMap `gorm:"type:json" json:"metadata,omitempty"` // Extra structured data for future extensibility

	// Timestamps
	PaidAt    time.Time `json:"paid_at"`                                                                                                                                      // When payment was actually made (can be backdated)
	CreatedAt time.Time `gorm:"index:idx_payments_order_deleted_created,priority:3" json:"created_at"`                                                                        // When this record was created in system
	UpdatedAt time.Time `json:"updated_at"`                                                                                                                                   // Last modification timestamp
	IsDeleted bool      `gorm:"default:false;index;index:idx_payments_order_deleted_created,priority:2;index:idx_payments_order_deleted_status,priority:2" json:"is_deleted"` // Soft delete flag

	// Relationships: Preload these for complete payment context
	Order           Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`                 // The order this payment belongs to
	SubmittedByUser User  `gorm:"foreignKey:SubmittedBy" json:"submitted_by_user,omitempty"` // User who submitted payment
}

// TableName specifies the table name for Payment model
func (Payment) TableName() string {
	return "payments"
}

// Payment Workflow Example:
// 1. Merchant submits payment (status=Pending, verified=false)
// 2. System auto-calculates order.PaidAmount and order.OutstandingAmount
// 3. Vendor views payment and marks as verified (verified=true, status=Completed)
// 4. System updates order.OrderPaid when all payments are verified
//
// Flexible Payment Examples:
// - Token Payment: amount=500, payment_type="token", order_total=5900 (still pending 5400)
// - COD: amount=5900, payment_mode=Cash, payment_type="cod"
// - Split Payment: Pay 2000 UPI + 2000 Card + 1900 Cash (3 separate Payment records)
// - Credit: Merchant receives goods first, pays later (order can be SHIPPED/DELIVERED with paid_amount=0)
