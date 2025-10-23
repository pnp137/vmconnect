package models

import "strconv"

type PaymentStatus int

const (
	PaymentStatusPending    PaymentStatus = 10 // Payment initiated but not processed
	PaymentStatusProcessing PaymentStatus = 20 // Payment being processed
	PaymentStatusCompleted  PaymentStatus = 30 // Payment completed successfully
	PaymentStatusFailed     PaymentStatus = 40 // Payment failed
	PaymentStatusCancelled  PaymentStatus = 50 // Payment cancelled by user/system
	PaymentStatusRefunded   PaymentStatus = 60 // Payment refunded
	PaymentStatusPartial    PaymentStatus = 70 // Partial payment made
	PaymentStatusExpired    PaymentStatus = 80 // Payment expired
)

func (paymentStatus PaymentStatus) String() string {
	switch paymentStatus {
	case PaymentStatusPending:
		return "pending"
	case PaymentStatusProcessing:
		return "processing"
	case PaymentStatusCompleted:
		return "completed"
	case PaymentStatusFailed:
		return "failed"
	case PaymentStatusCancelled:
		return "cancelled"
	case PaymentStatusRefunded:
		return "refunded"
	case PaymentStatusPartial:
		return "partial"
	case PaymentStatusExpired:
		return "expired"
	default:
		return "unknown status"
	}
}

var (
	paymentStatusStringMap = map[string]PaymentStatus{
		"pending":    PaymentStatusPending,
		"processing": PaymentStatusProcessing,
		"completed":  PaymentStatusCompleted,
		"failed":     PaymentStatusFailed,
		"cancelled":  PaymentStatusCancelled,
		"refunded":   PaymentStatusRefunded,
		"partial":    PaymentStatusPartial,
		"expired":    PaymentStatusExpired,
	}
)

func ParsePaymentStatusString(str string) (PaymentStatus, bool) {
	c, ok := paymentStatusStringMap[str]
	return c, ok
}

var PaymentStatusMap = map[PaymentStatus]int{
	PaymentStatusPending:    10,
	PaymentStatusProcessing: 20,
	PaymentStatusCompleted:  30,
	PaymentStatusFailed:     40,
	PaymentStatusCancelled:  50,
	PaymentStatusRefunded:   60,
	PaymentStatusPartial:    70,
	PaymentStatusExpired:    80,
}

func CheckIfPaymentStatusStringIsValid(str string) bool {
	status, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	for _, v := range PaymentStatusMap {
		if v == status {
			return true
		}
	}
	return false
}
