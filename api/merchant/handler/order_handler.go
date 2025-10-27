package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/service"
	"linksupply.io/vmconnect/api/response"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// GetMerchantOrders retrieves orders for a merchant
func (h *OrderHandler) GetMerchantOrders(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)

	var queryParams dto.OrderQueryParam
	if err := c.QueryParser(&queryParams); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid query parameters", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.GetMerchantOrders(c.Context(), uint(merchantID), &queryParams)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// GetOrderDetail retrieves order details
func (h *OrderHandler) GetOrderDetail(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	result, errDetails := h.orderService.GetOrderDetail(c.Context(), uint(merchantID), uint(orderID))
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// ManageCartItem adds or updates items in cart
func (h *OrderHandler) ManageCartItem(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	vendorID, _ := strconv.ParseUint(c.Params("vid"), 10, 32)

	var req dto.CartItemRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.ManageCartItem(c.Context(), uint(merchantID), uint(vendorID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// DeleteOrderItem deletes an item from order
func (h *OrderHandler) DeleteOrderItem(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)
	itemID, _ := strconv.ParseUint(c.Params("itemId"), 10, 32)

	result, errDetails := h.orderService.DeleteOrderItem(c.Context(), uint(merchantID), uint(orderID), uint(itemID))
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// PlaceOrder places an order (cart -> placed)
func (h *OrderHandler) PlaceOrder(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.PlaceOrderRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.PlaceOrder(c.Context(), uint(merchantID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// MarkReceived marks order as received
func (h *OrderHandler) MarkReceived(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.MarkReceivedRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.MarkReceived(c.Context(), uint(merchantID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// CompleteOrder marks order as completed
func (h *OrderHandler) CompleteOrder(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.CompleteOrderRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.CompleteOrder(c.Context(), uint(merchantID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// CancelOrder cancels an order
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.CancelOrderRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.CancelOrder(c.Context(), uint(merchantID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// UpdateOrderStatus updates order status (legacy endpoint)
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	merchantID, _ := strconv.ParseUint(c.Params("mid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.UpdateOrderStatus(c.Context(), uint(merchantID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}
