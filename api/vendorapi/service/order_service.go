package service

import (
	"context"
	"net/http"
	"time"

	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils"
)

// OrderService defines the interface for vendor order operations
type OrderService interface {
	ConfirmOrder(ctx context.Context, vendorID, orderID uint, req *dto.ConfirmOrderRequest) (*dto.ConfirmOrderResponse, *response.ErrorDetails)
	GenerateInvoice(ctx context.Context, vendorID, orderID uint, req *dto.GenerateInvoiceRequest) (*dto.GenerateInvoiceResponse, *response.ErrorDetails)
	DispatchOrder(ctx context.Context, vendorID, orderID uint, req *dto.DispatchOrderRequest) (*dto.DispatchOrderResponse, *response.ErrorDetails)
	CancelOrder(ctx context.Context, vendorID, orderID uint, req *dto.CancelOrderRequest) (*dto.CancelOrderResponse, *response.ErrorDetails)
}

type OrderServiceImpl struct {
	repository     repository.VendorRepository
	stateValidator *utils.OrderStateValidator
}

// NewOrderService creates a new instance of vendor order service
func NewOrderService(repo repository.VendorRepository) OrderService {
	return &OrderServiceImpl{
		repository:     repo,
		stateValidator: utils.NewOrderStateValidator(),
	}
}

// ConfirmOrder transitions order from PLACED to CONFIRMED
func (s *OrderServiceImpl) ConfirmOrder(ctx context.Context, vendorID, orderID uint, req *dto.ConfirmOrderRequest) (*dto.ConfirmOrderResponse, *response.ErrorDetails) {
	// Get vendor to obtain user_id
	vendor, err := s.repository.GetVendorByID(vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor details",
			Error:   err,
		}
	}

	// Get order details
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Order not found",
			Error:   err,
		}
	}

	// Verify order belongs to vendor
	if order.VendorID != vendorID {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this vendor",
			Error:   nil,
		}
	}

	targetStatus := models.ORDER_CONFIRMED

	// Check if status already exists in history
	existsInHistory, lastActivity, err := s.repository.CheckStatusInHistory(orderID, targetStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check order history",
			Error:   err,
		}
	}

	if existsInHistory {
		return &dto.ConfirmOrderResponse{
			OrderID:   orderID,
			Status:    targetStatus.String(),
			Message:   "Order was already confirmed",
			Timestamp: lastActivity.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Validate state transition
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_VENDOR); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Update order status
	err = s.repository.UpdateOrderStatus(orderID, targetStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update order status",
			Error:   err,
		}
	}

	// Create order activity
	actorRole := models.ACTOR_VENDOR
	metadata := models.JSONMap{
		"action":          "confirm_order",
		"previous_status": order.Status.String(),
	}

	if req.EstimatedDelivery != "" {
		metadata["estimated_delivery"] = req.EstimatedDelivery
	}

	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    vendor.UserID,
		ActorRole:  &actorRole,
		OrderState: &targetStatus,
		Remarks:    req.Notes,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
	}

	return &dto.ConfirmOrderResponse{
		OrderID:   orderID,
		Status:    targetStatus.String(),
		Message:   "Order confirmed successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// GenerateInvoice transitions order to INVOICED
func (s *OrderServiceImpl) GenerateInvoice(ctx context.Context, vendorID, orderID uint, req *dto.GenerateInvoiceRequest) (*dto.GenerateInvoiceResponse, *response.ErrorDetails) {
	// Get vendor to obtain user_id
	vendor, err := s.repository.GetVendorByID(vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor details",
			Error:   err,
		}
	}

	// Get order details
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Order not found",
			Error:   err,
		}
	}

	// Verify order belongs to vendor
	if order.VendorID != vendorID {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this vendor",
			Error:   nil,
		}
	}

	targetStatus := models.ORDER_INVOICED

	// Validate state transition
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_VENDOR); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Parse invoice date
	_, err = time.Parse(time.RFC3339, req.InvoiceDate)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: "Invalid invoice_date format. Use ISO 8601 format",
			Error:   err,
		}
	}

	// Update order status
	err = s.repository.UpdateOrderStatus(orderID, targetStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update order status",
			Error:   err,
		}
	}

	// Create order activity with invoice details
	actorRole := models.ACTOR_VENDOR
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    vendor.UserID,
		ActorRole:  &actorRole,
		OrderState: &targetStatus,
		Remarks:    req.Notes,
		Metadata: models.JSONMap{
			"action":          "generate_invoice",
			"previous_status": order.Status.String(),
			"invoice_info": map[string]interface{}{
				"invoice_no":   req.InvoiceNo,
				"invoice_date": req.InvoiceDate,
				"amount":       req.Amount,
				"tax_amount":   req.TaxAmount,
				"file_url":     req.FileURL,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
	}

	return &dto.GenerateInvoiceResponse{
		OrderID:   orderID,
		InvoiceID: 0, // Will be set when Invoice record is created
		InvoiceNo: req.InvoiceNo,
		Status:    targetStatus.String(),
		Message:   "Invoice generated successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// DispatchOrder transitions order to SHIPPED
func (s *OrderServiceImpl) DispatchOrder(ctx context.Context, vendorID, orderID uint, req *dto.DispatchOrderRequest) (*dto.DispatchOrderResponse, *response.ErrorDetails) {
	// Get vendor to obtain user_id
	vendor, err := s.repository.GetVendorByID(vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor details",
			Error:   err,
		}
	}

	// Get order details
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Order not found",
			Error:   err,
		}
	}

	// Verify order belongs to vendor
	if order.VendorID != vendorID {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this vendor",
			Error:   nil,
		}
	}

	targetStatus := models.ORDER_SHIPPED

	// Check if status already exists in history
	existsInHistory, lastActivity, err := s.repository.CheckStatusInHistory(orderID, targetStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check order history",
			Error:   err,
		}
	}

	if existsInHistory {
		return &dto.DispatchOrderResponse{
			OrderID:    orderID,
			TrackingID: req.TrackingID,
			Status:     targetStatus.String(),
			Message:    "Order was already dispatched",
			Timestamp:  lastActivity.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Validate state transition
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_VENDOR); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Parse dispatch date if provided
	var dispatchDate time.Time
	if req.DispatchDate != "" {
		dispatchDate, err = time.Parse(time.RFC3339, req.DispatchDate)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: "Invalid dispatch_date format. Use ISO 8601 format",
				Error:   err,
			}
		}
	} else {
		dispatchDate = time.Now()
	}

	// Update order status
	err = s.repository.UpdateOrderStatus(orderID, targetStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update order status",
			Error:   err,
		}
	}

	// Create order activity with dispatch details
	actorRole := models.ACTOR_VENDOR
	dispatchInfo := map[string]interface{}{
		"tracking_id":   req.TrackingID,
		"carrier":       req.Carrier,
		"dispatch_date": dispatchDate.Format(time.RFC3339),
	}

	if req.ExpectedDelivery != "" {
		dispatchInfo["expected_delivery"] = req.ExpectedDelivery
	}
	if req.Packages > 0 {
		dispatchInfo["packages"] = req.Packages
	}

	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    vendor.UserID,
		ActorRole:  &actorRole,
		OrderState: &targetStatus,
		Remarks:    req.Notes,
		Metadata: models.JSONMap{
			"action":          "dispatch_order",
			"previous_status": order.Status.String(),
			"dispatch_info":   dispatchInfo,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
	}

	return &dto.DispatchOrderResponse{
		OrderID:    orderID,
		TrackingID: req.TrackingID,
		Status:     targetStatus.String(),
		Message:    "Order dispatched successfully",
		Timestamp:  time.Now().Format(time.RFC3339),
	}, nil
}

// VerifyPayment updates payment status (doesn't change order status)
// CancelOrder transitions order to CANCELLED
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, vendorID, orderID uint, req *dto.CancelOrderRequest) (*dto.CancelOrderResponse, *response.ErrorDetails) {
	// Get vendor to obtain user_id
	vendor, err := s.repository.GetVendorByID(vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor details",
			Error:   err,
		}
	}

	// Get order details
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Order not found",
			Error:   err,
		}
	}

	// Verify order belongs to vendor
	if order.VendorID != vendorID {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this vendor",
			Error:   nil,
		}
	}

	targetStatus := models.ORDER_CANCELLED

	// Check if status already exists in history
	existsInHistory, lastActivity, err := s.repository.CheckStatusInHistory(orderID, targetStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check order history",
			Error:   err,
		}
	}

	if existsInHistory {
		return &dto.CancelOrderResponse{
			OrderID:   orderID,
			Status:    targetStatus.String(),
			Message:   "Order was already cancelled",
			Timestamp: lastActivity.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Validate state transition
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_VENDOR); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Update order status
	err = s.repository.UpdateOrderStatus(orderID, targetStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update order status",
			Error:   err,
		}
	}

	// Create order activity
	actorRole := models.ACTOR_VENDOR
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    vendor.UserID,
		ActorRole:  &actorRole,
		OrderState: &targetStatus,
		Remarks:    req.Notes,
		Metadata: models.JSONMap{
			"action":          "cancel_order",
			"previous_status": order.Status.String(),
			"reason":          req.Reason,
			"refund_required": req.RefundRequired,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
	}

	return &dto.CancelOrderResponse{
		OrderID:   orderID,
		Status:    targetStatus.String(),
		Message:   "Order cancelled successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
