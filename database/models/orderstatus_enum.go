package models

import "strconv"

type OrderStatus int

const (
	// Cart State
	ORDER_IN_CART OrderStatus = 5 // Items in cart, not yet placed

	// Initial States
	ORDER_PLACED          OrderStatus = 10 // Merchant placed order
	ORDER_PAYMENT_PENDING OrderStatus = 15 // Payment submitted, awaiting verification

	// Confirmation States
	ORDER_CONFIRMED OrderStatus = 20 // Vendor accepted order

	// Processing States
	ORDER_INVOICED OrderStatus = 30 // Invoice generated

	// Fulfillment States
	ORDER_SHIPPED   OrderStatus = 40 // Vendor sent goods
	ORDER_DELIVERED OrderStatus = 50 // Merchant received goods

	// Completion States
	ORDER_COMPLETED OrderStatus = 60 // Order fully completed (payment settled)

	// Exception States
	ORDER_CANCELLED OrderStatus = 70 // Order cancelled
	ORDER_ON_HOLD   OrderStatus = 75 // Order on hold (optional)
)

func (orderStatus OrderStatus) String() string {
	switch orderStatus {
	case ORDER_IN_CART:
		return "IN_CART"
	case ORDER_PLACED:
		return "PLACED"
	case ORDER_PAYMENT_PENDING:
		return "PAYMENT_PENDING"
	case ORDER_CONFIRMED:
		return "CONFIRMED"
	case ORDER_INVOICED:
		return "INVOICED"
	case ORDER_SHIPPED:
		return "SHIPPED"
	case ORDER_DELIVERED:
		return "DELIVERED"
	case ORDER_COMPLETED:
		return "COMPLETED"
	case ORDER_CANCELLED:
		return "CANCELLED"
	case ORDER_ON_HOLD:
		return "ON_HOLD"
	default:
		return "UNKNOWN STATUS"
	}
}

var (
	orderStatusStringMap = map[string]OrderStatus{
		"IN_CART":         ORDER_IN_CART,
		"PLACED":          ORDER_PLACED,
		"PAYMENT_PENDING": ORDER_PAYMENT_PENDING,
		"CONFIRMED":       ORDER_CONFIRMED,
		"INVOICED":        ORDER_INVOICED,
		"SHIPPED":         ORDER_SHIPPED,
		"DELIVERED":       ORDER_DELIVERED,
		"COMPLETED":       ORDER_COMPLETED,
		"CANCELLED":       ORDER_CANCELLED,
		"ON_HOLD":         ORDER_ON_HOLD,
	}
)

func ParseOrderStatusString(str string) (OrderStatus, bool) {
	c, ok := orderStatusStringMap[str]
	return c, ok
}

var OrderStatusMap = map[OrderStatus]int{
	ORDER_IN_CART:         5,
	ORDER_PLACED:          10,
	ORDER_PAYMENT_PENDING: 15,
	ORDER_CONFIRMED:       20,
	ORDER_INVOICED:        30,
	ORDER_SHIPPED:         40,
	ORDER_DELIVERED:       50,
	ORDER_COMPLETED:       60,
	ORDER_CANCELLED:       70,
	ORDER_ON_HOLD:         75,
}

func CheckIfOrderStatusStringIsValid(str string) bool {
	status, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	for _, v := range OrderStatusMap {
		if v == status {
			return true
		}
	}
	return false
}
