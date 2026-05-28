package handler

import (
	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/service"
	validator "linksupply.io/vmconnect/api/vendorapi/validator"
	"linksupply.io/vmconnect/logger"
)

type VendorMerchantHandler struct {
	merchantService service.VendorMerchantService
	validator       validator.VendorValidator
}

func NewVendorMerchantHandler(merchantService service.VendorMerchantService) *VendorMerchantHandler {
	return &VendorMerchantHandler{
		merchantService: merchantService,
		validator:       validator.NewVendorValidator(),
	}
}

// AddMerchant handles POST /vendors/{vid}/merchants
func (h *VendorMerchantHandler) AddMerchant(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.AddMerchant] "

	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	var req dto.AddMerchantRequest
	if err := c.BodyParser(&req); err != nil {
		log.Debug(logPrefix + "Invalid request body: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid request body")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	if err := h.validator.ValidateAddMerchantRequest(&req); err != nil {
		log.Debug(logPrefix + "Invalid add merchant request: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	merchantResponse, errDetails := h.merchantService.AddMerchant(c.UserContext(), vendorID, dto.ToAddMerchantRequestDto(&req))
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error adding merchant: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	body := &response.HTTPResponse{
		Content: merchantResponse,
	}
	return response.WriteHTTPResponse(c, 201, body)
}

// GetMerchants handles GET /vendors/{vid}/merchants
func (h *VendorMerchantHandler) GetMerchants(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.GetMerchants] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	var queryParams dto.MerchantQueryParam
	if err := c.QueryParser(&queryParams); err != nil {
		log.Debug(logPrefix + "Invalid query parameters: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid query parameters")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}
	if queryParams.Ordering == "" {
		queryParams.Ordering = "created_at"
	}

	log.Debug(logPrefix + "Getting merchants for vendor ID: " + vendorIDStr)

	merchantsResponse, errDetails := h.merchantService.GetMerchants(c.UserContext(), vendorID, &queryParams)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting merchants: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved merchants for vendor ID: " + vendorIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: merchantsResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// UpdateMerchant handles PUT /vendors/{vid}/merchants/{mid}
func (h *VendorMerchantHandler) UpdateMerchant(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.UpdateMerchant] "

	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	var req dto.UpdateMerchantRequest
	if err := c.BodyParser(&req); err != nil {
		log.Debug(logPrefix + "Invalid request body: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid request body")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	if err := h.validator.ValidateUpdateMerchantRequest(&req); err != nil {
		log.Debug(logPrefix + "Invalid update merchant request: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	merchantResponse, errDetails := h.merchantService.UpdateMerchant(c.UserContext(), vendorID, merchantID, dto.ToUpdateMerchantRequestDto(&req))
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error updating merchant: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	body := &response.HTTPResponse{
		Content: merchantResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// DeleteMerchant handles DELETE /vendors/{vid}/merchants/{mid}
func (h *VendorMerchantHandler) DeleteMerchant(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.DeleteMerchant] "

	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	deleteResponse, errDetails := h.merchantService.DeleteMerchant(c.UserContext(), vendorID, merchantID)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error deleting merchant: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	body := &response.HTTPResponse{
		Content: deleteResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// GetProductVisibility handles GET /vendors/{vid}/merchants/{mid}/visibility
func (h *VendorMerchantHandler) GetProductVisibility(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.GetProductVisibility] "

	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	visibilityResponse, errDetails := h.merchantService.GetProductVisibility(c.UserContext(), vendorID, merchantID)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting product visibility: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	body := &response.HTTPResponse{
		Content: visibilityResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// UpdateProductVisibility handles PATCH /vendors/{vid}/merchants/{mid}/visibility
func (h *VendorMerchantHandler) UpdateProductVisibility(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.UpdateProductVisibility] "

	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	var req []dto.UpdateProductVisibilityRequest
	if err := c.BodyParser(&req); err != nil {
		log.Debug(logPrefix + "Invalid request body: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid request body")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Updating product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	updateResponse, errDetails := h.merchantService.UpdateProductVisibility(c.UserContext(), vendorID, merchantID, req)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error updating product visibility: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully updated product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	body := &response.HTTPResponse{
		Content: updateResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}
