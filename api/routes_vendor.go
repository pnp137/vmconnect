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

func GetDefaultVendorProductHandler(server *APIServer) *vendorHandler.VendorProductHandler {
	// Initialize vendor product dependencies
	repo := vendorRepository.NewVendorProductRepository(server.dataSource)
	svc := vendorService.NewVendorProductService(repo)
	return vendorHandler.NewVendorProductHandler(svc)
}

func GetDefaultVendorMerchantHandler(server *APIServer) *vendorHandler.VendorMerchantHandler {
	// Initialize vendor merchant dependencies
	repo := vendorRepository.NewVendorMerchantRepository(server.dataSource)
	svc := vendorService.NewVendorMerchantService(repo)
	return vendorHandler.NewVendorMerchantHandler(svc)
}

func GetDefaultVendorOrderHandler(server *APIServer) *vendorHandler.OrderHandler {
	// Initialize vendor order dependencies
	repo := vendorRepository.NewVendorRepository(server.dataSource)
	orderSvc := vendorService.NewOrderService(repo)
	return vendorHandler.NewOrderHandler(orderSvc)
}

func GetDefaultVendorPaymentHandler(server *APIServer) *vendorHandler.PaymentHandler {
	// Initialize vendor payment dependencies
	repo := vendorRepository.NewVendorRepository(server.dataSource)
	paymentSvc := vendorService.NewPaymentService(repo)
	return vendorHandler.NewPaymentHandler(paymentSvc)
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

	// Get vendor info by ID
	vendorApi.Get("/:vid", func(c *fiber.Ctx) error {
		vh := GetDefaultVendorHandler(server)
		return vh.GetVendorInfo(c)
	})

	// Product management routes
	vendorApi.Post("/:vid/products", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.AddProduct(c)
	})

	vendorApi.Put("/:vid/products/:pid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.UpdateProduct(c)
	})

	vendorApi.Delete("/:vid/products/:pid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.DeleteProduct(c)
	})

	vendorApi.Get("/:vid/products/:pid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.GetProduct(c)
	})

	vendorApi.Get("/:vid/products", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.GetProducts(c)
	})

	vendorApi.Get("/:vid/products/categories", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.GetCategories(c)
	})

	// Merchant management routes
	vendorApi.Get("/:vid/merchants", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.GetMerchants(c)
	})

	vendorApi.Get("/:vid/merchants/:mid/visibility", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.GetProductVisibility(c)
	})

	vendorApi.Patch("/:vid/merchants/:mid/visibility", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.UpdateProductVisibility(c)
	})

	// Order state transition endpoints
	vendorApi.Post("/:vid/orders/:oid/confirm", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.ConfirmOrder(c)
	})

	vendorApi.Post("/:vid/orders/:oid/invoice", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.GenerateInvoice(c)
	})

	vendorApi.Post("/:vid/orders/:oid/dispatch", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.DispatchOrder(c)
	})

	// Payment endpoints (new flexible payment system)
	vendorApi.Post("/:vid/orders/:oid/mark-paid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorPaymentHandler(server)
		return ph.MarkOrderAsPaid(c)
	})

	vendorApi.Get("/:vid/orders/:oid/payments", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorPaymentHandler(server)
		return ph.GetOrderPayments(c)
	})

	vendorApi.Post("/:vid/orders/:oid/payments/:pid/verify", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorPaymentHandler(server)
		return ph.VerifyPayment(c)
	})

	vendorApi.Post("/:vid/orders/:oid/cancel", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.CancelOrder(c)
	})

}
