package dto

import "time"

// CartItemRequest represents the request to manage cart items (add/update only)
type CartItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"min=1"` // Must be greater than 0
}

// CartItemResponse represents the response for cart item operations
type CartItemResponse struct {
	OrderID   uint   `json:"order_id"`
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Action    string `json:"action"` // "added", "updated", "removed"
	Message   string `json:"message"`
}

// OrderInfo represents order information
type OrderInfo struct {
	ID          uint      `json:"id"`
	OrderNumber string    `json:"order_number"`
	VendorID    uint      `json:"vendor_id"`
	VendorName  string    `json:"vendor_name"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	ItemCount   int       `json:"item_count"` // Changed from ItemCount to use Order.ItemCount
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetOrdersResponse represents the response for getting merchant orders
type GetOrdersResponse struct {
	Orders []OrderInfo `json:"orders"`
	Count  int         `json:"count"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int         `json:"total"`
}

// OrderQueryParam represents query parameters for order listing
type OrderQueryParam struct {
	VendorID    uint     `query:"vendor_id" validate:"omitempty"`
	OrderStatus []string `query:"order_status" validate:"omitempty"`
	Limit       int      `query:"limit" validate:"min=1,max=100"`
	Offset      int      `query:"offset" validate:"min=0"`
	Ordering    string   `query:"ordering" validate:"omitempty,oneof=created_at updated_at total_amount"`
}

// UpdateOrderStatusRequest represents the request to update order status
type UpdateOrderStatusRequest struct {
	Status  string `json:"status" validate:"required"`
	Remarks string `json:"remarks,omitempty"`
}

// UpdateOrderStatusResponse represents the response for order status update
type UpdateOrderStatusResponse struct {
	OrderID     uint    `json:"order_id"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	Message     string  `json:"message"`
}

// OrderItemInfo represents order item information
type OrderItemInfo struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	SubTotal    float64 `json:"sub_total"`
}

// OrderDetailInfo represents detailed order information with items
type OrderDetailInfo struct {
	ID          uint            `json:"id"`
	OrderNumber string          `json:"order_number"`
	VendorID    uint            `json:"vendor_id"`
	VendorName  string          `json:"vendor_name"`
	Status      string          `json:"status"`
	TotalAmount float64         `json:"total_amount"`
	ItemCount   int             `json:"item_count"`
	Notes       string          `json:"notes,omitempty"`
	OrderItems  []OrderItemInfo `json:"order_items"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// GetOrderDetailResponse represents the response for getting order details
type GetOrderDetailResponse struct {
	Order OrderDetailInfo `json:"order"`
}

// DeleteOrderItemResponse represents the response for deleting order item
type DeleteOrderItemResponse struct {
	OrderItemID uint   `json:"order_item_id"`
	Message     string `json:"message"`
}

// ===== Order State Transition DTOs =====

// PlaceOrderRequest represents the request to place an order
type PlaceOrderRequest struct {
	Notes string `json:"notes,omitempty" validate:"max=500"`
}

// PlaceOrderResponse represents the response for placing an order
type PlaceOrderResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// MarkReceivedRequest represents the request to mark order as received
type MarkReceivedRequest struct {
	ReceivedDate string `json:"received_date,omitempty"` // ISO 8601 format
	Notes        string `json:"notes,omitempty" validate:"max=500"`
}

// MarkReceivedResponse represents the response for marking order as received
type MarkReceivedResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// CompleteOrderRequest represents the request to complete an order
type CompleteOrderRequest struct {
	Rating int    `json:"rating,omitempty" validate:"min=1,max=5"`
	Notes  string `json:"notes,omitempty" validate:"max=500"`
}

// CompleteOrderResponse represents the response for completing an order
type CompleteOrderResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// CancelOrderRequest represents the request to cancel an order
type CancelOrderRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
	Notes  string `json:"notes,omitempty" validate:"max=500"`
}

// CancelOrderResponse represents the response for cancelling an order
type CancelOrderResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
