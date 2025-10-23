package models

import "strconv"

type OrderStatus int

const (
	ORDER_IN_CART        OrderStatus = 10
	ORDER_PLACED         OrderStatus = 20
	ORDER_CONFIRMED      OrderStatus = 30
	ORDER_INVOICED       OrderStatus = 40
	ORDER_PARTIALLY_PAID OrderStatus = 50
	ORDER_PAID           OrderStatus = 60
	ORDER_READY_TO_SHIP  OrderStatus = 70
	ORDER_SHIPPED        OrderStatus = 80
	ORDER_DELIVERED      OrderStatus = 90
	ORDER_COMPLETED      OrderStatus = 100
	ORDER_CANCELLED      OrderStatus = 110
	ORDER_RETURNED       OrderStatus = 120
)

func (orderStatus OrderStatus) String() string {
	switch orderStatus {
	case ORDER_IN_CART:
		return "IN_CART"
	case ORDER_PLACED:
		return "PLACED"
	case ORDER_CONFIRMED:
		return "CONFIRMED"
	case ORDER_INVOICED:
		return "INVOICED"
	case ORDER_PARTIALLY_PAID:
		return "PARTIALLY_PAID"
	case ORDER_PAID:
		return "PAID"
	case ORDER_READY_TO_SHIP:
		return "READY_TO_SHIP"
	case ORDER_SHIPPED:
		return "SHIPPED"
	case ORDER_DELIVERED:
		return "DELIVERED"
	case ORDER_COMPLETED:
		return "COMPLETED"
	case ORDER_CANCELLED:
		return "CANCELLED"
	case ORDER_RETURNED:
		return "RETURNED"
	default:
		return "UNKNOWN STATUS"
	}
}

var (
	orderStatusStringMap = map[string]OrderStatus{
		"IN_CART":        ORDER_IN_CART,
		"PLACED":         ORDER_PLACED,
		"CONFIRMED":      ORDER_CONFIRMED,
		"INVOICED":       ORDER_INVOICED,
		"PARTIALLY_PAID": ORDER_PARTIALLY_PAID,
		"PAID":           ORDER_PAID,
		"READY_TO_SHIP":  ORDER_READY_TO_SHIP,
		"SHIPPED":        ORDER_SHIPPED,
		"DELIVERED":      ORDER_DELIVERED,
		"COMPLETED":      ORDER_COMPLETED,
		"CANCELLED":      ORDER_CANCELLED,
		"RETURNED":       ORDER_RETURNED,
	}
)

func ParseOrderStatusString(str string) (OrderStatus, bool) {
	c, ok := orderStatusStringMap[str]
	return c, ok
}

var OrderStatusMap = map[OrderStatus]int{
	ORDER_IN_CART:        10,
	ORDER_PLACED:         20,
	ORDER_CONFIRMED:      30,
	ORDER_INVOICED:       40,
	ORDER_PARTIALLY_PAID: 50,
	ORDER_PAID:           60,
	ORDER_READY_TO_SHIP:  70,
	ORDER_SHIPPED:        80,
	ORDER_DELIVERED:      90,
	ORDER_COMPLETED:      100,
	ORDER_CANCELLED:      110,
	ORDER_RETURNED:       120,
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
