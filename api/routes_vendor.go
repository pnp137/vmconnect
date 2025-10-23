package api

import (
	"github.com/gofiber/fiber/v2"
	vendorHandler "linksupply.io/vmconnect/api/vendorapi/handler"
	vendorRepository "linksupply.io/vmconnect/api/vendorapi/repository"
	vendorService "linksupply.io/vmconnect/api/vendorapi/service"
	vendorValidator "linksupply.io/vmconnect/api/vendorapi/validator"
)

func GetDefaultVendorHandler(server *APIServer) *vendorHandler.VendorHandler {
	// Initialize vendor dependencies
	repo := vendorRepository.NewVendorRepository(server.dataSource)
	svc := vendorService.NewVendorService(repo)
	val := vendorValidator.NewVendorValidator()
	vh := vendorHandler.NewVendorHandler(svc, val)
	return vh
}

// SetupVendorRoutes wires vendor dependencies and registers vendor routes
func SetupVendorRoutes(server *APIServer) {

	app := server.app
	vendorApi := app.Group("/api/v0/vendor")

	// Vendor authentication routes
	vendorApi.Post("/register", func(c *fiber.Ctx) error {
		vh := GetDefaultVendorHandler(server)
		return vh.Register(c)
	})

}
