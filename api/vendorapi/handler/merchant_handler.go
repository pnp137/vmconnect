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

	log.Debug(logPrefix + "Getting merchants for vendor ID: " + vendorIDStr)

	// Call service
	merchantsResponse, errDetails := h.merchantService.GetMerchants(c.UserContext(), vendorID)
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

// GetProductVisibility handles GET /vendors/{vid}/merchants/{mid}/visibility
func (h *VendorMerchantHandler) GetProductVisibility(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.GetProductVisibility] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	// Call service
	visibilityResponse, errDetails := h.merchantService.GetProductVisibility(c.UserContext(), vendorID, merchantID)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting product visibility: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: visibilityResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// UpdateProductVisibility handles PATCH /vendors/{vid}/merchants/{mid}/visibility
func (h *VendorMerchantHandler) UpdateProductVisibility(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorMerchantHandler.UpdateProductVisibility] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Parse request body
	var req []dto.UpdateProductVisibilityRequest
	if err := c.BodyParser(&req); err != nil {
		log.Debug(logPrefix + "Invalid request body: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid request body")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Updating product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	// Call service
	updateResponse, errDetails := h.merchantService.UpdateProductVisibility(c.UserContext(), vendorID, merchantID, req)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error updating product visibility: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully updated product visibility for vendor ID: " + vendorIDStr + ", merchant ID: " + merchantIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: updateResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}
