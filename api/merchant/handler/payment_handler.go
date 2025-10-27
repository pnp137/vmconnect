package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/service"
	"linksupply.io/vmconnect/api/response"
)

// PaymentHandler handles payment-related HTTP requests for merchants
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler creates a new payment handler instance
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// SubmitPayment godoc
// @Summary Submit payment for an order
// @Description Merchant submits a payment for an order (token, partial, full, COD, etc.)
// @Tags Merchant Payments
// @Accept json
// @Produce json
// @Param mid path int true "Merchant ID"
// @Param oid path int true "Order ID"
// @Param request body dto.SubmitPaymentRequest true "Payment details"
// @Success 200 {object} dto.SubmitPaymentResponse
// @Failure 400 {object} response.ErrorDetails
// @Failure 401 {object} response.ErrorDetails
// @Failure 500 {object} response.ErrorDetails
// @Router /api/v0/merchant/{mid}/orders/{oid}/payments [post]
func (h *PaymentHandler) SubmitPayment(c *fiber.Ctx) error {
	// Parse merchant ID from path
	merchantID, err := strconv.ParseUint(c.Params("mid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid merchant ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	// Parse order ID from path
	orderID, err := strconv.ParseUint(c.Params("oid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid order ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	// Parse request body
	var req dto.SubmitPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid request body", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	// Validate required fields
	if req.Amount <= 0 {
		body := response.GetErrorHTTPResponseBody(400, "Amount must be greater than 0", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	if req.PaymentMode == 0 {
		body := response.GetErrorHTTPResponseBody(400, "Payment mode is required", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	// Call service
	result, err := h.paymentService.SubmitPayment(
		c.Context(),
		uint(merchantID),
		uint(orderID),
		&req,
	)

	if err != nil {
		body := response.GetErrorHTTPResponseBody(500, err.Error(), nil)
		return response.WriteHTTPResponse(c, http.StatusInternalServerError, body)
	}

	body := &response.HTTPResponse{
		Content: result,
	}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}

// GetOrderPayments godoc
// @Summary Get all payments for an order
// @Description Retrieve all payment records for a specific order
// @Tags Merchant Payments
// @Accept json
// @Produce json
// @Param mid path int true "Merchant ID"
// @Param oid path int true "Order ID"
// @Success 200 {object} dto.GetPaymentsResponse
// @Failure 400 {object} response.ErrorDetails
// @Failure 401 {object} response.ErrorDetails
// @Failure 500 {object} response.ErrorDetails
// @Router /api/v0/merchant/{mid}/orders/{oid}/payments [get]
func (h *PaymentHandler) GetOrderPayments(c *fiber.Ctx) error {
	// Parse merchant ID from path
	merchantID, err := strconv.ParseUint(c.Params("mid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid merchant ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	// Parse order ID from path
	orderID, err := strconv.ParseUint(c.Params("oid"), 10, 32)
	if err != nil {
		body := response.GetErrorHTTPResponseBody(400, "Invalid order ID", nil)
		return response.WriteHTTPResponse(c, http.StatusBadRequest, body)
	}

	// Call service
	result, err := h.paymentService.GetOrderPayments(
		c.Context(),
		uint(merchantID),
		uint(orderID),
	)

	if err != nil {
		body := response.GetErrorHTTPResponseBody(500, err.Error(), nil)
		return response.WriteHTTPResponse(c, http.StatusInternalServerError, body)
	}

	body := &response.HTTPResponse{
		Content: result,
	}
	return response.WriteHTTPResponse(c, http.StatusOK, body)
}
