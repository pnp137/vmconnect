package handler

import (
	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/service"
	"linksupply.io/vmconnect/api/vendorapi/validator"
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
