package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/auth/dto"
	"linksupply.io/vmconnect/api/auth/service"
	"linksupply.io/vmconnect/api/auth/validator"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/logger"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	service   service.AuthService
	validator validator.AuthValidator
}

// NewAuthHandler creates a new instance of auth handler
func NewAuthHandler(service service.AuthService, validator validator.AuthValidator) *AuthHandler {
	return &AuthHandler{
		service:   service,
		validator: validator,
	}
}

// Login handles login requests for both merchants and vendors
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "handler.Login: "

	// Parse request body
	reqBody := &dto.LoginRequest{}
	if err := parseAndGetLoginRequest(c, reqBody); err != nil {
		log.Error(err, logPrefix+"Error in parsing the request body")
		body := response.GetErrorHTTPResponseBody(400, "Invalid Request Body")
		return response.WriteHTTPResponse(c, 400, body)
	}

	log.Debug(logPrefix + "Login request received for username: " + reqBody.Username)

	// Validate login request
	if err := h.validator.ValidateLoginRequest(reqBody); err != nil {
		log.Error(err, logPrefix+"Error in validating login request")
		body := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, body)
	}

	// Convert to DTO
	loginRequestDto := dto.ToLoginRequestDto(reqBody)

	// Process login
	authResponse, errorDetails := h.service.Login(c.UserContext(), loginRequestDto)
	if errorDetails != nil {
		log.Error(errorDetails.Error, logPrefix+"Error in login process")
		body := response.GetErrorHTTPResponseBody(errorDetails.Code, errorDetails.Message)
		return response.WriteHTTPResponse(c, errorDetails.Code, body)
	}

	log.Info(logPrefix + "Login successful for username: " + reqBody.Username)

	// Return success response
	body := &response.HTTPResponse{
		Content: authResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// parseAndGetLoginRequest parses the request body into LoginRequest
func parseAndGetLoginRequest(c *fiber.Ctx, body *dto.LoginRequest) error {
	return json.Unmarshal([]byte(c.Body()), body)
}
