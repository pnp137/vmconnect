# Payment Handling Design - Credit & Flexible Payment Scenarios

## 🎯 Design Philosophy

**Key Principle:** Decouple order fulfillment from payment tracking.

- **Order Status** = Physical goods flow (confirmed → shipped → delivered)
- **Payment Status** = Financial flow (tracked separately)
- **Payments** = Multiple records linked to order/invoice

---

## 📊 Current vs Enhanced Design

### Current Design (Tightly Coupled)
```
CONFIRMED → INVOICED → PAID → SHIPPED → DELIVERED
```
**Problem:** Forces payment before shipment

### Enhanced Design (Decoupled)
```
Order Flow:    CONFIRMED → READY_TO_SHIP → SHIPPED → DELIVERED → COMPLETED
                    ↓            ↓              ↓          ↓
Payment Flow:  [Payment 1] [Payment 2]   [Payment 3] [Payment 4]
               (Anytime before, during, or after delivery)
```

---

## 🆕 Proposed Changes

### 1. **New Order Statuses**

Add to `orderstatus_enum.go`:

```go
const (
    // ... existing statuses ...
    ORDER_CONFIRMED      OrderStatus = 30   // Existing
    ORDER_INVOICED       OrderStatus = 40   // Keep for backward compatibility
    ORDER_PARTIALLY_PAID OrderStatus = 50   // Keep for tracking
    ORDER_PAID           OrderStatus = 60   // Keep for fully paid
    ORDER_ON_CREDIT      OrderStatus = 65   // NEW: Goods on credit
    ORDER_READY_TO_SHIP  OrderStatus = 70   // Existing
    // ... rest
)
```

### 2. **Enhanced Payment Tracking**

#### Update Order Model
Add payment tracking fields to `Order`:

```go
type Order struct {
    // ... existing fields ...
    
    // Payment tracking
    TotalAmount       float64           `json:"total_amount"`
    PaidAmount        float64           `gorm:"default:0" json:"paid_amount"`
    OutstandingAmount float64           `gorm:"default:0" json:"outstanding_amount"`
    PaymentTerms      string            `json:"payment_terms,omitempty"` // "upfront", "credit", "cod", "custom"
    CreditDays        int               `gorm:"default:0" json:"credit_days"` // Credit period
    PaymentDueDate    *time.Time        `json:"payment_due_date,omitempty"`
    
    // Relationships
    Payments          []Payment         `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}
```

#### Enhance Payment Model

```go
type Payment struct {
    ID              uint64         `gorm:"primaryKey"`
    OrderID         uint           `gorm:"index;not null" json:"order_id"`
    InvoiceID       uint64         `gorm:"index" json:"invoice_id,omitempty"`
    Amount          float64        `gorm:"not null" json:"amount"`
    PaymentType     string         `json:"payment_type"` // "advance", "partial", "final", "credit_settlement"
    PaymentMode     string         `json:"payment_mode"` // "upi", "cash", "cheque", etc.
    TransactionID   string         `json:"transaction_id,omitempty"`
    UTR             string         `json:"utr,omitempty"`
    Status          *PaymentStatus `gorm:"type:int"`
    VerifiedBy      uint           `json:"verified_by,omitempty"` // User ID of verifier
    VerifiedAt      *time.Time     `json:"verified_at,omitempty"`
    PaidAt          *time.Time     `json:"paid_at"`
    Notes           string         `gorm:"type:text" json:"notes,omitempty"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
}
```

---

## 💼 Business Scenarios & Implementation

### Scenario 1: Full Upfront Payment (Existing)
```
Flow: PLACED → CONFIRMED → INVOICED → PAID → SHIPPED → DELIVERED

Merchant pays 100% before shipment
```

**API Call:**
```bash
POST /merchant/1/orders/1/payment
{
  "amount": 5900.00,
  "payment_mode": "upi",
  "payment_type": "full",
  "utr": "UTR123"
}
```

**Backend Logic:**
```go
if payment.Amount >= order.TotalAmount {
    order.PaidAmount = order.TotalAmount
    order.OutstandingAmount = 0
    order.Status = ORDER_PAID
}
```

---

### Scenario 2: Token/Advance Payment + Balance Later
```
Flow: PLACED → CONFIRMED → INVOICED → PARTIALLY_PAID → SHIPPED → 
      DELIVERED → [Balance Payment] → COMPLETED
```

**Step 1: Token Payment (20%)**
```bash
POST /merchant/1/orders/1/payment
{
  "amount": 1180.00,
  "payment_mode": "upi",
  "payment_type": "advance",
  "utr": "UTR001"
}
```

**Backend:**
```go
order.PaidAmount += 1180.00  // 1180
order.OutstandingAmount = order.TotalAmount - order.PaidAmount  // 4720
order.Status = ORDER_PARTIALLY_PAID

// Create Payment record
payment1 := Payment{
    OrderID: 1,
    Amount: 1180.00,
    PaymentType: "advance",
    Status: PaymentStatusPending,
}
```

**Step 2: Balance Payment (80%)**
```bash
POST /merchant/1/orders/1/payment
{
  "amount": 4720.00,
  "payment_mode": "bank_transfer",
  "payment_type": "final",
  "utr": "UTR002"
}
```

**Backend:**
```go
order.PaidAmount += 4720.00  // 5900
order.OutstandingAmount = 0
order.Status = ORDER_PAID

payment2 := Payment{
    OrderID: 1,
    Amount: 4720.00,
    PaymentType: "final",
}
```

---

### Scenario 3: Full Credit (Pay After Delivery) ⭐ NEW
```
Flow: PLACED → CONFIRMED → ON_CREDIT → READY_TO_SHIP → SHIPPED → 
      DELIVERED → [Payment] → COMPLETED
```

**Step 1: Vendor Approves Credit**
```bash
POST /vendor/1/orders/1/approve-credit
{
  "credit_days": 30,
  "payment_terms": "Net 30",
  "notes": "Approved for 30-day credit"
}
```

**Backend:**
```go
func (s *OrderService) ApproveCreditTerms(ctx, vendorID, orderID, req) {
    order.Status = ORDER_ON_CREDIT
    order.PaymentTerms = "credit"
    order.CreditDays = 30
    order.PaymentDueDate = time.Now().AddDate(0, 0, 30)
    order.PaidAmount = 0
    order.OutstandingAmount = order.TotalAmount
    
    // Create activity
    activity := OrderActivity{
        OrderState: ORDER_ON_CREDIT,
        ActorRole: ACTOR_VENDOR,
        Metadata: {
            "action": "approve_credit",
            "credit_days": 30,
            "payment_due": order.PaymentDueDate,
        },
    }
}
```

**Step 2: Ship Without Payment**
```bash
POST /vendor/1/orders/1/dispatch
{
  "tracking_id": "TRK123",
  "carrier": "BlueDart"
}
```

Status: `ON_CREDIT → SHIPPED` (no payment needed)

**Step 3: Payment After Delivery**
```bash
POST /merchant/1/orders/1/payment
{
  "amount": 5900.00,
  "payment_mode": "bank_transfer",
  "payment_type": "credit_settlement",
  "utr": "UTR123"
}
```

---

### Scenario 4: COD (Cash on Delivery)
```
Flow: PLACED → CONFIRMED → INVOICED → READY_TO_SHIP → SHIPPED → 
      DELIVERED → [COD Payment] → COMPLETED
```

**Step 1: Mark as COD**
```bash
POST /vendor/1/orders/1/set-payment-terms
{
  "payment_terms": "cod",
  "notes": "Cash on delivery"
}
```

**Backend:**
```go
order.PaymentTerms = "cod"
order.Status = ORDER_CONFIRMED  // Can ship without payment
```

**Step 2: Record COD Payment**
```bash
POST /vendor/1/orders/1/payment/record-cod
{
  "amount": 5900.00,
  "payment_mode": "cash",
  "collected_by": "delivery_agent",
  "notes": "Cash collected on delivery"
}
```

---

### Scenario 5: Multiple Installments
```
Flow: PLACED → CONFIRMED → PARTIALLY_PAID → [Install 1] → 
      SHIPPED → [Install 2] → DELIVERED → [Install 3] → COMPLETED
```

**Installment 1 (30%)**
```bash
POST /merchant/1/orders/1/payment
{
  "amount": 1770.00,
  "payment_type": "installment_1",
  "installment_number": 1,
  "total_installments": 3
}
```

**Installment 2 (30%)**
```bash
POST /merchant/1/orders/1/payment
{
  "amount": 1770.00,
  "payment_type": "installment_2",
  "installment_number": 2
}
```

**Installment 3 (40%)**
```bash
POST /merchant/1/orders/1/payment
{
  "amount": 2360.00,
  "payment_type": "installment_3_final",
  "installment_number": 3
}
```

---

## 🔧 Implementation Changes

### 1. New DTOs

#### Merchant DTOs
```go
// ApproveCredit request (can be merchant-initiated)
type RequestCreditRequest struct {
    CreditDays int    `json:"credit_days" validate:"required,min=1,max=90"`
    Notes      string `json:"notes,omitempty"`
}

// Multi-payment support
type SubmitPaymentRequest struct {
    Amount            float64 `json:"amount" validate:"required,gt=0"`
    PaymentMode       string  `json:"payment_mode" validate:"required"`
    PaymentType       string  `json:"payment_type"` // "advance", "partial", "final", "credit_settlement"
    InstallmentNumber int     `json:"installment_number,omitempty"`
    TotalInstallments int     `json:"total_installments,omitempty"`
    UTR               string  `json:"utr,omitempty"`
    Notes             string  `json:"notes,omitempty"`
}

// Payment status response
type PaymentStatusResponse struct {
    OrderID           uint      `json:"order_id"`
    TotalAmount       float64   `json:"total_amount"`
    PaidAmount        float64   `json:"paid_amount"`
    OutstandingAmount float64   `json:"outstanding_amount"`
    PaymentTerms      string    `json:"payment_terms"`
    PaymentDueDate    *string   `json:"payment_due_date,omitempty"`
    Payments          []Payment `json:"payments"`
}
```

#### Vendor DTOs
```go
type ApproveCreditRequest struct {
    CreditDays   int    `json:"credit_days" validate:"required,min=1,max=90"`
    PaymentTerms string `json:"payment_terms"` // "Net 30", "Net 60", etc.
    MaxCredit    float64 `json:"max_credit,omitempty"`
    Notes        string `json:"notes,omitempty"`
}

type RecordCODPaymentRequest struct {
    Amount       float64 `json:"amount" validate:"required,gt=0"`
    PaymentMode  string  `json:"payment_mode" validate:"required"`
    CollectedBy  string  `json:"collected_by"`
    Notes        string  `json:"notes,omitempty"`
}
```

---

### 2. Enhanced Service Methods

```go
// Merchant Service
func (s *OrderServiceImpl) SubmitPayment(ctx, merchantID, orderID, req) {
    order := s.repo.GetOrderByID(orderID)
    
    // Calculate new amounts
    newPaidAmount := order.PaidAmount + req.Amount
    newOutstanding := order.TotalAmount - newPaidAmount
    
    // Create payment record
    payment := Payment{
        OrderID:       orderID,
        Amount:        req.Amount,
        PaymentType:   req.PaymentType,
        PaymentMode:   req.PaymentMode,
        UTR:           req.UTR,
        Status:        PaymentStatusPending,
        PaidAt:        time.Now(),
    }
    s.repo.CreatePayment(&payment)
    
    // Update order
    order.PaidAmount = newPaidAmount
    order.OutstandingAmount = newOutstanding
    
    // Determine order status
    if newOutstanding == 0 {
        order.Status = ORDER_PAID
    } else if newPaidAmount > 0 {
        order.Status = ORDER_PARTIALLY_PAID
    }
    
    s.repo.UpdateOrder(order)
    
    // Create activity
    activity := OrderActivity{
        OrderState: order.Status,
        Metadata: {
            "action": "submit_payment",
            "payment_id": payment.ID,
            "amount": req.Amount,
            "payment_type": req.PaymentType,
            "paid_amount": newPaidAmount,
            "outstanding": newOutstanding,
        },
    }
    s.repo.CreateOrderActivity(activity)
}

// Vendor Service
func (s *OrderServiceImpl) ApproveCredit(ctx, vendorID, orderID, req) {
    order := s.repo.GetOrderByID(orderID)
    
    // Set credit terms
    order.PaymentTerms = "credit"
    order.CreditDays = req.CreditDays
    dueDate := time.Now().AddDate(0, 0, req.CreditDays)
    order.PaymentDueDate = &dueDate
    order.OutstandingAmount = order.TotalAmount
    order.Status = ORDER_ON_CREDIT
    
    s.repo.UpdateOrder(order)
    
    // Create activity
    activity := OrderActivity{
        OrderState: ORDER_ON_CREDIT,
        ActorRole: ACTOR_VENDOR,
        Metadata: {
            "action": "approve_credit",
            "credit_days": req.CreditDays,
            "payment_due": dueDate,
            "payment_terms": req.PaymentTerms,
        },
    }
    s.repo.CreateOrderActivity(activity)
}

func (s *OrderServiceImpl) RecordCODPayment(ctx, vendorID, orderID, req) {
    order := s.repo.GetOrderByID(orderID)
    
    // Create payment record
    payment := Payment{
        OrderID:     orderID,
        Amount:      req.Amount,
        PaymentType: "cod",
        PaymentMode: req.PaymentMode,
        Status:      PaymentStatusCompleted,
        VerifiedBy:  vendorID,
        VerifiedAt:  time.Now(),
        PaidAt:      time.Now(),
    }
    s.repo.CreatePayment(&payment)
    
    // Update order
    order.PaidAmount = req.Amount
    order.OutstandingAmount = 0
    order.Status = ORDER_PAID
    
    s.repo.UpdateOrder(order)
}
```

---

### 3. New API Endpoints

#### Merchant Endpoints
```
POST   /merchant/{mid}/orders/{oid}/payment          # Submit any payment
GET    /merchant/{mid}/orders/{oid}/payment-status   # Get payment details
POST   /merchant/{mid}/orders/{oid}/request-credit   # Request credit terms
```

#### Vendor Endpoints
```
POST   /vendor/{vid}/orders/{oid}/approve-credit     # Approve credit
POST   /vendor/{vid}/orders/{oid}/set-payment-terms  # Set COD/other terms
POST   /vendor/{vid}/orders/{oid}/payment/record-cod # Record COD payment
GET    /vendor/{vid}/orders/{oid}/payment-status     # Get payment details
POST   /vendor/{vid}/orders/{oid}/payment/verify     # Verify payment
```

---

### 4. Updated State Transition Rules

```go
// Allow shipment without full payment for credit/COD
var flexibleTransitions = []StateTransitionRule{
    {
        FromStatus:   ORDER_CONFIRMED,
        ToStatuses:   []OrderStatus{
            ORDER_INVOICED,           // Normal flow
            ORDER_PARTIALLY_PAID,     // Partial payment
            ORDER_ON_CREDIT,          // Credit approved
            ORDER_READY_TO_SHIP,      // Skip payment (COD)
        },
        AllowedRoles: []ActorRole{ACTOR_VENDOR, ACTOR_MERCHANT},
    },
    {
        FromStatus:   ORDER_ON_CREDIT,
        ToStatuses:   []OrderStatus{
            ORDER_READY_TO_SHIP,      // Can ship on credit
            ORDER_SHIPPED,            // Can ship on credit
            ORDER_PARTIALLY_PAID,     // Payment received
            ORDER_PAID,               // Full payment received
        },
        AllowedRoles: []ActorRole{ACTOR_VENDOR, ACTOR_MERCHANT},
    },
}
```

---

## 📊 Payment Tracking Dashboard

### Merchant View
```json
{
  "order_id": 1,
  "status": "SHIPPED",
  "total_amount": 5900.00,
  "paid_amount": 1770.00,
  "outstanding_amount": 4130.00,
  "payment_terms": "installment",
  "payment_due_date": "2025-11-30",
  "payments": [
    {
      "payment_id": 1,
      "amount": 1770.00,
      "payment_type": "installment_1",
      "status": "completed",
      "paid_at": "2025-10-26T10:00:00Z"
    }
  ]
}
```

### Vendor View (with verification status)
```json
{
  "order_id": 1,
  "merchant_name": "John's Shop",
  "total_amount": 5900.00,
  "paid_amount": 1770.00,
  "outstanding_amount": 4130.00,
  "payment_terms": "installment",
  "payments": [
    {
      "payment_id": 1,
      "amount": 1770.00,
      "utr": "UTR123",
      "status": "pending_verification",
      "verified_by": null,
      "paid_at": "2025-10-26T10:00:00Z"
    }
  ]
}
```

---

## ✅ Benefits of This Design

1. **Flexible** - Supports all payment scenarios
2. **Decoupled** - Order flow independent of payment
3. **Trackable** - Multiple payment records with history
4. **Real-world** - Matches actual business practices
5. **Credit Management** - Built-in credit term tracking
6. **COD Support** - Easy cash-on-delivery handling
7. **Installments** - Multiple partial payments
8. **Auditable** - Full payment history in OrderActivity

---

## 🔄 Example Flows

### Credit Flow
```
1. PLACED → CONFIRMED
2. Vendor approves credit → ON_CREDIT
3. Vendor ships → SHIPPED
4. Merchant receives → DELIVERED  
5. Merchant pays (before due date) → Payment recorded
6. Merchant completes → COMPLETED

Payment: Tracked separately, outstanding amount monitored
```

### COD Flow
```
1. PLACED → CONFIRMED
2. Vendor sets COD terms → CONFIRMED (payment_terms="cod")
3. Vendor ships → SHIPPED
4. Merchant receives → DELIVERED
5. Vendor records COD payment → Payment recorded
6. Merchant completes → COMPLETED
```

---

## 🚀 Next Steps

1. Update Order model with payment tracking fields
2. Update Payment model with enhanced fields
3. Add ORDER_ON_CREDIT to enum
4. Implement new service methods
5. Create new API endpoints
6. Update state transition rules
7. Add payment dashboard endpoints

Would you like me to implement these changes?

