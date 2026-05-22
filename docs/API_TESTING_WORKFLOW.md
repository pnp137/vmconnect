# Order Management API - Complete Testing Workflow

## Prerequisites
- Server running on `http://localhost:8080` (or your configured port)
- You have registered merchant and vendor accounts
- You have obtained JWT tokens for authentication

---

## Setup: Register Users & Get Tokens

### 1. Register a Merchant
```bash
curl -X POST http://localhost:8080/api/v0/merchant/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@shop.com",
    "phone": "9876543210",
    "password": "password123",
    "shop_name": "John Electronics",
    "address": "123 Main Street",
    "pincode": "110001"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Merchant registered successfully",
  "data": {
    "merchant_id": 1,
    "user_id": 1,
    "token": "eyJhbGc..."
  }
}
```

Save the `merchant_id` and `token`.

---

### 2. Register a Vendor
```bash
curl -X POST http://localhost:8080/api/v0/vendor/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@vendor.com",
    "phone": "9876543211",
    "password": "password123",
    "company_name": "Smith Suppliers",
    "gst_number": "29ABCDE1234F1Z5",
    "address": "456 Supplier Lane"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Vendor registered successfully",
  "data": {
    "vendor_id": 1,
    "user_id": 2,
    "vendor_code": "VS12345678",
    "token": "eyJhbGc..."
  }
}
```

Save the `vendor_id`, `vendor_code`, and `token`.

---

## Testing Flow: Complete Order Lifecycle

### Phase 1: Merchant Adds Items to Cart (IN_CART)

#### 3. Merchant Adds Vendor
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/vendors \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "vendor_code": "VS12345678"
  }'
```

#### 4. Add Product to Cart
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/vendors/1/cart/items \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 5
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "order_id": 1,
    "product_id": 1,
    "quantity": 5,
    "action": "added",
    "message": "Product added to cart"
  }
}
```

Save the `order_id`.

---

### Phase 2: Place Order (IN_CART → PLACED)

#### 5. Place Order
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/place \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Please deliver by Friday"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Order placed successfully",
  "data": {
    "order_id": 1,
    "status": "PLACED",
    "message": "Order placed successfully",
    "timestamp": "2025-10-26T10:30:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Merchant (User ID)
- Status: PLACED
- Metadata: `{"action": "place_order", "previous_status": "IN_CART"}`

---

### Phase 3: Vendor Confirms Order (PLACED → CONFIRMED)

#### 6. Vendor Confirms Order
```bash
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/confirm \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "estimated_delivery": "2025-11-05T18:00:00Z",
    "notes": "Order confirmed, will process today"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Order confirmed successfully",
  "data": {
    "order_id": 1,
    "status": "CONFIRMED",
    "message": "Order confirmed successfully",
    "timestamp": "2025-10-26T11:00:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Vendor (User ID)
- Status: CONFIRMED
- Metadata: `{"action": "confirm_order", "estimated_delivery": "...", "previous_status": "PLACED"}`

---

### Phase 4: Vendor Generates Invoice (CONFIRMED → INVOICED)

#### 7. Generate Invoice
```bash
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/invoice \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_no": "INV-2025-001",
    "invoice_date": "2025-10-26T12:00:00Z",
    "amount": 5000.00,
    "tax_amount": 900.00,
    "file_url": "https://s3.amazonaws.com/invoices/INV-2025-001.pdf",
    "notes": "Invoice generated and sent"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Invoice generated successfully",
  "data": {
    "order_id": 1,
    "invoice_id": 1,
    "invoice_no": "INV-2025-001",
    "status": "INVOICED",
    "message": "Invoice generated successfully",
    "timestamp": "2025-10-26T12:00:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Vendor (User ID)
- Status: INVOICED
- Metadata: `{"action": "generate_invoice", "invoice_info": {...}}`

---

### Phase 5: Merchant Submits Payment (INVOICED → PAID)

#### 8. Submit Payment
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/payment \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5900.00,
    "payment_mode": "upi",
    "utr": "UTR202510261234",
    "reference_number": "REF123456",
    "paid_date": "2025-10-26T13:00:00Z",
    "notes": "Payment completed via PhonePe"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Payment submitted successfully",
  "data": {
    "order_id": 1,
    "payment_id": 1,
    "status": "PAID",
    "amount": 5900.00,
    "message": "Payment submitted successfully",
    "timestamp": "2025-10-26T13:00:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Merchant (User ID)
- Status: PAID
- Metadata: `{"action": "submit_payment", "payment_info": {"utr": "...", "amount": ...}}`

---

### Phase 6: Vendor Verifies Payment

#### 9. Verify Payment
```bash
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/payment/verify \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": 1,
    "verified": true,
    "notes": "Payment received and verified in bank account"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Payment verified successfully",
  "data": {
    "order_id": 1,
    "payment_id": 1,
    "status": "PAID",
    "message": "Payment verified successfully",
    "timestamp": "2025-10-26T14:00:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Vendor (User ID)
- Status: PAID (no change)
- Metadata: `{"action": "verify_payment", "payment_id": 1, "verified": true}`

---

### Phase 7: Vendor Dispatches Order (PAID → SHIPPED)

#### 10. Dispatch Order
```bash
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/dispatch \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "tracking_id": "TRK202510261234",
    "carrier": "BlueDart",
    "dispatch_date": "2025-10-27T09:00:00Z",
    "expected_delivery": "2025-10-30T18:00:00Z",
    "packages": 2,
    "notes": "Dispatched via overnight service"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Order dispatched successfully",
  "data": {
    "order_id": 1,
    "tracking_id": "TRK202510261234",
    "status": "SHIPPED",
    "message": "Order dispatched successfully",
    "timestamp": "2025-10-27T09:00:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Vendor (User ID)
- Status: SHIPPED
- Metadata: `{"action": "dispatch_order", "dispatch_info": {"tracking_id": ..., "carrier": ...}}`

---

### Phase 8: Merchant Marks Received (SHIPPED → DELIVERED)

#### 11. Mark Order as Received
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/receive \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "received_date": "2025-10-30T16:30:00Z",
    "notes": "All items received in good condition"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Order marked as received successfully",
  "data": {
    "order_id": 1,
    "status": "DELIVERED",
    "message": "Order marked as received successfully",
    "timestamp": "2025-10-30T16:30:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Merchant (User ID)
- Status: DELIVERED
- Metadata: `{"action": "mark_received", "received_date": "..."}`

---

### Phase 9: Merchant Completes Order (DELIVERED → COMPLETED)

#### 12. Complete Order
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/complete \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 5,
    "notes": "Excellent service, will order again"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Order completed successfully",
  "data": {
    "order_id": 1,
    "status": "COMPLETED",
    "message": "Order completed successfully",
    "timestamp": "2025-10-30T17:00:00Z"
  }
}
```

**✅ OrderActivity Created:**
- Actor: Merchant (User ID)
- Status: COMPLETED
- Metadata: `{"action": "complete_order", "rating": 5}`

---

## Alternative Flows

### Flow 1: Cancel Order (Any State → CANCELLED)

#### By Merchant
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/cancel \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Changed requirements",
    "notes": "Will order next month"
  }'
```

#### By Vendor
```bash
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/cancel \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Out of stock",
    "refund_required": true,
    "notes": "Initiating refund process"
  }'
```

---

### Flow 2: Partial Payment
```bash
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/payment \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 3000.00,
    "payment_mode": "upi",
    "utr": "UTR202510261235",
    "paid_date": "2025-10-26T13:00:00Z",
    "notes": "Partial payment - 50% advance"
  }'
```

**This will set status to PARTIALLY_PAID**

---

### Flow 3: Flexible Order (Payment Before Dispatch)
```
1. PLACED → CONFIRMED (Vendor confirms)
2. CONFIRMED → PAID (Merchant pays directly)
3. PAID → SHIPPED (Vendor dispatches)
4. SHIPPED → DELIVERED (Merchant receives)
5. DELIVERED → COMPLETED (Merchant completes)
```

---

## Query & Inspection APIs

### Get Order Details
```bash
curl -X GET http://localhost:8080/api/v0/merchant/1/orders/1 \
  -H "Authorization: Bearer {MERCHANT_TOKEN}"
```

**Response includes OrderActivity history:**
```json
{
  "data": {
    "order": {
      "id": 1,
      "status": "COMPLETED",
      "total_amount": 5900.00,
      "orderActivity": [
        {
          "id": 1,
          "order_state": "PLACED",
          "actorId": 1,
          "actorRole": "merchant",
          "remarks": "Please deliver by Friday",
          "metadata": {"action": "place_order"},
          "createdAt": "2025-10-26T10:30:00Z"
        },
        {
          "id": 2,
          "order_state": "CONFIRMED",
          "actorId": 2,
          "actorRole": "vendor",
          "remarks": "Order confirmed",
          "createdAt": "2025-10-26T11:00:00Z"
        }
        // ... more activities
      ]
    }
  }
}
```

### List Orders
```bash
curl -X GET "http://localhost:8080/api/v0/merchant/1/orders?order_status=PLACED&order_status=CONFIRMED&limit=20&offset=0" \
  -H "Authorization: Bearer {MERCHANT_TOKEN}"
```

---

## Testing Idempotency

### Test 1: Try placing same order twice
```bash
# First call
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/place \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"notes": "First attempt"}'

# Second call (immediately after)
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/place \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Second attempt"}'
```

**Expected:** Second call returns the same result with message "Order was already placed"

---

## Error Scenarios

### 1. Invalid State Transition
```bash
# Try to deliver an order that's not shipped
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/receive \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response:**
```json
{
  "error": true,
  "message": "cannot transition from PLACED to DELIVERED"
}
```

### 2. Wrong Role Access
```bash
# Merchant trying to confirm order (only vendor can)
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/confirm \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response:**
```json
{
  "error": true,
  "message": "role merchant is not allowed to transition from PLACED to CONFIRMED"
}
```

---

## Summary: Complete Order Flow

```
IN_CART (Merchant adds items)
   ↓
PLACED (Merchant places order) ←─ POST /merchant/{mid}/orders/{oid}/place
   ↓
CONFIRMED (Vendor confirms) ←─────POST /vendor/{vid}/orders/{oid}/confirm
   ↓
INVOICED (Vendor generates) ←─────POST /vendor/{vid}/orders/{oid}/invoice
   ↓
PAID (Merchant pays) ←────────────POST /merchant/{mid}/orders/{oid}/payment
   ↓ [Optional: Vendor verifies] POST /vendor/{vid}/orders/{oid}/payment/verify
SHIPPED (Vendor dispatches) ←─────POST /vendor/{vid}/orders/{oid}/dispatch
   ↓
DELIVERED (Merchant receives) ←───POST /merchant/{mid}/orders/{oid}/receive
   ↓
COMPLETED (Merchant completes) ←──POST /merchant/{mid}/orders/{oid}/complete
```

**Each transition creates an OrderActivity record with full audit trail!** ✅

