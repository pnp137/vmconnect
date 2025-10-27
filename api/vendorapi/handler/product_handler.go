package handler

import (
	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/service"
	validator "linksupply.io/vmconnect/api/vendorapi/validator"
	"linksupply.io/vmconnect/logger"
)

type VendorProductHandler struct {
	productService service.VendorProductService
	validator      validator.VendorValidator
}

func NewVendorProductHandler(productService service.VendorProductService) *VendorProductHandler {
	return &VendorProductHandler{
		productService: productService,
		validator:      validator.NewVendorValidator(),
	}
}

// AddProduct handles POST /vendors/{vid}/products
func (h *VendorProductHandler) AddProduct(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorProductHandler.AddProduct] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Parse request body
	var reqBody dto.AddProductRequest
	if err := c.BodyParser(&reqBody); err != nil {
		log.Debug(logPrefix + "Invalid request body: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid request body")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Adding product for vendor ID: " + vendorIDStr)

	// Call service
	addResponse, errDetails := h.productService.AddProduct(c.UserContext(), vendorID, &reqBody)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error adding product: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully added product for vendor ID: " + vendorIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: addResponse,
	}
	return response.WriteHTTPResponse(c, 201, body)
}

// UpdateProduct handles PUT /vendors/{vid}/products/{pid}
func (h *VendorProductHandler) UpdateProduct(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorProductHandler.UpdateProduct] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get product ID from URL parameter
	productIDStr := c.Params("pid")
	productID, err := h.validator.ValidateIDParameter(productIDStr, "product ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid product ID: " + productIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Parse request body
	var reqBody dto.UpdateProductRequest
	if err := c.BodyParser(&reqBody); err != nil {
		log.Debug(logPrefix + "Invalid request body: " + err.Error())
		errorBody := response.GetErrorHTTPResponseBody(400, "Invalid request body")
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Updating product ID: " + productIDStr + " for vendor ID: " + vendorIDStr)

	// Call service
	updateResponse, errDetails := h.productService.UpdateProduct(c.UserContext(), vendorID, productID, &reqBody)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error updating product: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully updated product ID: " + productIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: updateResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// DeleteProduct handles DELETE /vendors/{vid}/products/{pid}
func (h *VendorProductHandler) DeleteProduct(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorProductHandler.DeleteProduct] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get product ID from URL parameter
	productIDStr := c.Params("pid")
	productID, err := h.validator.ValidateIDParameter(productIDStr, "product ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid product ID: " + productIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Deleting product ID: " + productIDStr + " for vendor ID: " + vendorIDStr)

	// Call service
	deleteResponse, errDetails := h.productService.DeleteProduct(c.UserContext(), vendorID, productID)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error deleting product: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully deleted product ID: " + productIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: deleteResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// GetProduct handles GET /vendors/{vid}/products/{pid}
func (h *VendorProductHandler) GetProduct(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorProductHandler.GetProduct] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	// Get product ID from URL parameter
	productIDStr := c.Params("pid")
	productID, err := h.validator.ValidateIDParameter(productIDStr, "product ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid product ID: " + productIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting product ID: " + productIDStr + " for vendor ID: " + vendorIDStr)

	// Call service
	productResponse, errDetails := h.productService.GetProduct(c.UserContext(), vendorID, productID)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting product: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved product ID: " + productIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: productResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// GetProducts handles GET /vendors/{vid}/products
func (h *VendorProductHandler) GetProducts(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorProductHandler.GetProducts] "

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
		queryParams.Ordering = "created_at"
	}

	log.Debug(logPrefix + "Getting products for vendor ID: " + vendorIDStr)

	// Call service
	productsResponse, errDetails := h.productService.GetProducts(c.UserContext(), vendorID, &queryParams)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting products: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved products for vendor ID: " + vendorIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: productsResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}

// GetCategories handles GET /vendors/{vid}/products/categories
func (h *VendorProductHandler) GetCategories(c *fiber.Ctx) error {
	log := logger.NewLogger()
	logPrefix := "[VendorProductHandler.GetCategories] "

	// Get vendor ID from URL parameter
	vendorIDStr := c.Params("vid")
	vendorID, err := h.validator.ValidateIDParameter(vendorIDStr, "vendor ID")
	if err != nil {
		log.Debug(logPrefix + "Invalid vendor ID: " + vendorIDStr)
		errorBody := response.GetErrorHTTPResponseBody(400, err.Error())
		return response.WriteHTTPResponse(c, 400, errorBody)
	}

	log.Debug(logPrefix + "Getting categories for vendor ID: " + vendorIDStr)

	// Call service
	categoriesResponse, errDetails := h.productService.GetCategories(c.UserContext(), vendorID)
	if errDetails != nil {
		log.Error(errDetails.Error, logPrefix+"Error getting categories: "+errDetails.Message)
		errorBody := response.GetErrorHTTPResponseBody(errDetails.Code, errDetails.Message)
		return response.WriteHTTPResponse(c, errDetails.Code, errorBody)
	}

	log.Info(logPrefix + "Successfully retrieved categories for vendor ID: " + vendorIDStr)

	// Return success response
	body := &response.HTTPResponse{
		Content: categoriesResponse,
	}
	return response.WriteHTTPResponse(c, 200, body)
}
