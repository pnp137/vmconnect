package dto

// ===== Vendor Order State Transition DTOs =====

// ConfirmOrderRequest represents the request to confirm an order
type ConfirmOrderRequest struct {
	EstimatedDelivery string `json:"estimated_delivery,omitempty"` // ISO 8601 format
	Notes             string `json:"notes,omitempty" validate:"max=500"`
}

// ConfirmOrderResponse represents the response for confirming an order
type ConfirmOrderResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// GenerateInvoiceRequest represents the request to generate invoice
type GenerateInvoiceRequest struct {
	InvoiceNo   string  `json:"invoice_no" validate:"required,max=50"`
	InvoiceDate string  `json:"invoice_date" validate:"required"` // ISO 8601 format
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	TaxAmount   float64 `json:"tax_amount,omitempty" validate:"gte=0"`
	FileURL     string  `json:"file_url,omitempty" validate:"max=500"`
	Notes       string  `json:"notes,omitempty" validate:"max=500"`
}

// GenerateInvoiceResponse represents the response for generating invoice
type GenerateInvoiceResponse struct {
	OrderID   uint   `json:"order_id"`
	InvoiceID uint   `json:"invoice_id"`
	InvoiceNo string `json:"invoice_no"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// DispatchOrderRequest represents the request to dispatch an order
type DispatchOrderRequest struct {
	TrackingID       string `json:"tracking_id" validate:"required,max=100"`
	Carrier          string `json:"carrier" validate:"required,max=100"`
	DispatchDate     string `json:"dispatch_date,omitempty"`     // ISO 8601 format
	ExpectedDelivery string `json:"expected_delivery,omitempty"` // ISO 8601 format
	Packages         int    `json:"packages,omitempty" validate:"min=1"`
	Notes            string `json:"notes,omitempty" validate:"max=500"`
}

// DispatchOrderResponse represents the response for dispatching an order
type DispatchOrderResponse struct {
	OrderID    uint   `json:"order_id"`
	TrackingID string `json:"tracking_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

// CancelOrderRequest represents the request to cancel an order
type CancelOrderRequest struct {
	Reason         string `json:"reason" validate:"required,max=500"`
	RefundRequired bool   `json:"refund_required,omitempty"`
	Notes          string `json:"notes,omitempty" validate:"max=500"`
}

// CancelOrderResponse represents the response for cancelling an order
type CancelOrderResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
