package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/utils/auth"
)

// JWTProtected is a middleware to protect routes with JWT authentication
func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			errorBody := response.GetErrorHTTPResponseBody(401, "Authorization header required")
			return response.WriteHTTPResponse(c, 401, errorBody)
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			errorBody := response.GetErrorHTTPResponseBody(401, "Invalid authorization format. Use 'Bearer <token>'")
			return response.WriteHTTPResponse(c, 401, errorBody)
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			errorBody := response.GetErrorHTTPResponseBody(401, "Token is required")
			return response.WriteHTTPResponse(c, 401, errorBody)
		}

		// Validate the token
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			errorBody := response.GetErrorHTTPResponseBody(401, "Invalid or expired token")
			return response.WriteHTTPResponse(c, 401, errorBody)
		}

		// Store user info in context for use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role_type", claims.RoleType)

		// Continue to the next handler
		return c.Next()
	}
}

// RequireRole is a middleware to check if the user has a specific role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role_type").(string)
		if !ok || userRole != requiredRole {
			errorBody := response.GetErrorHTTPResponseBody(403, "Insufficient permissions")
			return response.WriteHTTPResponse(c, 403, errorBody)
		}
		return c.Next()
	}
}

// RequireRoles allows access to users with any of the specified roles
func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role_type").(string)
		if !ok {
			errorBody := response.GetErrorHTTPResponseBody(403, "Role not found in token")
			return response.WriteHTTPResponse(c, 403, errorBody)
		}

		// Check if user role matches any of the required roles
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		errorBody := response.GetErrorHTTPResponseBody(403, "Insufficient permissions")
		return response.WriteHTTPResponse(c, 403, errorBody)
	}
}
