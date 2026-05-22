# Quick Start Testing Guide

## 🚀 Quick Test (5 Minutes)

### Step 1: Start Server
```bash
cd /Users/admin/Documents/gochat/Blockchain/vmconnect
go run main.go
```

### Step 2: Register Test Users

**Register Merchant:**
```bash
curl -X POST http://localhost:8080/api/v0/merchant/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Merchant",
    "email": "merchant@test.com",
    "phone": "9999999999",
    "password": "test123",
    "shop_name": "Test Shop",
    "address": "123 Test St",
    "pincode": "110001"
  }'
```

Save the token from response as `MERCHANT_TOKEN`

**Register Vendor:**
```bash
curl -X POST http://localhost:8080/api/v0/vendor/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Vendor",
    "email": "vendor@test.com",
    "phone": "9999999998",
    "password": "test123",
    "company_name": "Test Suppliers",
    "gst_number": "29TEST1234F1Z5",
    "address": "456 Vendor St"
  }'
```

Save the token as `VENDOR_TOKEN` and `vendor_code`

### Step 3: Test Complete Flow

Replace `{MERCHANT_TOKEN}`, `{VENDOR_TOKEN}`, `{VENDOR_CODE}` with actual values:

```bash
# 1. Merchant adds vendor
curl -X POST http://localhost:8080/api/v0/merchant/1/vendors \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"vendor_code": "{VENDOR_CODE}"}'

# 2. Add item to cart (assuming product_id 1 exists)
curl -X POST http://localhost:8080/api/v0/merchant/1/vendors/1/cart/items \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 5}'

# 3. Place order (order_id from step 2 response)
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/place \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Test order"}'

# 4. Vendor confirms
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/confirm \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Confirmed"}'

# 5. Vendor generates invoice
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/invoice \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_no": "INV-TEST-001",
    "invoice_date": "2025-10-26T12:00:00Z",
    "amount": 5000.00,
    "tax_amount": 900.00,
    "notes": "Test invoice"
  }'

# 6. Merchant pays
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/payment \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5900.00,
    "payment_mode": "upi",
    "utr": "UTR-TEST-123",
    "paid_date": "2025-10-26T13:00:00Z",
    "notes": "Test payment"
  }'

# 7. Vendor dispatches
curl -X POST http://localhost:8080/api/v0/vendor/1/orders/1/dispatch \
  -H "Authorization: Bearer {VENDOR_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "tracking_id": "TRK-TEST-123",
    "carrier": "Test Carrier",
    "notes": "Test dispatch"
  }'

# 8. Merchant receives
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/receive \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Test received"}'

# 9. Merchant completes
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/complete \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"rating": 5, "notes": "Test completed"}'

# 10. Get order with activity history
curl -X GET http://localhost:8080/api/v0/merchant/1/orders/1 \
  -H "Authorization: Bearer {MERCHANT_TOKEN}"
```

### Step 4: Verify in Database

```sql
-- Check order status
SELECT id, status FROM orders WHERE id = 1;

-- Check activity history (should have 8-9 records)
SELECT 
  id,
  order_id,
  actor_id,
  actor_role,
  order_state,
  remarks,
  metadata,
  created_at
FROM order_activities
WHERE order_id = 1
ORDER BY created_at ASC;
```

## ✅ Expected Results

After running all commands, you should see:

1. **Order Status:** `COMPLETED` (100)
2. **OrderActivity Records:** 8-9 entries showing full history
3. **Each Activity Has:**
   - Actor ID (user who performed action)
   - Actor Role (merchant/vendor)
   - Order State at that time
   - Metadata with action details
   - Timestamp

## 🧪 Test Idempotency

Try running step 3 (place order) twice:
```bash
# First call - should succeed
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/place \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"notes": "First"}'

# Second call - should return "Order was already placed"
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/place \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Second"}'
```

## 🐛 Test Error Handling

Try invalid transition:
```bash
# Try to deliver before shipping
curl -X POST http://localhost:8080/api/v0/merchant/1/orders/1/receive \
  -H "Authorization: Bearer {MERCHANT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{}'

# Expected: Error - "cannot transition from PLACED to DELIVERED"
```

## 📚 Full Documentation

- **Complete Guide:** `API_TESTING_WORKFLOW.md`
- **Implementation Details:** `IMPLEMENTATION_SUMMARY.md`

## 🎉 You're All Set!

The order state management system is fully implemented and ready to use!

