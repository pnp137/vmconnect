package api

import (
	"github.com/gofiber/fiber/v2"
	merchantHandler "linksupply.io/vmconnect/api/merchant/handler"
	merchantRepository "linksupply.io/vmconnect/api/merchant/repository"
	merchantService "linksupply.io/vmconnect/api/merchant/service"
	merchantValidator "linksupply.io/vmconnect/api/merchant/validator"
)

func GetDefaultMerchantVendorHandler(server *APIServer) *merchantHandler.MerchantVendorHandler {
	// Initialize merchant vendor dependencies
	repo := merchantRepository.NewMerchantRepository(server.dataSource)
	vendorSvc := merchantService.NewMerchantVendorService(repo)
	return merchantHandler.NewMerchantVendorHandler(vendorSvc)
}

func GetDefaultMerchantHandler(server *APIServer) *merchantHandler.MerchantHandler {
	// Initialize merchant dependencies
	repo := merchantRepository.NewMerchantRepository(server.dataSource)
	svc := merchantService.NewMerchantService(repo)
	val := merchantValidator.NewMerchantValidator()
	hand := merchantHandler.NewMerchantHandler(svc, val)
	return hand
}

func GetDefaultOrderHandler(server *APIServer) *merchantHandler.OrderHandler {
	// Initialize order dependencies
	repo := merchantRepository.NewMerchantRepository(server.dataSource)
	orderSvc := merchantService.NewOrderService(repo)
	return merchantHandler.NewOrderHandler(orderSvc)
}

func GetDefaultPaymentHandler(server *APIServer) *merchantHandler.PaymentHandler {
	// Initialize payment dependencies
	repo := merchantRepository.NewMerchantRepository(server.dataSource)
	paymentSvc := merchantService.NewPaymentService(repo)
	return merchantHandler.NewPaymentHandler(paymentSvc)
}

func SetupMerchantRoutes(server *APIServer) {
	app := server.app
	merchantApi := app.Group("/api/v0/merchant")

	// Merchant registration route
	merchantApi.Post("/register", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantHandler(server)
		return handler.Register(c)
	})

	// Get merchant info by ID
	merchantApi.Get("/:mid", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantHandler(server)
		return handler.GetMerchantInfo(c)
	})

	// Merchant vendor routes
	merchantApi.Get("/:mid/vendors", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantVendorHandler(server)
		return handler.GetMerchantVendors(c)
	})

	merchantApi.Post("/:mid/vendors", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantVendorHandler(server)
		return handler.AddVendorToMerchant(c)
	})

	merchantApi.Get("/:mid/vendors/:vid/products", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantVendorHandler(server)
		return handler.GetVendorProducts(c)
	})

	// Get vendor categories for merchant-vendor pair
	merchantApi.Get("/:mid/vendors/:vid/categories", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantVendorHandler(server)
		return handler.GetVendorCategories(c)
	})

	// Get paginated products by category for merchant-vendor pair
	merchantApi.Get("/:mid/vendors/:vid/products/category", func(c *fiber.Ctx) error {
		handler := GetDefaultMerchantVendorHandler(server)
		return handler.GetVendorProductsByCategory(c)
	})

	// Cart operations - New consolidated approach
	merchantApi.Post("/:mid/vendors/:vid/cart/items", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.ManageCartItem(c)
	})

	// Order operations
	merchantApi.Get("/:mid/orders", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.GetMerchantOrders(c)
	})

	merchantApi.Put("/:mid/orders/:oid/status", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.UpdateOrderStatus(c)
	})

	// Get order detail
	merchantApi.Get("/:mid/orders/:oid", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.GetOrderDetail(c)
	})

	// Delete cart item
	merchantApi.Delete("/:mid/cart/:oid/items/:itemId", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.DeleteOrderItem(c)
	})

	// State transition endpoints
	merchantApi.Post("/:mid/orders/:oid/place", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.PlaceOrder(c)
	})

	// Payment endpoints (new flexible payment system)
	merchantApi.Post("/:mid/orders/:oid/payments", func(c *fiber.Ctx) error {
		handler := GetDefaultPaymentHandler(server)
		return handler.SubmitPayment(c)
	})

	merchantApi.Get("/:mid/orders/:oid/payments", func(c *fiber.Ctx) error {
		handler := GetDefaultPaymentHandler(server)
		return handler.GetOrderPayments(c)
	})

	// Order state endpoints
	merchantApi.Post("/:mid/orders/:oid/receive", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.MarkReceived(c)
	})

	merchantApi.Post("/:mid/orders/:oid/complete", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.CompleteOrder(c)
	})

	merchantApi.Post("/:mid/orders/:oid/cancel", func(c *fiber.Ctx) error {
		handler := GetDefaultOrderHandler(server)
		return handler.CancelOrder(c)
	})
}
