package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/repository"
	validator "linksupply.io/vmconnect/api/merchant/validator"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils"
)

// OrderService defines the interface for order operations
type OrderService interface {
	ManageCartItem(ctx context.Context, merchantID, vendorID uint, req *dto.CartItemRequest) (*dto.CartItemResponse, *response.ErrorDetails)
	DeleteOrderItem(ctx context.Context, merchantID, orderID, orderItemID uint) (*dto.DeleteOrderItemResponse, *response.ErrorDetails)
	GetMerchantOrders(ctx context.Context, merchantID uint, queryParams *dto.OrderQueryParam) (*dto.GetOrdersResponse, *response.ErrorDetails)
	GetOrderDetail(ctx context.Context, merchantID, orderID uint) (*dto.GetOrderDetailResponse, *response.ErrorDetails)
	UpdateOrderStatus(ctx context.Context, merchantID, orderID uint, req *dto.UpdateOrderStatusRequest) (*dto.UpdateOrderStatusResponse, *response.ErrorDetails)
	// State transition methods
	PlaceOrder(ctx context.Context, merchantID, orderID uint, req *dto.PlaceOrderRequest) (*dto.PlaceOrderResponse, *response.ErrorDetails)
	MarkReceived(ctx context.Context, merchantID, orderID uint, req *dto.MarkReceivedRequest) (*dto.MarkReceivedResponse, *response.ErrorDetails)
	CompleteOrder(ctx context.Context, merchantID, orderID uint, req *dto.CompleteOrderRequest) (*dto.CompleteOrderResponse, *response.ErrorDetails)
	CancelOrder(ctx context.Context, merchantID, orderID uint, req *dto.CancelOrderRequest) (*dto.CancelOrderResponse, *response.ErrorDetails)
}

type OrderServiceImpl struct {
	repository     repository.MerchantRepository
	validator      validator.MerchantValidator
	stateValidator *utils.OrderStateValidator
}

// NewOrderService creates a new instance of order service
func NewOrderService(repo repository.MerchantRepository) OrderService {
	return &OrderServiceImpl{
		repository:     repo,
		validator:      validator.NewMerchantValidator(),
		stateValidator: utils.NewOrderStateValidator(),
	}
}

// GetMerchantOrders gets merchant orders with filtering and pagination
func (s *OrderServiceImpl) GetMerchantOrders(ctx context.Context, merchantID uint, queryParams *dto.OrderQueryParam) (*dto.GetOrdersResponse, *response.ErrorDetails) {
	// Validate query parameters
	if err := s.validator.ValidateOrderQueryParam(queryParams); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Set default values
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}
	if queryParams.Ordering == "" {
		queryParams.Ordering = "created_at"
	}

	// Parse order status values from array (convert enum names to integer strings)
	if len(queryParams.OrderStatus) > 0 {
		var statusInts []int
		for _, status := range queryParams.OrderStatus {
			status = strings.TrimSpace(status)
			if statusInt, isValid := models.ParseOrderStatusString(status); isValid {
				statusInts = append(statusInts, int(statusInt))
			}
		}
		// Convert to string array for query
		if len(statusInts) > 0 {
			var statusStrings []string
			for _, statusInt := range statusInts {
				statusStrings = append(statusStrings, strconv.Itoa(statusInt))
			}
			queryParams.OrderStatus = statusStrings
		} else {
			queryParams.OrderStatus = []string{}
		}
	}

	// Get orders
	orders, totalCount, err := s.repository.GetMerchantOrders(merchantID, queryParams)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant orders",
			Error:   err,
		}
	}

	// Convert to DTO
	orderInfos := make([]dto.OrderInfo, len(orders))
	for i, order := range orders {
		orderInfos[i] = dto.OrderInfo{
			ID:          order.ID,
			OrderNumber: strconv.FormatUint(uint64(order.ID), 10), // Simple order number
			VendorID:    order.VendorID,
			VendorName:  order.Vendor.CompanyName,
			Status:      order.Status.String(),
			TotalAmount: order.TotalAmount,
			ItemCount:   order.ItemCount, // Use the field from Order model
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
	}

	return &dto.GetOrdersResponse{
		Orders: orderInfos,
		Count:  len(orderInfos),
		Limit:  queryParams.Limit,
		Offset: queryParams.Offset,
		Total:  totalCount,
	}, nil
}

// updateOrderTotalAmount updates the total amount of an order
func (s *OrderServiceImpl) updateOrderTotalAmount(orderID uint) error {
	// This would typically involve calculating total from order items
	// For now, we'll implement a simple version
	var totalAmount float64
	var orderItems []models.OrderItem

	err := s.repository.GetOrderItemsByOrderID(orderID, &orderItems)
	if err != nil {
		return err
	}

	for _, item := range orderItems {
		totalAmount += item.SubTotal
	}

	// Update order total
	return s.repository.UpdateOrderTotalAmount(orderID, totalAmount)
}

// ManageCartItem manages cart items (add/update/remove)
func (s *OrderServiceImpl) ManageCartItem(ctx context.Context, merchantID, vendorID uint, req *dto.CartItemRequest) (*dto.CartItemResponse, *response.ErrorDetails) {
	// Get merchant to obtain user_id
	merchant, err := s.repository.GetMerchantByID(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant details",
			Error:   err,
		}
	}

	// Get or create cart order
	cartOrder, err := s.repository.GetOrCreateCartOrder(merchantID, vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get or create cart order",
			Error:   err,
		}
	}

	// Validate that order is in CART status for cart operations
	if err := s.validator.ValidateOrderStatusForCartOperation(cartOrder.Status); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Get product details
	product, err := s.repository.GetProductByID(req.ProductID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: "Product not found",
			Error:   err,
		}
	}

	var action string
	var message string

	// Add or update item (quantity is already validated to be > 0)
	existingItem, err := s.repository.GetOrderItemByProductID(cartOrder.ID, req.ProductID)
	if err == nil {
		// Update existing item
		_, err = s.repository.UpdateCartItemQuantity(existingItem.ID, req.Quantity)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update item quantity",
				Error:   err,
			}
		}
		action = "updated"
		message = "Product quantity updated in cart"
	} else {
		// Add new item
		_, err = s.repository.AddToCart(cartOrder.ID, req.ProductID, req.Quantity, getProductPrice(product))
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to add product to cart",
				Error:   err,
			}
		}
		action = "added"
		message = "Product added to cart"
	}

	// Update order total amount
	err = s.updateOrderTotalAmount(cartOrder.ID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update order total",
			Error:   err,
		}
	}

	orderStatus := models.ORDER_IN_CART
	actorRole := models.ACTOR_MERCHANT
	// Create order activity
	activity := &models.OrderActivity{
		OrderID:    cartOrder.ID,
		ActorID:    merchant.OwnerID, // Store user_id from merchant
		ActorRole:  &actorRole,
		OrderState: &orderStatus,
		Remarks:    message,
		Metadata:   nil, // Can be populated with additional data if needed
		CreatedAt:  time.Now(),
		IsDeleted:  false,
		UpdatedAt:  time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return &dto.CartItemResponse{
		OrderID:   cartOrder.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Action:    action,
		Message:   message,
	}, nil
}

// UpdateOrderStatus updates order status
func (s *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, merchantID, orderID uint, req *dto.UpdateOrderStatusRequest) (*dto.UpdateOrderStatusResponse, *response.ErrorDetails) {
	// Get merchant to obtain user_id
	merchant, err := s.repository.GetMerchantByID(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant details",
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

	// Verify order belongs to merchant
	if !orderBelongsToMerchant(order, merchantID) {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this merchant",
			Error:   nil,
		}
	}

	// Parse new status
	newStatus, isValid := models.ParseOrderStatusString(req.Status)
	if !isValid {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: "Invalid order status",
			Error:   nil,
		}
	}

	// Update order status
	err = s.repository.UpdateOrderStatus(orderID, newStatus)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update order status",
			Error:   err,
		}
	}

	// Create order activity
	remarks := req.Remarks
	if remarks == "" {
		remarks = "Order status updated to " + req.Status
	}

	actorRole := models.ACTOR_MERCHANT
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    merchant.OwnerID, // Store user_id from merchant
		ActorRole:  &actorRole,
		OrderState: &newStatus,
		Remarks:    remarks,
		Metadata:   nil, // Can be populated with additional data if needed
		CreatedAt:  time.Now(),
		IsDeleted:  false,
		UpdatedAt:  time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the status update
		// TODO: Add proper logging
	}

	return &dto.UpdateOrderStatusResponse{
		OrderID:     orderID,
		Status:      newStatus.String(),
		TotalAmount: order.TotalAmount,
		Message:     "Order status updated successfully",
	}, nil
}

// GetOrderDetail gets detailed order information with items
func (s *OrderServiceImpl) GetOrderDetail(ctx context.Context, merchantID, orderID uint) (*dto.GetOrderDetailResponse, *response.ErrorDetails) {
	// Get order with items
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Order not found",
			Error:   err,
		}
	}

	// Verify order belongs to merchant
	if !orderBelongsToMerchant(order, merchantID) {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this merchant",
			Error:   nil,
		}
	}

	// Convert order items to DTO
	orderItems := make([]dto.OrderItemInfo, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = dto.OrderItemInfo{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name, // Assuming Product is preloaded
			Quantity:    item.Quantity,
			Price:       item.Price,
			SubTotal:    item.SubTotal,
		}
	}

	// Convert to detailed DTO
	orderDetail := dto.OrderDetailInfo{
		ID:          order.ID,
		OrderNumber: strconv.FormatUint(uint64(order.ID), 10),
		VendorID:    order.VendorID,
		VendorName:  order.Vendor.CompanyName,
		Status:      order.Status.String(),
		TotalAmount: order.TotalAmount,
		ItemCount:   order.ItemCount,
		Notes:       order.Notes,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	return &dto.GetOrderDetailResponse{
		Order: orderDetail,
	}, nil
}

// DeleteOrderItem deletes an order item from cart
func (s *OrderServiceImpl) DeleteOrderItem(ctx context.Context, merchantID, orderID, orderItemID uint) (*dto.DeleteOrderItemResponse, *response.ErrorDetails) {
	// Get the order to verify it belongs to the merchant and check status
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Order not found",
			Error:   err,
		}
	}

	// Verify order belongs to merchant
	if !orderBelongsToMerchant(order, merchantID) {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this merchant",
			Error:   nil,
		}
	}

	// Validate that order is in CART status for cart operations
	if err := s.validator.ValidateOrderStatusForCartOperation(order.Status); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Get order item details and verify it belongs to this order
	orderItem, err := s.repository.GetOrderItemByID(orderItemID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Order item not found",
			Error:   err,
		}
	}

	// Verify order item belongs to the specified order
	if orderItem.OrderID != orderID {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: "Order item does not belong to the specified order",
			Error:   nil,
		}
	}

	// Hard delete the order item
	err = s.repository.HardDeleteOrderItem(orderItemID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete order item",
			Error:   err,
		}
	}

	// Update order total amount
	err = s.updateOrderTotalAmount(order.ID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update order total",
			Error:   err,
		}
	}

	return &dto.DeleteOrderItemResponse{
		OrderItemID: orderItemID,
		Message:     "Order item deleted successfully",
	}, nil
}

// ===== State Transition Methods =====

// PlaceOrder transitions order from IN_CART to PLACED
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, merchantID, orderID uint, req *dto.PlaceOrderRequest) (*dto.PlaceOrderResponse, *response.ErrorDetails) {
	// Get merchant to obtain user_id
	merchant, err := s.repository.GetMerchantByID(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant details",
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

	// Verify order belongs to merchant
	if !orderBelongsToMerchant(order, merchantID) {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this merchant",
			Error:   nil,
		}
	}

	targetStatus := models.ORDER_PLACED

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
		return &dto.PlaceOrderResponse{
			OrderID:   orderID,
			Status:    targetStatus.String(),
			Message:   "Order was already placed",
			Timestamp: lastActivity.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Validate state transition
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_MERCHANT); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Check if order has items
	if order.ItemCount == 0 {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: "Cannot place order with no items",
			Error:   nil,
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
	actorRole := models.ACTOR_MERCHANT
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    merchant.OwnerID,
		ActorRole:  &actorRole,
		OrderState: &targetStatus,
		Remarks:    req.Notes,
		Metadata: models.JSONMap{
			"action":          "place_order",
			"previous_status": order.Status.String(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return &dto.PlaceOrderResponse{
		OrderID:   orderID,
		Status:    targetStatus.String(),
		Message:   "Order placed successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// MarkReceived transitions order from SHIPPED to DELIVERED
func (s *OrderServiceImpl) MarkReceived(ctx context.Context, merchantID, orderID uint, req *dto.MarkReceivedRequest) (*dto.MarkReceivedResponse, *response.ErrorDetails) {
	// Get merchant to obtain user_id
	merchant, err := s.repository.GetMerchantByID(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant details",
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

	// Verify order belongs to merchant
	if !orderBelongsToMerchant(order, merchantID) {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this merchant",
			Error:   nil,
		}
	}

	targetStatus := models.ORDER_DELIVERED

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
		return &dto.MarkReceivedResponse{
			OrderID:   orderID,
			Status:    targetStatus.String(),
			Message:   "Order was already marked as delivered",
			Timestamp: lastActivity.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Validate state transition
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_MERCHANT); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Parse received date if provided
	var receivedDate time.Time
	if req.ReceivedDate != "" {
		receivedDate, err = time.Parse(time.RFC3339, req.ReceivedDate)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: "Invalid received_date format. Use ISO 8601 format",
				Error:   err,
			}
		}
	} else {
		receivedDate = time.Now()
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
	actorRole := models.ACTOR_MERCHANT
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    merchant.OwnerID,
		ActorRole:  &actorRole,
		OrderState: &targetStatus,
		Remarks:    req.Notes,
		Metadata: models.JSONMap{
			"action":          "mark_received",
			"previous_status": order.Status.String(),
			"received_date":   receivedDate.Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return &dto.MarkReceivedResponse{
		OrderID:   orderID,
		Status:    targetStatus.String(),
		Message:   "Order marked as received successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// CompleteOrder transitions order from DELIVERED to COMPLETED
func (s *OrderServiceImpl) CompleteOrder(ctx context.Context, merchantID, orderID uint, req *dto.CompleteOrderRequest) (*dto.CompleteOrderResponse, *response.ErrorDetails) {
	// Get merchant to obtain user_id
	merchant, err := s.repository.GetMerchantByID(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant details",
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

	// Verify order belongs to merchant
	if !orderBelongsToMerchant(order, merchantID) {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this merchant",
			Error:   nil,
		}
	}

	targetStatus := models.ORDER_COMPLETED

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
		return &dto.CompleteOrderResponse{
			OrderID:   orderID,
			Status:    targetStatus.String(),
			Message:   "Order was already completed",
			Timestamp: lastActivity.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Validate state transition
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_MERCHANT); err != nil {
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
	actorRole := models.ACTOR_MERCHANT
	metadata := models.JSONMap{
		"action":          "complete_order",
		"previous_status": order.Status.String(),
	}

	if req.Rating > 0 {
		metadata["rating"] = req.Rating
	}

	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    merchant.OwnerID,
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
		// TODO: Add proper logging
	}

	return &dto.CompleteOrderResponse{
		OrderID:   orderID,
		Status:    targetStatus.String(),
		Message:   "Order completed successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// CancelOrder transitions order to CANCELLED
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, merchantID, orderID uint, req *dto.CancelOrderRequest) (*dto.CancelOrderResponse, *response.ErrorDetails) {
	// Get merchant to obtain user_id
	merchant, err := s.repository.GetMerchantByID(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant details",
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

	// Verify order belongs to merchant
	if !orderBelongsToMerchant(order, merchantID) {
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Order does not belong to this merchant",
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
	if err := s.stateValidator.ValidateTransition(*order.Status, targetStatus, models.ACTOR_MERCHANT); err != nil {
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
	actorRole := models.ACTOR_MERCHANT
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    merchant.OwnerID,
		ActorRole:  &actorRole,
		OrderState: &targetStatus,
		Remarks:    req.Notes,
		Metadata: models.JSONMap{
			"action":          "cancel_order",
			"previous_status": order.Status.String(),
			"reason":          req.Reason,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.repository.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return &dto.CancelOrderResponse{
		OrderID:   orderID,
		Status:    targetStatus.String(),
		Message:   "Order cancelled successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func getProductPrice(product *models.Product) float64 {
	if product == nil || len(product.Variants) == 0 {
		return 0
	}
	return product.Variants[0].Price
}

func orderBelongsToMerchant(order *models.Order, merchantID uint) bool {
	return order != nil && order.MerchantID != nil && *order.MerchantID == merchantID
}
