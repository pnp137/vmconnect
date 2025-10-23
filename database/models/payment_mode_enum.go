package models

import "strconv"

type PaymentMode int

const (
	PaymentModeCash         PaymentMode = 10 // Cash payment
	PaymentModeCard         PaymentMode = 20 // Credit/Debit card
	PaymentModeUPI          PaymentMode = 30 // UPI payment
	PaymentModeNetBanking   PaymentMode = 40 // Net banking
	PaymentModeWallet       PaymentMode = 50 // Digital wallet (Paytm, PhonePe, etc.)
	PaymentModeCheque       PaymentMode = 60 // Cheque payment
	PaymentModeBankTransfer PaymentMode = 70 // Bank transfer/NEFT/RTGS
	PaymentModeEMI          PaymentMode = 80 // EMI payment
)

func (paymentMode PaymentMode) String() string {
	switch paymentMode {
	case PaymentModeCash:
		return "cash"
	case PaymentModeCard:
		return "card"
	case PaymentModeUPI:
		return "upi"
	case PaymentModeNetBanking:
		return "net_banking"
	case PaymentModeWallet:
		return "wallet"
	case PaymentModeCheque:
		return "cheque"
	case PaymentModeBankTransfer:
		return "bank_transfer"
	case PaymentModeEMI:
		return "emi"
	default:
		return "unknown mode"
	}
}

var (
	paymentModeStringMap = map[string]PaymentMode{
		"cash":          PaymentModeCash,
		"card":          PaymentModeCard,
		"upi":           PaymentModeUPI,
		"net_banking":   PaymentModeNetBanking,
		"wallet":        PaymentModeWallet,
		"cheque":        PaymentModeCheque,
		"bank_transfer": PaymentModeBankTransfer,
		"emi":           PaymentModeEMI,
	}
)

func ParsePaymentModeString(str string) (PaymentMode, bool) {
	c, ok := paymentModeStringMap[str]
	return c, ok
}

var PaymentModeMap = map[PaymentMode]int{
	PaymentModeCash:         10,
	PaymentModeCard:         20,
	PaymentModeUPI:          30,
	PaymentModeNetBanking:   40,
	PaymentModeWallet:       50,
	PaymentModeCheque:       60,
	PaymentModeBankTransfer: 70,
	PaymentModeEMI:          80,
}

func CheckIfPaymentModeStringIsValid(str string) bool {
	mode, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	for _, v := range PaymentModeMap {
		if v == mode {
			return true
		}
	}
	return false
}
