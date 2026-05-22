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

// GenerateOTP handles OTP generation requests.
func (h *AuthHandler) GenerateOTP(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "handler.GenerateOTP: "

	reqBody := &dto.GenerateOTPRequest{}
	if err := parseAndGetRequest(c, reqBody); err != nil {
		log.Error(err, logPrefix+"Error in parsing the request body")
		body := response.GetErrorHTTPResponseBody(400, "Invalid Request Body")
		return response.WriteHTTPResponse(c, 400, body)
	}

	log.Debug(logPrefix + "OTP generation request received for username: " + reqBody.Username)

	if err := h.validator.ValidateGenerateOTPRequest(reqBody); err != nil {
		log.Error(err, logPrefix+"Error in validating OTP generation request")
		body := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, body)
	}

	otpRequestDto := dto.ToGenerateOTPRequestDto(reqBody)
	otpResponse, errorDetails := h.service.GenerateOTP(c.UserContext(), otpRequestDto)
	if errorDetails != nil {
		log.Error(errorDetails.Error, logPrefix+"Error in OTP generation process")
		body := response.GetErrorHTTPResponseBody(errorDetails.Code, errorDetails.Message)
		return response.WriteHTTPResponse(c, errorDetails.Code, body)
	}

	log.Info(logPrefix + "OTP generated for username: " + reqBody.Username)

	body := &response.HTTPResponse{
		Content: otpResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// ValidateOTP handles OTP validation requests and returns a token response.
func (h *AuthHandler) ValidateOTP(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "handler.ValidateOTP: "

	reqBody := &dto.ValidateOTPRequest{}
	if err := parseAndGetRequest(c, reqBody); err != nil {
		log.Error(err, logPrefix+"Error in parsing the request body")
		body := response.GetErrorHTTPResponseBody(400, "Invalid Request Body")
		return response.WriteHTTPResponse(c, 400, body)
	}

	log.Debug(logPrefix + "OTP validation request received for username: " + reqBody.Username)

	if err := h.validator.ValidateValidateOTPRequest(reqBody); err != nil {
		log.Error(err, logPrefix+"Error in validating OTP validation request")
		body := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, body)
	}

	otpRequestDto := dto.ToValidateOTPRequestDto(reqBody)
	authResponse, errorDetails := h.service.ValidateOTP(c.UserContext(), otpRequestDto)
	if errorDetails != nil {
		log.Error(errorDetails.Error, logPrefix+"Error in OTP validation process")
		body := response.GetErrorHTTPResponseBody(errorDetails.Code, errorDetails.Message)
		return response.WriteHTTPResponse(c, errorDetails.Code, body)
	}

	log.Info(logPrefix + "OTP login successful for username: " + reqBody.Username)

	body := &response.HTTPResponse{
		Content: authResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// parseAndGetLoginRequest parses the request body into LoginRequest
func parseAndGetLoginRequest(c *fiber.Ctx, body *dto.LoginRequest) error {
	return parseAndGetRequest(c, body)
}

func parseAndGetRequest(c *fiber.Ctx, body interface{}) error {
	return json.Unmarshal([]byte(c.Body()), body)
}
