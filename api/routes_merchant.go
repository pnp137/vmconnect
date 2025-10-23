package api

import (
	"github.com/gofiber/fiber/v2"
	merchantHandler "linksupply.io/vmconnect/api/merchant/handler"
	merchantRepository "linksupply.io/vmconnect/api/merchant/repository"
	merchantService "linksupply.io/vmconnect/api/merchant/service"
	merchantValidator "linksupply.io/vmconnect/api/merchant/validator"
)

func GetDefaultMerchantHandler(server *APIServer) *merchantHandler.MerchantHandler {
	// Initialize merchant dependencies
	repo := merchantRepository.NewMerchantRepository(server.dataSource)
	svc := merchantService.NewMerchantService(repo)
	val := merchantValidator.NewMerchantValidator()
	hand := merchantHandler.NewMerchantHandler(svc, val)
	return hand
}

func SetupMerchantRoutes(server *APIServer) {
	// Merchant registration route

	app := server.app
	merchantApi := app.Group("/api/v0/merchant")

	merchantApi.Post("/register", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantHandler(server)
		return handler.Register(c)
	})
}
