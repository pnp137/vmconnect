package dto

import "time"

// SubmitPaymentRequest - Merchant submits a payment for an order
type SubmitPaymentRequest struct {
	Amount          float64                `json:"amount" validate:"required,gt=0"`
	PaymentMode     int                    `json:"payment_mode" validate:"required"` // Enum value: 10=cash, 20=card, 30=upi, etc.
	PaymentType     string                 `json:"payment_type,omitempty"`           // "token", "partial", "full", "cod_collection", "settlement"
	UTR             string                 `json:"utr,omitempty"`                    // For UPI/Bank transfers
	ReferenceNumber string                 `json:"reference_number,omitempty"`       // For cheque/other modes
	TransactionID   string                 `json:"transaction_id,omitempty"`         // For online payments
	Notes           string                 `json:"notes,omitempty"`                  // Additional notes
	Metadata        map[string]interface{} `json:"metadata,omitempty"`               // Additional metadata
}

// SubmitPaymentResponse - Response after submitting a payment
type SubmitPaymentResponse struct {
	PaymentID     uint64    `json:"payment_id"`
	OrderID       uint      `json:"order_id"`
	Amount        float64   `json:"amount"`
	PaymentMode   string    `json:"payment_mode"`
	PaymentStatus string    `json:"payment_status"`
	TotalPaid     float64   `json:"total_paid"`   // Total amount paid so far
	Outstanding   float64   `json:"outstanding"`  // Remaining amount to be paid
	OrderStatus   string    `json:"order_status"` // Current order status
	OrderPaid     bool      `json:"order_paid"`   // Vendor confirmation flag
	SubmittedAt   time.Time `json:"submitted_at"`
	Message       string    `json:"message"`
}

// GetPaymentsResponse - Response for fetching all payments of an order
type GetPaymentsResponse struct {
	OrderID           uint               `json:"order_id"`
	OrderStatus       string             `json:"order_status"`
	TotalAmount       float64            `json:"total_amount"`
	PaidAmount        float64            `json:"paid_amount"`
	OutstandingAmount float64            `json:"outstanding_amount"`
	OrderPaid         bool               `json:"order_paid"` // Vendor confirmation
	OrderPaidAt       *time.Time         `json:"order_paid_at,omitempty"`
	Payments          []PaymentSummary   `json:"payments"`
	PaymentSummary    PaymentModeSummary `json:"payment_summary"`
}

// PaymentSummary - Individual payment details
type PaymentSummary struct {
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
	PaidAt          time.Time  `json:"paid_at"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
}

// PaymentModeSummary - Summary of payments by mode
type PaymentModeSummary struct {
	ByMode              map[string]float64 `json:"by_mode"`              // Amount by payment mode
	Verified            float64            `json:"verified"`             // Total verified amount
	PendingVerification float64            `json:"pending_verification"` // Total pending verification
}
