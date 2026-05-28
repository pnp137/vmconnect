package repository

import (
	"errors"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

// VendorRepository defines the interface for vendor repository
// VendorRepository defines the interface for vendor repository operations
type VendorRepository interface {
	GetDB() *gorm.DB

	// User operations
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmailOrPhone(email, phone string) (*models.User, error)
	// Vendor operations
	CreateVendor(vendor *models.Vendor) (*models.Vendor, error)
	GetVendorByID(vendorID uint) (*models.Vendor, error)

	// Order operations
	GetProductVariantsByIDs(variantIDs []uint) ([]models.ProductVariant, error)
	CreateOrder(tx *gorm.DB, order *models.Order) error
	CreateOrderItem(tx *gorm.DB, item *models.OrderItem) error
	ReduceProductVariantStock(tx *gorm.DB, variantID uint, quantity int) error
	GetOrderByID(orderID uint) (*models.Order, error)
	UpdateOrderStatus(orderID uint, status models.OrderStatus) error
	CreateOrderActivity(activity *models.OrderActivity) (*models.OrderActivity, error)
	CheckStatusInHistory(orderID uint, status models.OrderStatus) (bool, *models.OrderActivity, error)
	UpdateOrder(order *models.Order) error

	// Payment operations
	CreatePayment(payment *models.Payment) (*models.Payment, error)
	UpdatePayment(payment *models.Payment) error
	GetPaymentByID(paymentID uint64) (*models.Payment, error)
	GetOrderPayments(orderID uint) ([]models.Payment, error)
	RecalculateOrderPayments(orderID uint) error
}

type VendorRepositoryImpl struct {
	db *gorm.DB
}

// NewVendorRepository creates a new instance of vendor repository
func NewVendorRepository(db *system.DataSource) VendorRepository {
	return &VendorRepositoryImpl{db: db.Db}
}

func (r *VendorRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}

// CreateUser creates a new user
func (r *VendorRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateVendor creates a new vendor profile
func (r *VendorRepositoryImpl) CreateVendor(vendor *models.Vendor) (*models.Vendor, error) {
	if err := r.db.Create(vendor).Error; err != nil {
		return nil, err
	}
	return vendor, nil
}

// GetUserByEmailOrPhone checks if user exists with email or phone
func (r *VendorRepositoryImpl) GetUserByEmailOrPhone(email, phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("(email = ? OR phone = ?) AND is_deleted = ?", email, phone, false).First(&user).Error
	return &user, err
}

// GetVendorByID gets vendor by ID with user info preloaded
func (r *VendorRepositoryImpl) GetVendorByID(vendorID uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Preload("User").Where("id = ? AND is_deleted = ?", vendorID, false).First(&vendor).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (r *VendorRepositoryImpl) GetProductVariantsByIDs(variantIDs []uint) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	err := r.db.Preload("Product").Where("id IN ? AND is_active = ?", variantIDs, true).Find(&variants).Error
	return variants, err
}

func (r *VendorRepositoryImpl) CreateOrder(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func (r *VendorRepositoryImpl) CreateOrderItem(tx *gorm.DB, item *models.OrderItem) error {
	return tx.Create(item).Error
}

func (r *VendorRepositoryImpl) ReduceProductVariantStock(tx *gorm.DB, variantID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}

	result := tx.Model(&models.ProductVariant{}).
		Where("id = ? AND stock >= ?", variantID, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

// GetOrderByID gets order by ID with related data
func (r *VendorRepositoryImpl) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Model(&models.Order{}).Preload("Vendor").Preload("Merchant").Preload("OrderItems.Product").
		Where("id = ? AND is_deleted = ?", orderID, false).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus updates order status
func (r *VendorRepositoryImpl) UpdateOrderStatus(orderID uint, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// CreateOrderActivity creates order activity record
func (r *VendorRepositoryImpl) CreateOrderActivity(activity *models.OrderActivity) (*models.OrderActivity, error) {
	err := r.db.Create(activity).Error
	if err != nil {
		return nil, err
	}
	return activity, nil
}

// CheckStatusInHistory checks if a status was previously recorded in order activity
func (r *VendorRepositoryImpl) CheckStatusInHistory(orderID uint, status models.OrderStatus) (bool, *models.OrderActivity, error) {
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
func (r *VendorRepositoryImpl) CreatePayment(payment *models.Payment) (*models.Payment, error) {
	if err := r.db.Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePayment updates an existing payment record
func (r *VendorRepositoryImpl) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// GetPaymentByID retrieves a payment by ID
func (r *VendorRepositoryImpl) GetPaymentByID(paymentID uint64) (*models.Payment, error) {
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
func (r *VendorRepositoryImpl) GetOrderPayments(orderID uint) ([]models.Payment, error) {
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
func (r *VendorRepositoryImpl) RecalculateOrderPayments(orderID uint) error {
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

// UpdateOrder updates an order record
func (r *VendorRepositoryImpl) UpdateOrder(order *models.Order) error {
	return r.db.Save(order).Error
}
