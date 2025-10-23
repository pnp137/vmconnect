package models

import "strconv"

type InvoiceStatus int

const (
	InvoiceStatusPending   InvoiceStatus = 10 // Vendor created draft invoice
	InvoiceStatusSent      InvoiceStatus = 20 // Vendor sent to merchant for review
	InvoiceStatusConfirmed InvoiceStatus = 30 // Vendor confirmed invoice details
	InvoiceStatusRejected  InvoiceStatus = 40 // Vendor rejected (maybe pricing/quantity mismatch)
	InvoiceStatusAccepted  InvoiceStatus = 50 // Merchant accepted invoice
	InvoiceStatusPaid      InvoiceStatus = 60 // Fully paid
	InvoiceStatusPartial   InvoiceStatus = 70 // Partially paid
	InvoiceStatusCancelled InvoiceStatus = 80 // Cancelled (by vendor or system)
)

func (invoiceStatus InvoiceStatus) String() string {
	switch invoiceStatus {
	case InvoiceStatusPending:
		return "pending"
	case InvoiceStatusSent:
		return "sent"
	case InvoiceStatusConfirmed:
		return "confirmed"
	case InvoiceStatusRejected:
		return "rejected"
	case InvoiceStatusAccepted:
		return "accepted"
	case InvoiceStatusPaid:
		return "paid"
	case InvoiceStatusPartial:
		return "partial"
	case InvoiceStatusCancelled:
		return "cancelled"
	default:
		return "unknown status"
	}
}

var (
	invoiceStatusStringMap = map[string]InvoiceStatus{
		"pending":   InvoiceStatusPending,
		"sent":      InvoiceStatusSent,
		"confirmed": InvoiceStatusConfirmed,
		"rejected":  InvoiceStatusRejected,
		"accepted":  InvoiceStatusAccepted,
		"paid":      InvoiceStatusPaid,
		"partial":   InvoiceStatusPartial,
		"cancelled": InvoiceStatusCancelled,
	}
)

func ParseInvoiceStatusString(str string) (InvoiceStatus, bool) {
	c, ok := invoiceStatusStringMap[str]
	return c, ok
}

var InvoiceStatusMap = map[InvoiceStatus]int{
	InvoiceStatusPending:   10,
	InvoiceStatusSent:      20,
	InvoiceStatusConfirmed: 30,
	InvoiceStatusRejected:  40,
	InvoiceStatusAccepted:  50,
	InvoiceStatusPaid:      60,
	InvoiceStatusPartial:   70,
	InvoiceStatusCancelled: 80,
}

func CheckIfInvoiceStatusStringIsValid(str string) bool {
	status, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	for _, v := range InvoiceStatusMap {
		if v == status {
			return true
		}
	}
	return false
}
