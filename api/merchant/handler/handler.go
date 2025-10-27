package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/service"
	"linksupply.io/vmconnect/api/merchant/validator"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/logger"
)

// MerchantHandler handles merchant-specific requests
type MerchantHandler struct {
	service   service.MerchantService
	validator validator.MerchantValidator
}

// NewMerchantHandler creates a new instance of merchant handler
func NewMerchantHandler(service service.MerchantService, validator validator.MerchantValidator) *MerchantHandler {
	return &MerchantHandler{
		service:   service,
		validator: validator,
	}
}

// Register handles merchant registration
func (h *MerchantHandler) Register(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "handler.MerchantRegister: "

	// Parse request body
	reqBody := &dto.MerchantRegisterRequest{}
	if err := parseAndGetRegisterRequest(c, reqBody); err != nil {
		log.Error(err, logPrefix+"Error in parsing the request body")
		body := response.GetErrorHTTPResponseBody(400, "Invalid Request Body")
		return response.WriteHTTPResponse(c, 400, body)
	}

	log.Debug(logPrefix + "Merchant registration request received for email: " + reqBody.Email)

	// Validate registration request
	if err := h.validator.ValidateRegisterRequest(reqBody); err != nil {
		log.Error(err, logPrefix+"Error in validating registration request")
		body := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, body)
	}

	// Convert to DTO
	registerRequestDto := dto.ToMerchantRegisterRequestDto(reqBody)

	// Process registration
	merchantResponse, errorDetails := h.service.RegisterMerchant(c.UserContext(), registerRequestDto)
	if errorDetails != nil {
		log.Error(errorDetails.Error, logPrefix+"Error in merchant registration")
		body := response.GetErrorHTTPResponseBody(errorDetails.Code, errorDetails.Message)
		return response.WriteHTTPResponse(c, errorDetails.Code, body)
	}

	log.Info(logPrefix + "Merchant registration successful for email: " + reqBody.Email)

	// Return success response
	body := &response.HTTPResponse{
		Content: merchantResponse,
	}
	return response.WriteHTTPResponse(c, 201, body)
}

// GetMerchantInfo handles GET /merchants/{mid}
func (h *MerchantHandler) GetMerchantInfo(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[MerchantHandler.GetMerchantInfo] "

	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting merchant info for ID: " + merchantIDStr)

	// Call service
	merchantInfo, errDetails := h.service.GetMerchantInfo(c.UserContext(), uint(merchantID))
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting merchant info: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved merchant info for ID: " + merchantIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: merchantInfo,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// parseAndGetRegisterRequest parses the request body into MerchantRegisterRequest
func parseAndGetRegisterRequest(c *fiber.Ctx, body *dto.MerchantRegisterRequest) error {
	return json.Unmarshal([]byte(c.Body()), body)
}
