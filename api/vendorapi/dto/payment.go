package dto

import "time"

// MarkOrderPaidRequest - Vendor marks order as paid (bulk payment verification)
type MarkOrderPaidRequest struct {
	Notes string `json:"notes,omitempty"` // "Received cash settlement", "Cleared via bank", etc.
}

// MarkOrderPaidResponse - Response after marking order as paid
type MarkOrderPaidResponse struct {
	OrderID          uint      `json:"order_id"`
	OrderPaid        bool      `json:"order_paid"`
	PaymentsVerified int       `json:"payments_verified"` // Number of payments verified
	TotalAmount      float64   `json:"total_amount"`      // Total amount verified
	MarkedAt         time.Time `json:"marked_at"`
	OrderStatus      string    `json:"order_status"` // Updated order status
	Message          string    `json:"message"`
}

// VerifyPaymentRequest - Vendor verifies individual payment (optional)
type VerifyPaymentRequest struct {
	Verified bool   `json:"verified"`        // true = accept, false = reject
	Notes    string `json:"notes,omitempty"` // Verification notes
}

// VerifyPaymentResponse - Response after verifying individual payment
type VerifyPaymentResponse struct {
	PaymentID     uint64    `json:"payment_id"`
	OrderID       uint      `json:"order_id"`
	Amount        float64   `json:"amount"`
	PaymentStatus string    `json:"payment_status"` // "completed" or "failed"
	Verified      bool      `json:"verified"`
	VerifiedAt    time.Time `json:"verified_at"`
	Message       string    `json:"message"`
}

// GetOrderPaymentsResponse - Vendor view of order payments
type GetOrderPaymentsResponse struct {
	OrderID           uint                 `json:"order_id"`
	MerchantName      string               `json:"merchant_name"`
	OrderStatus       string               `json:"order_status"`
	TotalAmount       float64              `json:"total_amount"`
	PaidAmount        float64              `json:"paid_amount"`
	OutstandingAmount float64              `json:"outstanding_amount"`
	OrderPaid         bool                 `json:"order_paid"`
	OrderPaidAt       *time.Time           `json:"order_paid_at,omitempty"`
	Payments          []PaymentDetail      `json:"payments"`
	Summary           PaymentStatusSummary `json:"summary"`
}

// PaymentDetail - Detailed payment information for vendor
type PaymentDetail struct {
	ID              uint64     `json:"id"`
	Amount          float64    `json:"amount"`
	PaymentMode     string     `json:"payment_mode"`
	PaymentType     string     `json:"payment_type,omitempty"`
	PaymentStatus   string     `json:"payment_status"`
	Verified        bool       `json:"verified"`
	UTR             string     `json:"utr,omitempty"`
	ReferenceNumber string     `json:"reference_number,omitempty"`
	TransactionID   string     `json:"transaction_id,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	SubmittedAt     time.Time  `json:"submitted_at"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
}

// PaymentStatusSummary - Summary of payment status
type PaymentStatusSummary struct {
	TotalSubmitted      float64 `json:"total_submitted"`      // Total payment submitted
	Verified            float64 `json:"verified"`             // Total verified
	PendingVerification float64 `json:"pending_verification"` // Total pending
	PaymentsCount       int     `json:"payments_count"`
}
