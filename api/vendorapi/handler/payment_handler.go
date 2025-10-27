package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/service"
)

// PaymentHandler handles payment-related HTTP requests for vendors
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler creates a new payment handler instance
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// MarkOrderAsPaid godoc
// @Summary Mark order as paid (bulk verify all payments)
// @Description Vendor marks all payments of an order as verified in one action
// @Tags Vendor Payments
// @Accept json
// @Produce json
// @Param vid path int true "Vendor ID"
// @Param oid path int true "Order ID"
// @Param request body dto.MarkOrderPaidRequest true "Mark paid details"
// @Success 200 {object} dto.MarkOrderPaidResponse
// @Failure 400 {object} response.ErrorDetails
// @Failure 401 {object} response.ErrorDetails
// @Failure 500 {object} response.ErrorDetails
// @Router /api/v0/vendor/{vid}/orders/{oid}/mark-paid [post]
func (h *PaymentHandler) MarkOrderAsPaid(c *fiber.Ctx) error {
	vendorID, err := strconv.ParseUint(c.Params("vid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid vendor ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	orderID, err := strconv.ParseUint(c.Params("oid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid order ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	var req dto.MarkOrderPaidRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, err := h.paymentService.MarkOrderAsPaid(c.Context(), uint(vendorID), uint(orderID), &req)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(500, err.Error(), nil)
		return response.WriteHTTPResponse(c, http.StatusInternalServerError, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// VerifyPayment godoc
// @Summary Verify individual payment (optional)
// @Description Vendor verifies or rejects a specific payment
// @Tags Vendor Payments
// @Accept json
// @Produce json
// @Param vid path int true "Vendor ID"
// @Param oid path int true "Order ID"
// @Param pid path int true "Payment ID"
// @Param request body dto.VerifyPaymentRequest true "Verification details"
// @Success 200 {object} dto.VerifyPaymentResponse
// @Failure 400 {object} response.ErrorDetails
// @Failure 401 {object} response.ErrorDetails
// @Failure 500 {object} response.ErrorDetails
// @Router /api/v0/vendor/{vid}/orders/{oid}/payments/{pid}/verify [post]
func (h *PaymentHandler) VerifyPayment(c *fiber.Ctx) error {
	vendorID, err := strconv.ParseUint(c.Params("vid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid vendor ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	orderID, err := strconv.ParseUint(c.Params("oid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid order ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	paymentID, err := strconv.ParseUint(c.Params("pid"), 10, 64)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid payment ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	var req dto.VerifyPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, err := h.paymentService.VerifyPayment(c.Context(), uint(vendorID), uint(orderID), paymentID, &req)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(500, err.Error(), nil)
		return response.WriteHTTPResponse(c, http.StatusInternalServerError, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// GetOrderPayments godoc
// @Summary Get all payments for an order (vendor view)
// @Description Retrieve all payment records for a specific order
// @Tags Vendor Payments
// @Accept json
// @Produce json
// @Param vid path int true "Vendor ID"
// @Param oid path int true "Order ID"
// @Success 200 {object} dto.GetOrderPaymentsResponse
// @Failure 400 {object} response.ErrorDetails
// @Failure 401 {object} response.ErrorDetails
// @Failure 500 {object} response.ErrorDetails
// @Router /api/v0/vendor/{vid}/orders/{oid}/payments [get]
func (h *PaymentHandler) GetOrderPayments(c *fiber.Ctx) error {
	vendorID, err := strconv.ParseUint(c.Params("vid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid vendor ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	orderID, err := strconv.ParseUint(c.Params("oid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid order ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	result, err := h.paymentService.GetOrderPayments(c.Context(), uint(vendorID), uint(orderID))
	if err != nil {
		body := response.GetErrorHTTPResponseBody(500, err.Error(), nil)
		return response.WriteHTTPResponse(c, http.StatusInternalServerError, body)
	}

	body := &response.HTTPResponse{Content: result}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}
