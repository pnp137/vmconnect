# 🧪 Flexible Payment System - API Testing Guide

## 🎯 Complete Order + Payment Flow Testing

This guide provides step-by-step API testing for the flexible payment system.

---

## 📋 Prerequisites

1. ✅ Database migration completed (`Payment` table created)
2. ✅ Server running on `http://localhost:8080`
3. ✅ Registered merchant and vendor accounts
4. ✅ Products added to vendor catalog
5. ✅ Merchant has added vendor

---

## 🔐 Test Setup

### Get IDs (Replace with your actual IDs):
```bash
MERCHANT_ID=1
VENDOR_ID=1
ORDER_ID=1
PRODUCT_ID=1
```

---

## 📦 **Scenario 1: Token Payment → Vendor Confirms → Ship → Settle**

### Step 1: Merchant Places Order
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/place" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Need these items urgently"
  }'
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "status": "PLACED",
    "total_amount": 5900.00,
    "paid_amount": 0,
    "outstanding_amount": 5900.00
  }
}
```

---

### Step 2: Merchant Submits Token Payment (₹500 via UPI)
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 500.00,
    "payment_mode": 30,
    "payment_type": "token",
    "utr": "UTR202510270001",
    "notes": "Token payment - will pay rest later"
  }'
```

**Payment Mode Values:**
- `10` = Cash
- `20` = Card  
- `30` = UPI
- `40` = Net Banking
- `50` = Wallet
- `60` = Cheque
- `70` = Bank Transfer
- `80` = EMI

**Expected Response:**
```json
{
  "content": {
    "payment_id": 1,
    "order_id": 1,
    "amount": 500.00,
    "payment_mode": "upi",
    "payment_status": "pending",
    "total_paid": 500.00,
    "outstanding": 5400.00,
    "order_status": "PAYMENT_PENDING",
    "order_paid": false,
    "submitted_at": "2025-10-27T10:15:00Z",
    "message": "Payment recorded successfully. Balance: ₹5400.00"
  }
}
```

**✅ Order Status Changed:** `PLACED` → `PAYMENT_PENDING`

---

### Step 3: Vendor Views Payments
```bash
curl -X GET "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/payments"
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "merchant_name": "John's Shop",
    "order_status": "PAYMENT_PENDING",
    "total_amount": 5900.00,
    "paid_amount": 500.00,
    "outstanding_amount": 5400.00,
    "order_paid": false,
    "payments": [
      {
        "id": 1,
        "amount": 500.00,
        "payment_mode": "upi",
        "payment_type": "token",
        "payment_status": "pending",
        "verified": false,
        "utr": "UTR202510270001",
        "notes": "Token payment",
        "submitted_at": "2025-10-27T10:15:00Z"
      }
    ],
    "summary": {
      "total_submitted": 500.00,
      "verified": 0,
      "pending_verification": 500.00,
      "payments_count": 1
    }
  }
}
```

---

### Step 4: Vendor Marks Order as Paid
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/mark-paid" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Token payment received in bank account"
  }'
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "order_paid": true,
    "payments_verified": 1,
    "total_amount": 500.00,
    "marked_at": "2025-10-27T11:30:00Z",
    "order_status": "CONFIRMED",
    "message": "Order marked as paid successfully. 1 payments verified."
  }
}
```

**✅ What Happened:**
- All `pending` payments → `verified: true`, `status: completed`
- `order_paid` → `true`
- Order status changed: `PAYMENT_PENDING` → `CONFIRMED`

---

### Step 5: Vendor Ships Order
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/dispatch" \
  -H "Content-Type: application/json" \
  -d '{
    "tracking_id": "TRK123456789",
    "carrier": "BlueDart",
    "notes": "Shipped on credit basis"
  }'
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "status": "SHIPPED",
    "tracking_id": "TRK123456789",
    "carrier": "BlueDart",
    "shipped_at": "2025-10-28T11:00:00Z"
  }
}
```

---

### Step 6: Merchant Receives Order
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/receive" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Goods received in good condition"
  }'
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "status": "DELIVERED",
    "delivered_at": "2025-10-30T16:00:00Z"
  }
}
```

---

### Step 7: Merchant Submits Remaining Payment (₹5400 via Cash)
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5400.00,
    "payment_mode": 10,
    "payment_type": "settlement",
    "notes": "Final payment - cash to delivery agent"
  }'
```

**Expected Response:**
```json
{
  "content": {
    "payment_id": 2,
    "order_id": 1,
    "amount": 5400.00,
    "payment_mode": "cash",
    "payment_status": "pending",
    "total_paid": 5900.00,
    "outstanding": 0,
    "order_status": "DELIVERED",
    "order_paid": false,
    "message": "Payment recorded successfully. Balance: ₹0.00"
  }
}
```

**✅ Note:** Order status stays `DELIVERED` (payment independent of order flow)

---

### Step 8: Vendor Marks Final Payment as Received
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/mark-paid" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Final settlement received"
  }'
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "order_paid": true,
    "payments_verified": 1,
    "total_amount": 5900.00,
    "marked_at": "2025-11-02T11:00:00Z",
    "order_status": "DELIVERED",
    "message": "Order marked as paid successfully. 1 payments verified."
  }
}
```

---

### Step 9: View Complete Payment History
```bash
curl -X GET "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments"
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "order_status": "DELIVERED",
    "total_amount": 5900.00,
    "paid_amount": 5900.00,
    "outstanding_amount": 0,
    "order_paid": true,
    "order_paid_at": "2025-11-02T11:00:00Z",
    "payments": [
      {
        "id": 1,
        "amount": 500.00,
        "payment_mode": "upi",
        "payment_type": "token",
        "payment_status": "completed",
        "verified": true,
        "utr": "UTR202510270001",
        "paid_at": "2025-10-27T10:15:00Z",
        "verified_at": "2025-10-27T11:30:00Z"
      },
      {
        "id": 2,
        "amount": 5400.00,
        "payment_mode": "cash",
        "payment_type": "settlement",
        "payment_status": "completed",
        "verified": true,
        "paid_at": "2025-10-30T16:30:00Z",
        "verified_at": "2025-11-02T11:00:00Z"
      }
    ],
    "payment_summary": {
      "by_mode": {
        "upi": 500.00,
        "cash": 5400.00
      },
      "verified": 5900.00,
      "pending_verification": 0
    }
  }
}
```

---

## 📦 **Scenario 2: Credit → Ship → Pay Later**

### Step 1: Merchant Places Order (No Payment)
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/place" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Need on credit"
  }'
```

### Step 2: Vendor Confirms Order (No Payment Verification)
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/confirm" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Confirmed on credit basis"
  }'
```

**Expected Response:**
```json
{
  "content": {
    "order_id": 1,
    "status": "CONFIRMED",
    "message": "Order confirmed successfully"
  }
}
```

**✅ Note:** Order moved from `PLACED` → `CONFIRMED` without any payment!

### Step 3: Vendor Ships on Credit
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/dispatch" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Shipping on credit - trusted merchant"
  }'
```

### Step 4: Delivered
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/receive" \
  -H "Content-Type: application/json" \
  -d '{}'
```

### Step 5: Merchant Pays After 10 Days
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5900.00,
    "payment_mode": 70,
    "payment_type": "credit_settlement",
    "reference_number": "NEFT123456",
    "notes": "Credit payment after 10 days"
  }'
```

### Step 6: Vendor Confirms Payment
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/mark-paid" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Credit settled"
  }'
```

---

## 📦 **Scenario 3: Multiple Partial Payments**

### Payment 1: ₹1000 UPI
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000.00,
    "payment_mode": 30,
    "payment_type": "partial",
    "utr": "UTR001"
  }'
```

### Payment 2: ₹2000 Bank Transfer
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 2000.00,
    "payment_mode": 70,
    "payment_type": "partial",
    "reference_number": "NEFT789"
  }'
```

### Payment 3: ₹1500 Cash
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1500.00,
    "payment_mode": 10,
    "payment_type": "partial"
  }'
```

### Payment 4: ₹1400 Cheque
```bash
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1400.00,
    "payment_mode": 60,
    "payment_type": "final",
    "reference_number": "CHQ123456"
  }'
```

### Verify All at Once
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/mark-paid" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "All 4 payments verified and received"
  }'
```

**Expected:** All 4 payments marked as `verified: true`, `status: completed`

---

## 🧪 **Advanced: Individual Payment Verification**

### Verify Single Payment (Optional - Rarely Used)
```bash
PAYMENT_ID=1

curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/payments/${PAYMENT_ID}/verify" \
  -H "Content-Type: application/json" \
  -d '{
    "verified": true,
    "notes": "UTR confirmed in bank statement"
  }'
```

### Reject Payment
```bash
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/payments/${PAYMENT_ID}/verify" \
  -H "Content-Type: application/json" \
  -d '{
    "verified": false,
    "notes": "UTR not found - please resubmit"
  }'
```

---

## 📊 **Testing Checklist**

### ✅ Merchant Actions:
- [ ] Submit payment at PLACED status → status changes to PAYMENT_PENDING
- [ ] Submit payment at SHIPPED status → status stays SHIPPED
- [ ] Submit multiple payments → total_paid auto-calculates
- [ ] View payment history
- [ ] See payment summary by mode

### ✅ Vendor Actions:
- [ ] View all payments for an order
- [ ] Mark order as paid (bulk action) → all payments verified
- [ ] Order status updates from PAYMENT_PENDING → CONFIRMED
- [ ] Ship on credit (without payment verification)
- [ ] Individual payment verification (optional)

### ✅ System Behavior:
- [ ] Auto-calculation: paid_amount, outstanding_amount
- [ ] Order activity logs every action
- [ ] Payment independent of order flow
- [ ] Can pay at any order status
- [ ] Mix multiple payment modes in single order

---

## 🎯 **Quick Test Commands (Copy-Paste)**

```bash
# Set your IDs
export MERCHANT_ID=1
export VENDOR_ID=1
export ORDER_ID=1

# 1. Place order
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/place" \
  -H "Content-Type: application/json" -d '{}'

# 2. Submit payment
curl -X POST "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments" \
  -H "Content-Type: application/json" \
  -d '{"amount":500,"payment_mode":30,"payment_type":"token","utr":"UTR001"}'

# 3. View payments (Vendor)
curl -X GET "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/payments"

# 4. Mark as paid
curl -X POST "http://localhost:8080/api/v0/vendor/${VENDOR_ID}/orders/${ORDER_ID}/mark-paid" \
  -H "Content-Type: application/json" -d '{"notes":"Received"}'

# 5. View payments (Merchant)
curl -X GET "http://localhost:8080/api/v0/merchant/${MERCHANT_ID}/orders/${ORDER_ID}/payments"
```

---

## 🎉 **Success Indicators:**

1. ✅ Payment submission changes status to `PAYMENT_PENDING`
2. ✅ Mark as paid changes status to `CONFIRMED`
3. ✅ `paid_amount` and `outstanding_amount` auto-calculate correctly
4. ✅ Can mix multiple payment modes
5. ✅ Order can progress without payment (credit scenario)
6. ✅ Complete payment history visible
7. ✅ Every action logged in order_activity table

---

**🚀 Your flexible payment system is ready to use!**


