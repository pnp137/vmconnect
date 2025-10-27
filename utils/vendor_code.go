package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// GenerateVendorCode returns a unique vendor code like SHR1098
func GenerateVendorCode(vendorName string) (string, error) {
	// 1️⃣ Take first 3 uppercase letters of vendor name
	prefix := strings.ToUpper(strings.TrimSpace(vendorName))
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}

	// 2️⃣ Generate cryptographically secure random 4-digit number (1000–9999)
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return "", err
	}
	codeNum := int(n.Int64()) + 1000

	// 3️⃣ Return formatted code
	return fmt.Sprintf("%s%d", prefix, codeNum), nil
}
