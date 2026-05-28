package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

// MerchantRepository defines the interface for merchant repository operations
type MerchantRepository interface {
	// User operations
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmailOrPhone(email, phone string) (*models.User, error)

	// Merchant operations
	CreateMerchant(merchant *models.Merchant) (*models.Merchant, error)
	GetMerchantByID(merchantID uint) (*models.Merchant, error)
	GetMerchantByUserID(userID uint) (*models.Merchant, error)

	// Vendor operations
	GetAndValidateVendorByCode(vendorCode string) (*models.Vendor, error)
	CreateVendorMerchantMapping(mapping *models.VendorMerchantMapping) (*models.VendorMerchantMapping, error)
	GetMerchantVendors(merchantID uint) ([]models.Vendor, error)
	GetVendorProductsForMerchant(merchantID, vendorID uint) ([]models.Product, error)
	GetVendorCategoriesForMerchant(merchantID, vendorID uint) ([]models.Category, error)
	GetVendorProductsByCategoryForMerchant(merchantID, vendorID uint, queryParams *dto.ProductQueryParam) ([]models.Product, int, error)
	// Cart operations
	GetOrCreateCartOrder(merchantID, vendorID uint) (*models.Order, error)
	AddToCart(orderID, productID uint, quantity int, price float64) (*models.OrderItem, error)
	UpdateCartItemQuantity(orderItemID uint, quantity int) (*models.OrderItem, error)
	UpdateOrderStatus(orderID uint, status models.OrderStatus) error
	CreateOrderActivity(activity *models.OrderActivity) (*models.OrderActivity, error)
	GetMerchantOrders(merchantID uint, queryParams *dto.OrderQueryParam) ([]models.Order, int, error)
	GetProductByID(productID uint) (*models.Product, error)
	GetOrderByID(orderID uint) (*models.Order, error)
	GetOrderItemsByOrderID(orderID uint, orderItems *[]models.OrderItem) error
	UpdateOrderTotalAmount(orderID uint, totalAmount float64) error
	GetOrderItemByProductID(orderID, productID uint) (*models.OrderItem, error)
	GetOrderItemByID(orderItemID uint) (*models.OrderItem, error)
	SoftDeleteOrderItem(orderItemID uint) error
	HardDeleteOrderItem(orderItemID uint) error
	UpdateOrderItemCount(orderID uint) error
	// Order Activity operations
	CheckStatusInHistory(orderID uint, status models.OrderStatus) (bool, *models.OrderActivity, error)
	// Payment operations
	CreatePayment(payment *models.Payment) (*models.Payment, error)
	UpdatePayment(payment *models.Payment) error
	GetPaymentByID(paymentID uint64) (*models.Payment, error)
	GetOrderPayments(orderID uint) ([]models.Payment, error)
	RecalculateOrderPayments(orderID uint) error
}

type MerchantRepositoryImpl struct {
	db *gorm.DB
}

// NewMerchantRepository creates a new instance of merchant repository
func NewMerchantRepository(db *system.DataSource) MerchantRepository {
	return &MerchantRepositoryImpl{db: db.Db}
}

// CreateUser creates a new user
func (r *MerchantRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateMerchant creates a new merchant profile
func (r *MerchantRepositoryImpl) CreateMerchant(merchant *models.Merchant) (*models.Merchant, error) {
	if err := r.db.Create(merchant).Error; err != nil {
		return nil, err
	}
	return merchant, nil
}

// GetUserByEmailOrPhone checks if user exists with email or phone
func (r *MerchantRepositoryImpl) GetUserByEmailOrPhone(email, phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("(email = ? OR phone = ?) AND is_deleted = ?", email, phone, false).First(&user).Error
	return &user, err
}

// GetAndValidateVendorByCode validates vendor code and returns vendor if valid
func (r *MerchantRepositoryImpl) GetAndValidateVendorByCode(vendorCode string) (*models.Vendor, error) {
	if vendorCode == "" {
		return nil, gorm.ErrRecordNotFound // Empty vendor code is invalid
	}

	var vendor models.Vendor
	err := r.db.Where("vendor_code = ? AND is_deleted = ?", vendorCode, false).First(&vendor).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

// CreateVendorMerchantMapping creates a new vendor-merchant mapping
func (r *MerchantRepositoryImpl) CreateVendorMerchantMapping(mapping *models.VendorMerchantMapping) (*models.VendorMerchantMapping, error) {
	if err := r.db.Create(mapping).Error; err != nil {
		return nil, err
	}
	return mapping, nil
}

// GetMerchantVendors gets all vendors linked to a merchant
func (r *MerchantRepositoryImpl) GetMerchantVendors(merchantID uint) ([]models.Vendor, error) {
	var vendors []models.Vendor
	err := r.db.Table("vendors").
		Joins("JOIN vendor_merchant_mappings ON vendors.id = vendor_merchant_mappings.vendor_id").
		Where("vendor_merchant_mappings.merchant_id = ? AND vendor_merchant_mappings.is_deleted = ? AND vendors.is_deleted = ?",
			merchantID, false, false).
		Find(&vendors).Error
	return vendors, err
}

// GetVendorProductsForMerchant gets products for merchant-vendor pair with smart logic
func (r *MerchantRepositoryImpl) GetVendorProductsForMerchant(merchantID, vendorID uint) ([]models.Product, error) {
	var products []models.Product

	// First, try to get products with custom visibility
	err := r.db.Model(&models.Product{}).
		Preload("Category").
		Preload("Variants", "is_active = ?", true).
		Joins("JOIN merchant_product_visibilities ON products.id = merchant_product_visibilities.product_id").
		Where("merchant_product_visibilities.merchant_id = ? AND merchant_product_visibilities.vendor_id = ? AND merchant_product_visibilities.is_deleted = ?",
			merchantID, vendorID, false).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	// If no products found with custom visibility, get all vendor products
	if len(products) == 0 {
		err = r.db.Model(&models.Product{}).
			Preload("Category").
			Preload("Variants", "is_active = ?", true).
			Where("vendor_id = ?", vendorID).Find(&products).Error
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

// GetMerchantByID gets merchant by ID with user info preloaded
func (r *MerchantRepositoryImpl) GetMerchantByID(merchantID uint) (*models.Merchant, error) {
	var merchant models.Merchant
	err := r.db.Model(&models.Merchant{}).Preload("Owner").Where("id = ? AND is_deleted = ?", merchantID, false).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// GetMerchantByUserID gets merchant by user ID
func (r *MerchantRepositoryImpl) GetMerchantByUserID(userID uint) (*models.Merchant, error) {
	var merchant models.Merchant
	err := r.db.Where("user_id = ? AND is_deleted = ?", userID, false).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// GetVendorCategoriesForMerchant gets unique categories for merchant-vendor pair
func (r *MerchantRepositoryImpl) GetVendorCategoriesForMerchant(merchantID, vendorID uint) ([]models.Category, error) {
	var categories []models.Category

	// First, try to get categories from products with custom visibility
	err := r.db.Table("categories").
		Joins("JOIN products ON categories.id = products.category_id").
		Joins("JOIN merchant_product_visibilities ON products.id = merchant_product_visibilities.product_id").
		Where("merchant_product_visibilities.merchant_id = ? AND merchant_product_visibilities.vendor_id = ? AND merchant_product_visibilities.is_deleted = ? AND products.is_deleted = ? AND categories.is_deleted = ?",
			merchantID, vendorID, false, false, false).
		Distinct("categories.id, categories.vendor_id, categories.name, categories.is_deleted, categories.created_at, categories.updated_at").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	// If no categories found with custom visibility, get all vendor categories
	if len(categories) == 0 {
		err = r.db.Where("vendor_id = ? AND is_deleted = ?", vendorID, false).Find(&categories).Error
		if err != nil {
			return nil, err
		}
	}

	return categories, nil
}

// GetVendorProductsByCategoryForMerchant gets paginated products for specific category
func (r *MerchantRepositoryImpl) GetVendorProductsByCategoryForMerchant(merchantID, vendorID uint, queryParams *dto.ProductQueryParam) ([]models.Product, int, error) {
	var products []models.Product
	var totalCount int64

	// Build base query
	query := r.db.Model(&models.Product{}).
		Preload("Category").
		Preload("Variants", "is_active = ?", true).
		Joins("JOIN categories ON products.category_id = categories.id").
		Where("categories.id = ? AND categories.vendor_id = ? AND categories.is_deleted = ?",
			queryParams.CategoryID, vendorID, false)

	// Apply state filter
	if queryParams.State == "active" {
		query = query.Where("products.is_active = ?", true)
	} else if queryParams.State == "featured" {
		query = query.Where("products.is_featured = ?", true)
	}

	// First, try to get products with custom visibility
	customQuery := query.Joins("JOIN merchant_product_visibilities ON products.id = merchant_product_visibilities.product_id").
		Where("merchant_product_visibilities.merchant_id = ? AND merchant_product_visibilities.vendor_id = ? AND merchant_product_visibilities.is_deleted = ?",
			merchantID, vendorID, false)

	// Count total for custom visibility
	err := customQuery.Model(&models.Product{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// If no products found with custom visibility, get all vendor products for this category
	if totalCount == 0 {
		// Count total for all products
		err = query.Model(&models.Product{}).Count(&totalCount).Error
		if err != nil {
			return nil, 0, err
		}

		// Apply ordering
		switch queryParams.Ordering {
		case "name":
			query = query.Order("products.name ASC")
		case "price":
			query = query.Joins("LEFT JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).Order("pv.price ASC")
		case "created_at":
			query = query.Order("products.created_at DESC")
		default:
			query = query.Order("products.name ASC")
		}

		// Apply pagination
		err = query.Limit(queryParams.Limit).Offset(queryParams.Offset).Find(&products).Error
	} else {
		// Apply ordering for custom visibility
		switch queryParams.Ordering {
		case "name":
			customQuery = customQuery.Order("products.name ASC")
		case "price":
			customQuery = customQuery.Joins("LEFT JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).Order("pv.price ASC")
		case "created_at":
			customQuery = customQuery.Order("products.created_at DESC")
		default:
			customQuery = customQuery.Order("products.name ASC")
		}

		// Apply pagination for custom visibility
		err = customQuery.Limit(queryParams.Limit).Offset(queryParams.Offset).Find(&products).Error
	}

	if err != nil {
		return nil, 0, err
	}

	return products, int(totalCount), nil
}

// GetOrCreateCartOrder gets existing cart order or creates new one
func (r *MerchantRepositoryImpl) GetOrCreateCartOrder(merchantID, vendorID uint) (*models.Order, error) {
	var order models.Order

	// Try to find existing cart order
	err := r.db.Where("merchant_id = ? AND vendor_id = ? AND status = ? AND is_deleted = ?",
		merchantID, vendorID, models.ORDER_IN_CART, false).First(&order).Error

	if err == nil {
		return &order, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	orderStatus := models.ORDER_IN_CART
	// Create new cart order
	order = models.Order{
		MerchantID:  &merchantID,
		VendorID:    vendorID,
		Status:      &orderStatus,
		TotalAmount: 0,
	}

	err = r.db.Create(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// AddToCart adds product to cart or updates existing item
func (r *MerchantRepositoryImpl) AddToCart(orderID, productID uint, quantity int, price float64) (*models.OrderItem, error) {
	var orderItem models.OrderItem

	// Check if item already exists in cart
	err := r.db.Where("order_id = ? AND product_id = ? AND is_deleted = ?",
		orderID, productID, false).First(&orderItem).Error

	if err == nil {
		// Update existing item
		orderItem.Quantity += float64(quantity)
		orderItem.SubTotal = orderItem.Quantity * orderItem.Price

		err = r.db.Save(&orderItem).Error
		if err != nil {
			return nil, err
		}

		// Update order item count
		err = r.UpdateOrderItemCount(orderID)
		if err != nil {
			return nil, err
		}

		return &orderItem, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new order item
	orderItem = models.OrderItem{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  float64(quantity),
		Price:     price,
		SubTotal:  float64(quantity) * price,
	}

	err = r.db.Create(&orderItem).Error
	if err != nil {
		return nil, err
	}

	// Update order item count
	err = r.UpdateOrderItemCount(orderID)
	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}

// UpdateOrderItemCount updates the item_count field in the order table
func (r *MerchantRepositoryImpl) UpdateOrderItemCount(orderID uint) error {
	// Count active order items for this order
	var count int64
	err := r.db.Model(&models.OrderItem{}).
		Where("order_id = ? AND is_deleted = ?", orderID, false).
		Count(&count).Error
	if err != nil {
		return err
	}

	// Update the order's item_count field
	err = r.db.Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("item_count", count).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateCartItemQuantity updates quantity of cart item
func (r *MerchantRepositoryImpl) UpdateCartItemQuantity(orderItemID uint, quantity int) (*models.OrderItem, error) {
	var orderItem models.OrderItem

	err := r.db.Where("id = ? AND is_deleted = ?", orderItemID, false).First(&orderItem).Error
	if err != nil {
		return nil, err
	}

	orderItem.Quantity = float64(quantity)
	orderItem.SubTotal = orderItem.Quantity * orderItem.Price

	err = r.db.Save(&orderItem).Error
	if err != nil {
		return nil, err
	}

	// Update order item count
	err = r.UpdateOrderItemCount(orderItem.OrderID)
	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}

// UpdateOrderStatus updates order status
func (r *MerchantRepositoryImpl) UpdateOrderStatus(orderID uint, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// CreateOrderActivity creates order activity record
func (r *MerchantRepositoryImpl) CreateOrderActivity(activity *models.OrderActivity) (*models.OrderActivity, error) {
	err := r.db.Create(activity).Error
	if err != nil {
		return nil, err
	}
	return activity, nil
}

// GetMerchantOrders gets merchant orders with filtering and pagination
func (r *MerchantRepositoryImpl) GetMerchantOrders(merchantID uint, queryParams *dto.OrderQueryParam) ([]models.Order, int, error) {
	var orders []models.Order
	var totalCount int64

	// Build base query
	query := r.db.Model(&models.Order{}).Preload("Vendor").
		Where("merchant_id = ? AND is_deleted = ?", merchantID, false)

	// Apply vendor filter
	if queryParams.VendorID > 0 {
		query = query.Where("vendor_id = ?", queryParams.VendorID)
	}

	// Apply status filter (array of status strings)
	if len(queryParams.OrderStatus) > 0 {
		query = query.Where("status IN (?)", queryParams.OrderStatus)
	}

	// Count total
	err := query.Model(&models.Order{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply ordering
	switch queryParams.Ordering {
	case "created_at":
		query = query.Order("created_at DESC")
	case "updated_at":
		query = query.Order("updated_at DESC")
	case "total_amount":
		query = query.Order("total_amount DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	err = query.Limit(queryParams.Limit).Offset(queryParams.Offset).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, int(totalCount), nil
}

// GetProductByID gets product by ID
func (r *MerchantRepositoryImpl) GetProductByID(productID uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Variants", "is_active = ?", true).Where("id = ?", productID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetOrderByID gets order by ID
func (r *MerchantRepositoryImpl) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Model(&models.Order{}).Preload("Vendor").Preload("OrderItems.Product").
		Where("id = ? AND is_deleted = ?", orderID, false).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderItemsByOrderID gets order items by order ID
func (r *MerchantRepositoryImpl) GetOrderItemsByOrderID(orderID uint, orderItems *[]models.OrderItem) error {
	return r.db.Where("order_id = ? AND is_deleted = ?", orderID, false).Find(orderItems).Error
}

// UpdateOrderTotalAmount updates order total amount
func (r *MerchantRepositoryImpl) UpdateOrderTotalAmount(orderID uint, totalAmount float64) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("total_amount", totalAmount).Error
}

// GetOrderItemByProductID gets order item by order ID and product ID
func (r *MerchantRepositoryImpl) GetOrderItemByProductID(orderID, productID uint) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	err := r.db.Where("order_id = ? AND product_id = ? AND is_deleted = ?", orderID, productID, false).First(&orderItem).Error
	if err != nil {
		return nil, err
	}
	return &orderItem, nil
}

// GetOrderItemByID gets order item by ID
func (r *MerchantRepositoryImpl) GetOrderItemByID(orderItemID uint) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	err := r.db.Where("id = ? AND is_deleted = ?", orderItemID, false).First(&orderItem).Error
	if err != nil {
		return nil, err
	}
	return &orderItem, nil
}

// SoftDeleteOrderItem soft deletes an order item
func (r *MerchantRepositoryImpl) SoftDeleteOrderItem(orderItemID uint) error {
	// First get the order ID before deleting
	var orderItem models.OrderItem
	err := r.db.Where("id = ?", orderItemID).First(&orderItem).Error
	if err != nil {
		return err
	}

	// Soft delete the order item
	err = r.db.Model(&models.OrderItem{}).Where("id = ?", orderItemID).Update("is_deleted", true).Error
	if err != nil {
		return err
	}

	// Update order item count
	err = r.UpdateOrderItemCount(orderItem.OrderID)
	if err != nil {
		return err
	}

	return nil
}

// HardDeleteOrderItem permanently deletes an order item
func (r *MerchantRepositoryImpl) HardDeleteOrderItem(orderItemID uint) error {
	// First get the order ID before deleting
	var orderItem models.OrderItem
	err := r.db.Where("id = ?", orderItemID).First(&orderItem).Error
	if err != nil {
		return err
	}

	// Hard delete the order item
	err = r.db.Delete(&orderItem).Error
	if err != nil {
		return err
	}

	// Update order item count
	err = r.UpdateOrderItemCount(orderItem.OrderID)
	if err != nil {
		return err
	}

	return nil
}

// InitializeOrderItemCounts initializes item_count for all existing orders
// This is useful for migration or fixing existing data
func (r *MerchantRepositoryImpl) InitializeOrderItemCounts() error {
	var orders []models.Order
	err := r.db.Find(&orders).Error
	if err != nil {
		return err
	}

	for _, order := range orders {
		err = r.UpdateOrderItemCount(order.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckStatusInHistory checks if a status was previously recorded in order activity
func (r *MerchantRepositoryImpl) CheckStatusInHistory(orderID uint, status models.OrderStatus) (bool, *models.OrderActivity, error) {
	var activity models.OrderActivity

	err := r.db.Where(
		"order_id = ? AND order_state = ? AND is_deleted = ?",
		orderID,
		int(status),
		false,
	).Order("created_at DESC").First(&activity).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, nil // Not found in history
		}
		return false, nil, err // Database error
	}

	return true, &activity, nil
}

// CreatePayment creates a new payment record
func (r *MerchantRepositoryImpl) CreatePayment(payment *models.Payment) (*models.Payment, error) {
	if err := r.db.Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePayment updates an existing payment record
func (r *MerchantRepositoryImpl) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// GetPaymentByID retrieves a payment by ID
func (r *MerchantRepositoryImpl) GetPaymentByID(paymentID uint64) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("id = ? AND is_deleted = ?", paymentID, false).
		Preload("Order").
		Preload("SubmittedByUser").
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetOrderPayments retrieves all payments for an order
func (r *MerchantRepositoryImpl) GetOrderPayments(orderID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("order_id = ? AND is_deleted = ?", orderID, false).
		Order("created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// RecalculateOrderPayments recalculates total paid and outstanding amounts for an order
func (r *MerchantRepositoryImpl) RecalculateOrderPayments(orderID uint) error {
	// Get all non-failed, non-cancelled payments
	var payments []models.Payment
	err := r.db.Where("order_id = ? AND is_deleted = ? AND status NOT IN (?)",
		orderID, false, []int{int(models.PaymentStatusFailed), int(models.PaymentStatusCancelled)}).
		Find(&payments).Error
	if err != nil {
		return err
	}

	// Calculate total paid
	totalPaid := 0.0
	for _, payment := range payments {
		totalPaid += payment.Amount
	}

	// Update order
	var order models.Order
	err = r.db.First(&order, orderID).Error
	if err != nil {
		return err
	}

	order.PaidAmount = totalPaid
	order.OutstandingAmount = order.TotalAmount - totalPaid

	return r.db.Save(&order).Error
}
