package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/service"
	"linksupply.io/vmconnect/api/vendorapi/validator"
	"linksupply.io/vmconnect/logger"
)

type VendorHandler struct {
	vendorService   service.VendorService
	vendorValidator validator.VendorValidator
}

func NewVendorHandler(vendorService service.VendorService, validator validator.VendorValidator) *VendorHandler {
	return &VendorHandler{
		vendorService:   vendorService,
		vendorValidator: validator,
	}
}

// Register handles vendor registration
func (h *VendorHandler) Register(c *fiber.Ctx) error {
	var req dto.VendorRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid JSON format")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Validate request
	if err := h.vendorValidator.ValidateRegisterRequest(&req); err != nil {
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Convert to internal DTO
	reqDto := dto.ToVendorRegisterRequestDto(&req)

	// Register vendor
	vendorResponse, err := h.vendorService.RegisterVendor(c.UserContext(), reqDto)
	if err != nil {
		errorBody := response.GetErrorHTTPResponseBody(err.Code, err.Message)
		return response.WriteHTTPResponse(c, err.Code, errorBody)
	}

	// Return success response
	body := &response.HTTPResponse{
		Content: vendorResponse,
	}
	return response.WriteHTTPResponse(c, 201, body)
}

// GetVendorInfo handles GET /vendors/{vid}
func (h *VendorHandler) GetVendorInfo(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorHandler.GetVendorInfo] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := strconv.ParseUint(vendorIDStr, 10, 32)
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid vendor ID")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting vendor info for ID: " + vendorIDStr)

	// Call service
	vendorInfo, errDetails := h.vendorService.GetVendorInfo(c.UserContext(), uint(vendorID))
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting vendor info: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved vendor info for ID: " + vendorIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: vendorInfo,
	}
	return response.WriteHTTPResponse(c, 200, body)
}
