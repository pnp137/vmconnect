package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/service"
	validator "linksupply.io/vmconnect/api/vendorapi/validator"
)

type OrderHandler struct {
	orderService service.OrderService
	validator    validator.VendorValidator
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		validator:    validator.NewVendorValidator(),
	}
}

// CreateOrder places a new order for the scoped vendor.
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	vendorID, err := h.validator.ValidateIDParameter(c.Params("vid"), "vendor ID")
	if err != nil {
		body := response.GetErrorHTTPResponseBody(http.StatusBadRequest, err.Error(), nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(http.StatusBadRequest, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.CreateOrder(c.Context(), vendorID, &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{
		Code:    http.StatusOK,
		Message: "Order placed successfully",
		Content: result,
	}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// ConfirmOrder confirms an order
func (h *OrderHandler) ConfirmOrder(c *fiber.Ctx) error {
	vendorID, _ := strconv.ParseUint(c.Params("vid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.ConfirmOrderRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.ConfirmOrder(c.Context(), uint(vendorID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// GenerateInvoice generates an invoice for an order
func (h *OrderHandler) GenerateInvoice(c *fiber.Ctx) error {
	vendorID, _ := strconv.ParseUint(c.Params("vid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.GenerateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.GenerateInvoice(c.Context(), uint(vendorID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// DispatchOrder dispatches an order
func (h *OrderHandler) DispatchOrder(c *fiber.Ctx) error {
	vendorID, _ := strconv.ParseUint(c.Params("vid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.DispatchOrderRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.DispatchOrder(c.Context(), uint(vendorID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// CancelOrder cancels an order
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	vendorID, _ := strconv.ParseUint(c.Params("vid"), 10, 32)
	orderID, _ := strconv.ParseUint(c.Params("oid"), 10, 32)

	var req dto.CancelOrderRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, errDetails := h.orderService.CancelOrder(c.Context(), uint(vendorID), uint(orderID), &req)
	if errDetails != nil {
		body := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message, nil)
		return response.WriteHTTPResponse(c, errDetails.Code, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}
