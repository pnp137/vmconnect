package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	authRepository "linksupply.io/vmconnect/api/auth/repository"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/database/models"
)

// VendorAPIProtected returns middleware that validates JWT and requires vendor role.
func VendorAPIProtected() []fiber.Handler {
	return []fiber.Handler{
		JWTProtected(),
		RequireRole(models.ROLE_VENDOR.String()),
	}
}

// RequireVendorOwnership ensures the authenticated vendor can only access their own :vid routes.
func RequireVendorOwnership(authRepo authRepository.AuthRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vendorIDStr := c.Params("vid")
		requestedVendorID, err := strconv.ParseUint(vendorIDStr, 10, 64)
		if err != nil || requestedVendorID == 0 {
			errorBody := response.GetErrorHTTPResponseBody(400, "Invalid vendor ID")
			return response.WriteHTTPResponse(c, 400, errorBody)
		}

		userID, ok := c.Locals("user_id").(uint)
		if !ok || userID == 0 {
			errorBody := response.GetErrorHTTPResponseBody(401, "User not authenticated")
			return response.WriteHTTPResponse(c, 401, errorBody)
		}

		vendor, err := authRepo.GetVendorByUserID(userID)
		if err != nil {
			errorBody := response.GetErrorHTTPResponseBody(403, "Vendor profile not found")
			return response.WriteHTTPResponse(c, 403, errorBody)
		}

		if vendor.ID != uint(requestedVendorID) {
			errorBody := response.GetErrorHTTPResponseBody(403, "Access denied: vendor ID mismatch")
			return response.WriteHTTPResponse(c, 403, errorBody)
		}

		c.Locals("vendor_id", vendor.ID)
		return c.Next()
	}
}
