package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/service"
	validator "linksupply.io/vmconnect/api/merchant/validator"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/logger"
)

type MerchantVendorHandler struct {
	vendorService service.MerchantVendorService
	validator     validator.MerchantValidator
}

// NewMerchantVendorHandler creates a new instance of merchant vendor handler
func NewMerchantVendorHandler(vendorService service.MerchantVendorService) *MerchantVendorHandler {
	return &MerchantVendorHandler{
		vendorService: vendorService,
		validator:     validator.NewMerchantValidator(),
	}
}

// GetMerchantVendors handles GET /merchants/{mid}/vendors
func (h *MerchantVendorHandler) GetMerchantVendors(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[MerchantVendorHandler.GetMerchantVendors] "

	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting vendors for merchant ID: " + merchantIDStr)

	// Call service
	vendorsResponse, errDetails := h.vendorService.GetMerchantVendors(c.UserContext(), uint(merchantID))
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting merchant vendors: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved " + strconv.Itoa(vendorsResponse.Count) + " vendors for merchant ID: " + merchantIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: vendorsResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// AddVendorToMerchant handles POST /merchants/{mid}/vendors
func (h *MerchantVendorHandler) AddVendorToMerchant(c *fiber.Ctx) error {
	logPrefix := "[MerchantVendorHandler.AddVendorToMerchant] "
	log := logger.NewLogger()
	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Parse request body
	var reqBody dto.AddVendorRequest
	if err := c.BodyParser(&reqBody); err != nil {
		log.Debug(logPrefix + "Invalid request body: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid request body")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Validate request body
	if err := h.validator.ValidateAddVendorRequest(&reqBody); err != nil {
		log.Debug(logPrefix + "Validation error: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Adding vendor with code: " + reqBody.VendorCode + " to merchant ID: " + merchantIDStr)

	// Call service
	addResponse, errDetails := h.vendorService.AddVendorToMerchant(c.UserContext(), uint(merchantID), &reqBody)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error adding vendor to merchant: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully added vendor " + reqBody.VendorCode + " to merchant ID: " + merchantIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: addResponse,
	}
	return response.WriteHTTPResponse(c, 201, body)
}

// GetVendorProducts handles GET /merchants/{mid}/vendors/{vid}/products
func (h *MerchantVendorHandler) GetVendorProducts(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[MerchantVendorHandler.GetVendorProducts] "

	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting products for merchant ID: " + merchantIDStr + " and vendor ID: " + vendorIDStr)

	// Call service
	productsResponse, errDetails := h.vendorService.GetVendorProducts(c.UserContext(), uint(merchantID), uint(vendorID))
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting vendor products: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved " + strconv.Itoa(productsResponse.Count) + " products for merchant ID: " + merchantIDStr + " and vendor ID: " + vendorIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: productsResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// GetVendorCategories handles GET /merchants/{mid}/vendors/{vid}/categories
func (h *MerchantVendorHandler) GetVendorCategories(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[MerchantVendorHandler.GetVendorCategories] "

	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting categories for merchant ID: " + merchantIDStr + " and vendor ID: " + vendorIDStr)

	// Call service
	categoriesResponse, errDetails := h.vendorService.GetVendorCategories(c.UserContext(), uint(merchantID), uint(vendorID))
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting vendor categories: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved " + strconv.Itoa(categoriesResponse.Count) + " categories for merchant ID: " + merchantIDStr + " and vendor ID: " + vendorIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: categoriesResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// GetVendorProductsByCategory handles GET /merchants/{mid}/vendors/{vid}/products/category
func (h *MerchantVendorHandler) GetVendorProductsByCategory(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[MerchantVendorHandler.GetVendorProductsByCategory] "

	// Get merchant ID from URL parameter
	merchantIDStr := c.Params("mid")
	merchantID, err := h.validator.ValidateIDParameter(merchantIDStr, "merchant ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid merchant ID: " + merchantIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Parse query parameters
	var queryParams dto.ProductQueryParam
	if err := c.QueryParser(&queryParams); err != nil {
		log.Debug(logPrefix + "Invalid query parameters: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid query parameters")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Set default values
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}
	if queryParams.Ordering == "" {
		queryParams.Ordering = "name"
	}

	// Validate query parameters
	if err := h.validator.ValidateProductQueryParam(&queryParams); err != nil {
		log.Debug(logPrefix + "Validation error: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting products for merchant ID: " + merchantIDStr + ", vendor ID: " + vendorIDStr + ", category ID: " + strconv.FormatUint(uint64(queryParams.CategoryID), 10))

	// Call service
	productsResponse, errDetails := h.vendorService.GetVendorProductsByCategory(c.UserContext(), uint(merchantID), uint(vendorID), &queryParams)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting vendor products by category: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved " + strconv.Itoa(productsResponse.Count) + " products for merchant ID: " + merchantIDStr + ", vendor ID: " + vendorIDStr + ", category ID: " + strconv.FormatUint(uint64(queryParams.CategoryID), 10))

	// Return success response
	body := &response.HTTPResponse{
		Content: productsResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}
