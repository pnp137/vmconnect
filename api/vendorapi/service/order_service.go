package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	validator "linksupply.io/vmconnect/api/vendorapi/validator"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils"
)

// OrderService defines the interface for vendor order operations
type OrderService interface {
	CreateOrder(ctx context.Context, vendorID uint, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *response.ErrorDetails)
	ConfirmOrder(ctx context.Context, vendorID, orderID uint, req *dto.ConfirmOrderRequest) (*dto.ConfirmOrderResponse, *response.ErrorDetails)
	GenerateInvoice(ctx context.Context, vendorID, orderID uint, req *dto.GenerateInvoiceRequest) (*dto.GenerateInvoiceResponse, *response.ErrorDetails)
	DispatchOrder(ctx context.Context, vendorID, orderID uint, req *dto.DispatchOrderRequest) (*dto.DispatchOrderResponse, *response.ErrorDetails)
	CancelOrder(ctx context.Context, vendorID, orderID uint, req *dto.CancelOrderRequest) (*dto.CancelOrderResponse, *response.ErrorDetails)
	GetOrders(ctx context.Context, vendorID uint, params *dto.OrderQueryParam) (*dto.GetOrdersResponse, *response.ErrorDetails)
	GetOrder(ctx context.Context, vendorID, orderID uint) (*dto.GetOrderResponse, *response.ErrorDetails)
}

type OrderServiceImpl struct {
	repository     repository.VendorRepository
	stateValidator *utils.OrderStateValidator
	validator      validator.VendorValidator
}

// NewOrderService creates a new instance of vendor order service
func NewOrderService(repo repository.VendorRepository) OrderService {
	return &OrderServiceImpl{
		repository:     repo,
		stateValidator: utils.NewOrderStateValidator(),
		validator:      validator.NewVendorValidator(),
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, vendorID uint, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *response.ErrorDetails) {
	if err := s.validator.ValidateCreateOrderRequest(req); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	if _, err := s.repository.GetVendorByID(vendorID); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: "vendor not found",
			Error:   err,
		}
	}

	variantIDs := uniqueVariantIDs(req.Items)
	variants, err := s.repository.GetProductVariantsByIDs(variantIDs)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "failed to fetch product variants",
			Error:   err,
		}
	}

	variantMap := make(map[uint]models.ProductVariant, len(variants))
	for _, variant := range variants {
		variantMap[variant.ID] = variant
	}

	orderItems := make([]models.OrderItem, 0, len(req.Items))
	responseItems := make([]dto.OrderItemResponse, 0, len(req.Items))
	orderTotal := 0.0

	for _, item := range req.Items {
		variant, found := variantMap[item.ProductVariantID]
		if !found {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("product variant %d not found", item.ProductVariantID),
				Error:   fmt.Errorf("variant %d not found", item.ProductVariantID),
			}
		}

		if !variant.IsActive {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("product variant %s is not available", variant.Name),
				Error:   fmt.Errorf("variant inactive"),
			}
		}

		if variant.Product.VendorID != vendorID {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("variant %s does not belong to vendor %d", variant.Name, vendorID),
				Error:   fmt.Errorf("variant vendor mismatch"),
			}
		}

		if item.Quantity < variant.MOQ {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("minimum order quantity for %s is %.0f", variant.Name, variant.MOQ),
				Error:   fmt.Errorf("quantity below moq"),
			}
		}

		if item.Quantity > float64(variant.Stock) {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("insufficient stock for variant %s", variant.Name),
				Error:   fmt.Errorf("insufficient stock"),
			}
		}

		quantityInt := int(item.Quantity)
		if float64(quantityInt) != item.Quantity {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("quantity for variant %s must be a whole number", variant.Name),
				Error:   fmt.Errorf("quantity must be whole number"),
			}
		}

		itemTotal := item.Quantity * variant.Price
		orderTotal += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:        variant.ProductID,
			ProductVariantID: variant.ID,
			ProductName:      variant.Product.Name,
			VariantName:      variant.Name,
			Quantity:         item.Quantity,
			UnitPrice:        variant.Price,
			TotalPrice:       itemTotal,
			Price:            variant.Price,
			SubTotal:         itemTotal,
		})

		responseItems = append(responseItems, dto.OrderItemResponse{
			ProductName: variant.Product.Name,
			VariantName: variant.Name,
			Quantity:    item.Quantity,
			UnitPrice:   variant.Price,
			TotalPrice:  itemTotal,
		})
	}

	orderStatus := models.ORDER_PLACED
	var latitude *float64
	var longitude *float64
	googleMapsURL := ""
	if req.Location != nil {
		latitude = &req.Location.Latitude
		longitude = &req.Location.Longitude
		googleMapsURL = buildGoogleMapsURL(req.Location.Latitude, req.Location.Longitude)
	}

	order := &models.Order{
		OrderNo:         fmt.Sprintf("ORDER-%d", time.Now().UnixNano()),
		VendorID:        vendorID,
		OrderFor:        req.OrderFor,
		CustomerName:    req.CustomerName,
		CustomerMobile:  req.CustomerMobile,
		Latitude:        latitude,
		Longitude:       longitude,
		GoogleMapsURL:   googleMapsURL,
		BusinessName:    req.BusinessName,
		DeliveryAddress: req.DeliveryAddress,
		Notes:           req.Notes,
		Status:          &orderStatus,
		TotalAmount:     orderTotal,
		ItemCount:       len(orderItems),
	}

	tx := s.repository.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "failed to start database transaction",
			Error:   tx.Error,
		}
	}

	if err := s.repository.CreateOrder(tx, order); err != nil {
		tx.Rollback()
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "failed to create order",
			Error:   err,
		}
	}

	for index := range orderItems {
		orderItems[index].OrderID = order.ID
		if err := s.repository.CreateOrderItem(tx, &orderItems[index]); err != nil {
			tx.Rollback()
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "failed to create order item",
				Error:   err,
			}
		}
	}

	for _, item := range orderItems {
		quantityInt := int(item.Quantity)
		if err := s.repository.ReduceProductVariantStock(tx, item.ProductVariantID, quantityInt); err != nil {
			tx.Rollback()
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Error:   err,
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "failed to commit order transaction",
			Error:   err,
		}
	}

	var locationResponse *dto.OrderLocationResponse
	if order.Latitude != nil && order.Longitude != nil {
		locationResponse = &dto.OrderLocationResponse{
			Latitude:      order.Latitude,
			Longitude:     order.Longitude,
			GoogleMapsURL: order.GoogleMapsURL,
		}
	}

	return &dto.CreateOrderResponse{
		OrderID:     order.ID,
		Status:      order.Status.String(),
		TotalAmount: orderTotal,
		Location:    locationResponse,
		Items:       responseItems,
		CreatedAt:   order.CreatedAt,
	}, nil
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

// GetOrders returns paginated list of orders for a vendor with optional filters
func (s *OrderServiceImpl) GetOrders(ctx context.Context, vendorID uint, params *dto.OrderQueryParam) (*dto.GetOrdersResponse, *response.ErrorDetails) {
	// set defaults
	if params == nil {
		params = &dto.OrderQueryParam{Limit: 20}
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.Ordering == "" {
		params.Ordering = "created_at"
	}

	orders, total, err := s.repository.GetOrders(vendorID, params)
	if err != nil {
		return nil, &response.ErrorDetails{Code: http.StatusInternalServerError, Message: "Failed to get orders", Error: err}
	}

	infos := make([]dto.OrderInfo, 0, len(orders))
	for _, o := range orders {
		status := ""
		if o.Status != nil {
			status = o.Status.String()
		}
		infos = append(infos, dto.OrderInfo{
			ID:             o.ID,
			OrderNo:        o.OrderNo,
			CustomerName:   o.CustomerName,
			CustomerMobile: o.CustomerMobile,
			BusinessName:   o.BusinessName,
			Status:         status,
			TotalAmount:    o.TotalAmount,
			CreatedAt:      o.CreatedAt,
		})
	}

	return &dto.GetOrdersResponse{
		Orders: infos,
		Count:  len(infos),
		Limit:  params.Limit,
		Offset: params.Offset,
		Total:  int(total),
	}, nil
}

// GetOrder returns a single order if it belongs to vendor
func (s *OrderServiceImpl) GetOrder(ctx context.Context, vendorID, orderID uint) (*dto.GetOrderResponse, *response.ErrorDetails) {
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, &response.ErrorDetails{Code: http.StatusNotFound, Message: "Order not found", Error: err}
	}
	if order.VendorID != vendorID {
		return nil, &response.ErrorDetails{Code: http.StatusForbidden, Message: "Order does not belong to this vendor", Error: nil}
	}
	return &dto.GetOrderResponse{Order: order}, nil
}

func uniqueVariantIDs(items []dto.OrderItemRequest) []uint {
	ids := make([]uint, 0, len(items))
	seen := make(map[uint]struct{})
	for _, item := range items {
		if _, ok := seen[item.ProductVariantID]; !ok {
			ids = append(ids, item.ProductVariantID)
			seen[item.ProductVariantID] = struct{}{}
		}
	}
	return ids
}

func buildGoogleMapsURL(latitude, longitude float64) string {
	return fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", latitude, longitude)
}
