package api

import (
	"github.com/gofiber/fiber/v2"
	authRepository "linksupply.io/vmconnect/api/auth/repository"
	vendorHandler "linksupply.io/vmconnect/api/vendorapi/handler"
	vendorRepository "linksupply.io/vmconnect/api/vendorapi/repository"
	vendorService "linksupply.io/vmconnect/api/vendorapi/service"
	vendorValidator "linksupply.io/vmconnect/api/vendorapi/validator"
	"linksupply.io/vmconnect/middleware"
)

const enableDeferredVendorOrderFlow = false

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
	listRepo := vendorRepository.NewVendorMerchantListRepository(server.dataSource)
	svc := vendorService.NewVendorMerchantService(repo, listRepo)
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

	// Public vendor routes
	vendorApi.Post("/register", func(c *fiber.Ctx) error {
		vh := GetDefaultVendorHandler(server)
		return vh.Register(c)
	})

	vendorApi.Post("/:vid/orders", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.CreateOrder(c)
	})

	vendorApi.Get("/:vid/products/categories", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.GetCategories(c)
	})

	authRepo := authRepository.NewAuthRepository(server.dataSource)
	vendorProtected := vendorApi.Group("", middleware.VendorAPIProtected()...)
	vendorScoped := vendorProtected.Group("/:vid", middleware.RequireVendorOwnership(authRepo))

	// Get vendor info by ID
	vendorScoped.Get("", func(c *fiber.Ctx) error {
		vh := GetDefaultVendorHandler(server)
		return vh.GetVendorInfo(c)
	})

	// Product management routes
	vendorScoped.Post("/products", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.AddProduct(c)
	})

	vendorScoped.Put("/products/:pid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.UpdateProduct(c)
	})

	vendorScoped.Delete("/products/:pid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.DeleteProduct(c)
	})

	vendorScoped.Get("/products/:pid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.GetProduct(c)
	})

	vendorScoped.Get("/products", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorProductHandler(server)
		return ph.GetProducts(c)
	})

	// Merchant management routes
	vendorScoped.Post("/merchants", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.AddMerchant(c)
	})

	vendorScoped.Get("/merchants", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.GetMerchants(c)
	})

	vendorScoped.Put("/merchants/:mid", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.UpdateMerchant(c)
	})

	vendorScoped.Delete("/merchants/:mid", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.DeleteMerchant(c)
	})

	// Deferred until product plan for vendor-managed order flow is finalized.
	if enableDeferredVendorOrderFlow {
		registerDeferredVendorOrderFlowRoutes(vendorScoped, server)
	}
}

func registerDeferredVendorOrderFlowRoutes(vendorScoped fiber.Router, server *APIServer) {
	vendorScoped.Get("/merchants/:mid/visibility", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.GetProductVisibility(c)
	})

	vendorScoped.Patch("/merchants/:mid/visibility", func(c *fiber.Ctx) error {
		mh := GetDefaultVendorMerchantHandler(server)
		return mh.UpdateProductVisibility(c)
	})

	vendorScoped.Post("/orders/:oid/confirm", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.ConfirmOrder(c)
	})

	vendorScoped.Post("/orders/:oid/invoice", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.GenerateInvoice(c)
	})

	vendorScoped.Post("/orders/:oid/dispatch", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.DispatchOrder(c)
	})

	vendorScoped.Post("/orders/:oid/mark-paid", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorPaymentHandler(server)
		return ph.MarkOrderAsPaid(c)
	})

	vendorScoped.Get("/orders/:oid/payments", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorPaymentHandler(server)
		return ph.GetOrderPayments(c)
	})

	vendorScoped.Post("/orders/:oid/payments/:pid/verify", func(c *fiber.Ctx) error {
		ph := GetDefaultVendorPaymentHandler(server)
		return ph.VerifyPayment(c)
	})

	vendorScoped.Post("/orders/:oid/cancel", func(c *fiber.Ctx) error {
		oh := GetDefaultVendorOrderHandler(server)
		return oh.CancelOrder(c)
	})
}
