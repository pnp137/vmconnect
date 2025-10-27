# Ultra-Flexible Payment System Design - Unorganized Market

## 🎯 Design Philosophy: "Pay Anytime, Any Mode, Any Amount"

**Key Principles:**
1. **No rigid workflow** - Merchant can pay whenever, however
2. **Multiple payment modes** - Mix UPI, cash, COD in single order
3. **Ledger-based tracking** - Just record all transactions
4. **Auto-calculation** - System calculates totals automatically
5. **Vendor decides fulfillment** - Ship based on trust/policy, not strict rules

---

## 📊 Simplified Data Model

### Order Model (Tracking Only)
```go
type Order struct {
    // ... existing fields ...
    
    // Financial tracking (auto-calculated from Payment records)
    TotalAmount       float64    `json:"total_amount"`
    PaidAmount        float64    `gorm:"default:0" json:"paid_amount"`        // Auto-calculated
    OutstandingAmount float64    `gorm:"default:0" json:"outstanding_amount"` // Auto-calculated
    
    // No strict payment terms - just preferences
    PreferredPaymentMode string `json:"preferred_payment_mode,omitempty"` // Merchant's preference
    
    // Relationships
    Payments []Payment `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}
```

### Payment Model (Ledger Entry)
```go
type Payment struct {
    ID              uint64         `gorm:"primaryKey" json:"id"`
    OrderID         uint           `gorm:"index;not null" json:"order_id"`
    
    // Who & When
    SubmittedBy     uint           `json:"submitted_by"` // User ID who submitted
    ActorRole       *ActorRole     `gorm:"type:int;not null" json:"actor_role"` // merchant or vendor
    
    // Payment Details
    Amount          float64        `gorm:"not null" json:"amount"`
    PaymentMode     *PaymentMode   `gorm:"type:int;not null" json:"payment_mode"` // Cash, Card, UPI, NetBanking, etc.
    PaymentType     string         `json:"payment_type"`  // "token", "partial", "full", "cod_collection"
    
    // Transaction Info
    UTR             string         `json:"utr,omitempty"`
    ReferenceNumber string         `json:"reference_number,omitempty"`
    TransactionID   string         `json:"transaction_id,omitempty"`
    
    // Verification
    Status          *PaymentStatus `gorm:"type:int;not null" json:"status"` // Pending, Processing, Completed, Failed, etc.
    VerifiedBy      uint           `json:"verified_by,omitempty"`
    VerifiedAt      *time.Time     `json:"verified_at,omitempty"`
    
    // Metadata
    Notes           string         `gorm:"type:text" json:"notes,omitempty"`
    Metadata        JSONMap        `gorm:"type:json" json:"metadata,omitempty"`
    
    // Timestamps
    PaidAt          time.Time      `json:"paid_at"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    
    // Relationships
    Order           Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// Available Payment Modes (from payment_mode_enum.go):
// - PaymentModeCash (10)         - Cash payment
// - PaymentModeCard (20)         - Credit/Debit card
// - PaymentModeUPI (30)          - UPI payment
// - PaymentModeNetBanking (40)   - Net banking
// - PaymentModeWallet (50)       - Digital wallet
// - PaymentModeCheque (60)       - Cheque payment
// - PaymentModeBankTransfer (70) - Bank transfer/NEFT/RTGS
// - PaymentModeEMI (80)          - EMI payment

// Available Payment Status (from payment_status_enum.go):
// - PaymentStatusPending (10)    - Payment initiated but not processed
// - PaymentStatusProcessing (20) - Payment being processed
// - PaymentStatusCompleted (30)  - Payment completed successfully
// - PaymentStatusFailed (40)     - Payment failed
// - PaymentStatusCancelled (50)  - Payment cancelled
// - PaymentStatusRefunded (60)   - Payment refunded
// - PaymentStatusPartial (70)    - Partial payment made
// - PaymentStatusExpired (80)    - Payment expired
```

---

## 🔄 Complete User Flow

### Scenario: Merchant's Journey

#### **Step 1: Place Order (No Payment Yet)**
```bash
# Merchant places order without payment
POST /merchant/1/orders/1/place
{
  "notes": "Need these items urgently",
  "preferred_payment_mode": "upi"  # Just a preference
}
```

**Response:**
```json
{
  "order_id": 1,
  "status": "PLACED",
  "total_amount": 5900.00,
  "paid_amount": 0,
  "outstanding_amount": 5900.00,
  "message": "Order placed. You can make payment anytime."
}
```

#### **Step 2: Make Token Payment (₹500 via UPI)**
Merchant makes small token payment when they have cash:
```bash
POST /merchant/1/orders/1/payments
{
  "amount": 500.00,
  "payment_mode": "upi",
  "payment_type": "token",
  "utr": "UTR202510270001",
  "notes": "Token payment - will pay rest later"
}
```

**Backend Automatically:**
```go
// Create payment record
payment1 := Payment{
    OrderID:     1,
    SubmittedBy: merchantUserID,
    ActorRole:   "merchant",
    Amount:      500.00,
    PaymentMode: "upi",
    PaymentType: "token",
    Status:      "pending",
    PaidAt:      time.Now(),
}
db.Create(&payment1)

// Auto-calculate order amounts
order.PaidAmount = CalculateTotalPaid(orderID)  // 500.00
order.OutstandingAmount = order.TotalAmount - order.PaidAmount  // 5400.00
db.Save(&order)
```

**Response:**
```json
{
  "payment_id": 1,
  "order_id": 1,
  "amount": 500.00,
  "payment_mode": "upi",
  "total_paid": 500.00,
  "outstanding": 5400.00,
  "message": "Payment recorded. Balance: ₹5400"
}
```

#### **Step 3: Vendor Ships (Based on Trust/Policy)**
Vendor sees merchant made token payment, decides to ship:
```bash
POST /vendor/1/orders/1/dispatch
{
  "tracking_id": "TRK123",
  "carrier": "BlueDart",
  "notes": "Shipping on credit basis"
}
```

Order Status: `PLACED → SHIPPED` (regardless of payment)

#### **Step 4: Another Payment (₹2000 via Bank Transfer)**
Few days later, merchant gets money and pays more:
```bash
POST /merchant/1/orders/1/payments
{
  "amount": 2000.00,
  "payment_mode": "bank_transfer",
  "payment_type": "partial",
  "reference_number": "REF123456",
  "notes": "Partial payment via bank"
}
```

**Backend:**
```go
payment2 := Payment{
    OrderID:         1,
    Amount:          2000.00,
    PaymentMode:     "bank_transfer",
    PaymentType:     "partial",
    Status:          "pending",
}

// Auto-recalculate
order.PaidAmount = 500 + 2000 = 2500.00
order.OutstandingAmount = 5900 - 2500 = 3400.00
```

**Response:**
```json
{
  "payment_id": 2,
  "total_paid": 2500.00,
  "outstanding": 3400.00,
  "payments_count": 2,
  "message": "Payment recorded. Balance: ₹3400"
}
```

#### **Step 5: COD Payment (₹1000 to Delivery Agent)**
When goods delivered, merchant gives cash to delivery agent:
```bash
POST /merchant/1/orders/1/payments
{
  "amount": 1000.00,
  "payment_mode": "cash",
  "payment_type": "cod_collection",
  "notes": "Cash given to delivery agent",
  "collected_by": "Agent Ram Kumar"
}
```

**Backend:**
```go
payment3 := Payment{
    Amount:      1000.00,
    PaymentMode: "cash",
    PaymentType: "cod_collection",
    Metadata: {
        "collected_by": "Agent Ram Kumar",
        "collection_method": "delivery_agent",
    },
}

order.PaidAmount = 500 + 2000 + 1000 = 3500.00
order.OutstandingAmount = 5900 - 3500 = 2400.00
```

#### **Step 6: Final Settlement (₹2400 via Cheque)**
Later, merchant clears balance with cheque:
```bash
POST /merchant/1/orders/1/payments
{
  "amount": 2400.00,
  "payment_mode": "cheque",
  "payment_type": "final",
  "reference_number": "CHQ789012",
  "notes": "Final settlement via cheque"
}
```

**Backend:**
```go
payment4 := Payment{
    Amount:      2400.00,
    PaymentMode: "cheque",
    PaymentType: "final",
}

order.PaidAmount = 500 + 2000 + 1000 + 2400 = 5900.00
order.OutstandingAmount = 0
// Order status can be marked COMPLETED by merchant
```

---

## 🎯 Key Features

### 1. **No Pre-defined Payment Plan**
- Merchant doesn't commit to payment mode upfront
- Can change mind anytime
- Mix multiple modes in single order

### 2. **Flexible Payment Recording**
```bash
# Merchant can record ANY payment at ANY time:
POST /merchant/1/orders/1/payments
{
  "amount": <any_amount>,
  "payment_mode": <any_mode>,
  "payment_type": <any_type>,
  "notes": <free_text>
}
```

### 3. **Auto-Calculation**
```go
// Backend automatically calculates on every payment
func RecalculateOrderPayments(orderID uint) {
    var payments []Payment
    db.Where("order_id = ? AND status != ?", orderID, "failed").Find(&payments)
    
    totalPaid := 0.0
    for _, p := range payments {
        totalPaid += p.Amount
    }
    
    order.PaidAmount = totalPaid
    order.OutstandingAmount = order.TotalAmount - totalPaid
    db.Save(&order)
}
```

### 4. **Payment Dashboard**
```bash
GET /merchant/1/orders/1/payments

Response:
{
  "order_id": 1,
  "status": "DELIVERED",
  "total_amount": 5900.00,
  "paid_amount": 5900.00,
  "outstanding_amount": 0,
  "payments": [
    {
      "id": 1,
      "amount": 500.00,
      "payment_mode": "upi",
      "payment_type": "token",
      "status": "verified",
      "paid_at": "2025-10-27T10:00:00Z",
      "notes": "Token payment"
    },
    {
      "id": 2,
      "amount": 2000.00,
      "payment_mode": "bank_transfer",
      "payment_type": "partial",
      "status": "verified",
      "paid_at": "2025-10-29T14:30:00Z"
    },
    {
      "id": 3,
      "amount": 1000.00,
      "payment_mode": "cash",
      "payment_type": "cod_collection",
      "status": "pending",
      "paid_at": "2025-10-30T16:00:00Z",
      "notes": "Cash to delivery agent"
    },
    {
      "id": 4,
      "amount": 2400.00,
      "payment_mode": "cheque",
      "payment_type": "final",
      "status": "pending",
      "paid_at": "2025-11-02T11:00:00Z",
      "notes": "Final settlement"
    }
  ],
  "payment_summary": {
    "by_mode": {
      "upi": 500.00,
      "bank_transfer": 2000.00,
      "cash": 1000.00,
      "cheque": 2400.00
    },
    "verified": 2500.00,
    "pending_verification": 3400.00
  }
}
```

---

## 🔧 Simple API Implementation

### Payment Submission API
```go
// Single endpoint for ANY payment
type SubmitPaymentRequest struct {
    Amount          float64                `json:"amount" validate:"required,gt=0"`
    PaymentMode     int                    `json:"payment_mode" validate:"required"` // Enum value: 10=cash, 20=card, 30=upi, etc.
    PaymentType     string                 `json:"payment_type,omitempty"`           // token, partial, final, cod_collection
    UTR             string                 `json:"utr,omitempty"`
    ReferenceNumber string                 `json:"reference_number,omitempty"`
    Notes           string                 `json:"notes,omitempty"`
    Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// In handler: Parse and validate enum
paymentMode := models.PaymentMode(req.PaymentMode)
if !models.CheckIfPaymentModeStringIsValid(strconv.Itoa(req.PaymentMode)) {
    return errors.New("invalid payment mode")
}

func (s *OrderService) SubmitPayment(ctx, merchantID, orderID, req) {
    // 1. Get merchant user
    merchant := s.repo.GetMerchantByID(merchantID)
    
    // 2. Get order
    order := s.repo.GetOrderByID(orderID)
    
    // 3. Validate and parse payment mode enum
    paymentMode := models.PaymentMode(req.PaymentMode)
    if !models.CheckIfPaymentModeStringIsValid(strconv.Itoa(req.PaymentMode)) {
        return nil, errors.New("invalid payment mode")
    }
    
    // 4. Create payment record with typed enums
    actorRole := models.ACTOR_MERCHANT
    paymentStatus := models.PaymentStatusPending
    
    payment := Payment{
        OrderID:         orderID,
        SubmittedBy:     merchant.UserID,
        ActorRole:       &actorRole,
        Amount:          req.Amount,
        PaymentMode:     &paymentMode,
        PaymentType:     req.PaymentType,
        UTR:             req.UTR,
        ReferenceNumber: req.ReferenceNumber,
        Status:          &paymentStatus,
        Notes:           req.Notes,
        Metadata:        req.Metadata,
        PaidAt:          time.Now(),
    }
    s.repo.CreatePayment(&payment)
    
    // 5. Auto-recalculate totals
    s.RecalculateOrderPayments(orderID)
    
    // 6. Get updated amounts
    updatedOrder := s.repo.GetOrderByID(orderID)
    
    // 7. Create order activity
    activity := OrderActivity{
        OrderID:    orderID,
        ActorID:    merchant.UserID,
        ActorRole:  &actorRole,
        OrderState: updatedOrder.Status, // Status doesn't change
        Remarks:    fmt.Sprintf("Payment of ₹%.2f via %s", req.Amount, paymentMode.String()),
        Metadata: JSONMap{
            "action": "submit_payment",
            "payment_id": payment.ID,
            "amount": req.Amount,
            "payment_mode": paymentMode.String(),
            "payment_mode_int": int(paymentMode),
            "total_paid": updatedOrder.PaidAmount,
            "outstanding": updatedOrder.OutstandingAmount,
        },
    }
    s.repo.CreateOrderActivity(&activity)
    
    return PaymentResponse{
        PaymentID:         payment.ID,
        TotalPaid:         updatedOrder.PaidAmount,
        Outstanding:       updatedOrder.OutstandingAmount,
        Message:           fmt.Sprintf("Payment recorded. Balance: ₹%.2f", updatedOrder.OutstandingAmount),
    }
}

// Helper: Recalculate payments
func (s *OrderService) RecalculateOrderPayments(orderID uint) {
    payments := s.repo.GetOrderPayments(orderID)
    
    totalPaid := 0.0
    for _, p := range payments {
        if p.Status != "failed" && p.Status != "cancelled" {
            totalPaid += p.Amount
        }
    }
    
    order := s.repo.GetOrderByID(orderID)
    order.PaidAmount = totalPaid
    order.OutstandingAmount = order.TotalAmount - totalPaid
    s.repo.UpdateOrder(order)
}
```

---

## 🎨 Vendor View & Verification

### Vendor Sees All Payments
```bash
GET /vendor/1/orders/1/payments

Response:
{
  "order_id": 1,
  "merchant_name": "John's Shop",
  "total_amount": 5900.00,
  "payments": [
    {
      "id": 1,
      "amount": 500.00,
      "mode": "upi",
      "utr": "UTR...",
      "status": "pending",
      "submitted_at": "2025-10-27T10:00:00Z"
    },
    {
      "id": 2,
      "amount": 2000.00,
      "mode": "bank_transfer",
      "ref": "REF123456",
      "status": "pending",
      "submitted_at": "2025-10-29T14:30:00Z"
    }
    // ... more
  ],
  "summary": {
    "total_submitted": 3500.00,
    "verified": 500.00,
    "pending_verification": 3000.00
  }
}
```

### Vendor Verifies Payment
```bash
POST /vendor/1/orders/1/payments/1/verify
{
  "verified": true,
  "notes": "Payment confirmed in bank account"
}
```

**Backend:**
```go
func (s *VendorOrderService) VerifyPayment(ctx, vendorID, paymentID, req) {
    payment := s.repo.GetPaymentByID(paymentID)
    vendor := s.repo.GetVendorByID(vendorID)
    
    if req.Verified {
        status := models.PaymentStatusCompleted // Changed from "verified" to enum
        payment.Status = &status
        payment.VerifiedBy = vendor.UserID
        now := time.Now()
        payment.VerifiedAt = &now
    } else {
        status := models.PaymentStatusFailed
        payment.Status = &status
    }
    
    s.repo.UpdatePayment(payment)
    
    // Create activity
    actorRole := models.ACTOR_VENDOR
    activity := OrderActivity{
        OrderID:    payment.OrderID,
        ActorID:    vendor.UserID,
        ActorRole:  &actorRole,
        Remarks:    fmt.Sprintf("Payment #%d %s", paymentID, payment.Status.String()),
        Metadata: JSONMap{
            "action": "verify_payment",
            "payment_id": paymentID,
            "verified": req.Verified,
            "payment_status": payment.Status.String(),
        },
    }
    s.repo.CreateOrderActivity(&activity)
}
```

---

## 📱 Mobile/Web UI Flow

### Merchant Payment Screen
```
┌─────────────────────────────────────┐
│  Order #123                         │
│  Total: ₹5900                       │
│  Paid: ₹2500    Outstanding: ₹3400 │
├─────────────────────────────────────┤
│  💰 Make Payment                    │
│                                     │
│  Amount: [________] ₹               │
│                                     │
│  Payment Mode:                      │
│  ○ UPI                              │
│  ○ Cash                             │
│  ○ Bank Transfer                    │
│  ○ Cheque                           │
│  ○ COD                              │
│                                     │
│  UTR/Ref: [________________]        │
│                                     │
│  Notes: [____________________]      │
│                                     │
│  [Submit Payment] →                 │
├─────────────────────────────────────┤
│  Payment History:                   │
│  • ₹500 (UPI) - Verified ✓         │
│  • ₹2000 (Bank) - Pending ⏳       │
└─────────────────────────────────────┘
```

---

## ✅ Benefits of This Design

1. ✅ **Maximum Flexibility** - Pay anytime, any mode, any amount
2. ✅ **Real Market Behavior** - Matches unorganized market reality
3. ✅ **Simple API** - One endpoint for all payments
4. ✅ **Auto-Calculation** - System handles all math
5. ✅ **Complete Tracking** - Every rupee accounted for
6. ✅ **Mixed Modes** - UPI + Cash + Cheque in same order
7. ✅ **Vendor Control** - Decide when to ship based on trust
8. ✅ **Audit Trail** - Full payment history in OrderActivity
9. ✅ **No Commitments** - Merchant not locked into payment plan
10. ✅ **Easy Reconciliation** - Clear breakdown by mode

---

## 🚀 API Summary

### Merchant APIs
```bash
# Submit any payment anytime
POST   /merchant/{mid}/orders/{oid}/payments

# View all payments
GET    /merchant/{mid}/orders/{oid}/payments

# View payment summary
GET    /merchant/{mid}/orders/{oid}/payment-summary
```

### Vendor APIs
```bash
# View all payments for order
GET    /vendor/{vid}/orders/{oid}/payments

# Verify a payment
POST   /vendor/{vid}/orders/{oid}/payments/{payment_id}/verify

# Mark payment as failed
POST   /vendor/{vid}/orders/{oid}/payments/{payment_id}/reject
```

---

## 💡 Example: Real Unorganized Market Scenario

```
Day 1: Merchant places order (₹5900)
  → No payment yet
  → Status: PLACED

Day 2: Merchant gets ₹500, pays via UPI
  → POST /payments {"amount": 500, "mode": "upi"}
  → Paid: ₹500, Outstanding: ₹5400

Day 3: Vendor ships (trusts merchant)
  → Status: SHIPPED
  → Outstanding: ₹5400 (doesn't matter for shipping)

Day 5: Merchant pays ₹2000 via bank
  → POST /payments {"amount": 2000, "mode": "bank_transfer"}
  → Paid: ₹2500, Outstanding: ₹3400

Day 7: Delivery happens, gives ₹1000 cash to agent
  → POST /payments {"amount": 1000, "mode": "cash"}
  → Paid: ₹3500, Outstanding: ₹2400
  → Status: DELIVERED

Day 15: Merchant settles balance with cheque
  → POST /payments {"amount": 2400, "mode": "cheque"}
  → Paid: ₹5900, Outstanding: ₹0
  → Status: COMPLETED
```

**4 different payments, 4 different modes, across 2 weeks - all tracked perfectly!** ✅

---

This design gives you **complete freedom** while maintaining **full accountability**. Perfect for unorganized markets! 🎉

