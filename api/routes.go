package api

import (
	"linksupply.io/vmconnect/api/health"

	// Auth modules
	authHandler "linksupply.io/vmconnect/api/auth/handler"
	authRepository "linksupply.io/vmconnect/api/auth/repository"
	authService "linksupply.io/vmconnect/api/auth/service"
	authValidator "linksupply.io/vmconnect/api/auth/validator"
)

func SetupRoutes(server *APIServer) {
	// Setup common auth routes
	SetupCommonRoutes(server)

	/*------------------------------------Merchant-----------------------------------------*/
	SetupMerchantRoutes(server)
	// /*------------------------------------Vendor-----------------------------------------*/
	SetupVendorRoutes(server)
	/*------------------------------------Health-----------------------------------------*/
	SetupHealthRoutes(server)
}

// SetupCommonRoutes sets up common authentication routes
func SetupCommonRoutes(server *APIServer) {
	// Initialize auth dependencies
	authRepo := authRepository.NewAuthRepository(server.dataSource)
	authSvc := authService.NewAuthService(authRepo)
	authVal := authValidator.NewAuthValidator()
	authHand := authHandler.NewAuthHandler(authSvc, authVal)

	// Common authentication routes
	server.app.Post("/api/auth/login", authHand.Login)
	server.app.Post("/api/auth/generate-otp", authHand.GenerateOTP)
	server.app.Post("/api/auth/validate-otp", authHand.ValidateOTP)
}

func SetupHealthRoutes(server *APIServer) {
	// Setup healthcheck routes and default routes
	app := server.app
	/*------------------------------------Health-----------------------------------------*/
	app.Get("/health", health.Health)
}
