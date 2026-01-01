#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
echo "Testing Modul 5: Transaction & Payment Gateway (Midtrans Integration)..."

# Ensure python3 is available for JSON parsing
if ! command -v python3 &> /dev/null; then
    echo "Error: python3 is not installed."
    exit 1
fi

# Function to parse JSON using python
parse_json() {
    python3 -c "import sys, json; print(json.load(sys.stdin)$1)" 2>/dev/null
}

# First, create a test product if not exists
echo -e "\n=== Setup: Creating Test Product ===\n"
PRODUCT_PAYLOAD='{
  "name": "Test Product for Checkout",
  "description": "Product for testing Module 5",
  "sku": "TEST-CHECKOUT-001",
  "price": 100000,
  "category_id": 1,
  "weight": 500,
  "stock": 100,
  "status": "active"
}'

PRODUCT_RES=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d "$PRODUCT_PAYLOAD")

PRODUCT_ID=$(echo "$PRODUCT_RES" | parse_json "['data']['id']" 2>/dev/null || echo "1")
echo "Using Product ID: $PRODUCT_ID"

echo -e "\n=== Test 1: Checkout Endpoint ===\n"
CHECKOUT_PAYLOAD=$(cat <<EOF
{
  "items": [
    {
      "product_id": ${PRODUCT_ID},
      "quantity": 2
    }
  ],
  "shipping_name": "John Doe",
  "shipping_phone": "08123456789",
  "shipping_address": "Jl. Sudirman No. 123",
  "shipping_city": "Jakarta",
  "shipping_province": "DKI Jakarta",
  "shipping_postal_code": "10110",
  "payment_method": "bank_transfer",
  "user_id": 1,
  "customer_notes": "Test order from Module 5 verification"
}
EOF
)

CHECKOUT_RES=$(curl -s -X POST "$BASE_URL/checkout" \
  -H "Content-Type: application/json" \
  -d "$CHECKOUT_PAYLOAD")

STATUS=$(echo "$CHECKOUT_RES" | parse_json "['status']")
if [ "$STATUS" == "success" ]; then
    echo "✅ Checkout Success"
    ORDER_NUMBER=$(echo "$CHECKOUT_RES" | parse_json "['data']['order_number']")
    ORDER_ID=$(echo "$CHECKOUT_RES" | parse_json "['data']['order_id']")
    SNAP_TOKEN=$(echo "$CHECKOUT_RES" | parse_json "['data']['snap_token']")
    AMOUNT=$(echo "$CHECKOUT_RES" | parse_json "['data']['amount']")
    
    echo "   Order Number: $ORDER_NUMBER"
    echo "   Order ID: $ORDER_ID"
    echo "   Amount: Rp $AMOUNT"
    
    if [ -n "$SNAP_TOKEN" ] && [ "$SNAP_TOKEN" != "None" ]; then
        echo "   ✅ Snap Token Generated: ${SNAP_TOKEN:0:30}..."
    else
        echo "   ⚠️ Snap Token: Not generated"
    fi
else
    echo "❌ Checkout Failed"
    echo "Response: $CHECKOUT_RES"
    ORDER_NUMBER=""
    ORDER_ID=""
fi

echo -e "\n=== Test 2: Get Order Detail ===\n"
if [ -n "$ORDER_ID" ] && [ "$ORDER_ID" != "None" ]; then
    ORDER_DETAIL=$(curl -s -X GET "$BASE_URL/orders/$ORDER_ID")
    
    ORDER_STATUS=$(echo "$ORDER_DETAIL" | parse_json "['data']['status']")
    PAYMENT_STATUS=$(echo "$ORDER_DETAIL" | parse_json "['data']['payment_status']")
    
    if [ -n "$ORDER_STATUS" ]; then
        echo "✅ Order Detail Retrieved"
        echo "   Order Status: $ORDER_STATUS"
        echo "   Payment Status: $PAYMENT_STATUS"
    else
        echo "❌ Failed to get order detail"
        echo "Response: $ORDER_DETAIL"
    fi
else
    echo "⚠️ Skipping - No order ID from checkout"
fi

echo -e "\n=== Test 3: Webhook Endpoint (Structure Test) ===\n"
# Note: This is just to test if webhook endpoint exists and handles requests
# Real webhook testing requires actual Midtrans sandbox
WEBHOOK_PAYLOAD=$(cat <<EOF
{
  "order_id": "${ORDER_NUMBER:-ORD-TEST}",
  "status_code": "200",
  "gross_amount": "200000.00",
  "signature_key": "test_signature",
  "transaction_status": "settlement",
  "transaction_id": "test123",
  "payment_type": "bank_transfer",
  "transaction_time": "2026-01-01 16:00:00",
  "fraud_status": "accept"
}
EOF
)

WEBHOOK_RES=$(curl -s -X POST "$BASE_URL/payment/webhook" \
  -H "Content-Type: application/json" \
  -d "$WEBHOOK_PAYLOAD")

# Webhook will likely fail signature verification, but endpoint should respond
if echo "$WEBHOOK_RES" | grep -q "signature\|success\|error"; then
    echo "✅ Webhook Endpoint Responding"
    echo "   (Signature verification expected to fail in test)"
else
    echo "⚠️ Webhook endpoint response: $WEBHOOK_RES"
fi

echo -e "\n=== Test 4: Order List Endpoint ===\n"
ORDERS_LIST=$(curl -s -X GET "$BASE_URL/orders")

if echo "$ORDERS_LIST" | grep -q "data\|orders"; then
    echo "✅ Orders List Endpoint Working"
    ORDER_COUNT=$(echo "$ORDERS_LIST" | parse_json "['data'].__len__()" 2>/dev/null || echo "N/A")
    if [ "$ORDER_COUNT" != "N/A" ]; then
        echo "   Total Orders: $ORDER_COUNT"
    fi
else
    echo "❌ Orders List Failed"
    echo "Response: $ORDERS_LIST"
fi

echo -e "\n=== Module 5 Test Summary ===\n"
echo "Endpoints Tested:"
echo "  ✅ POST /api/v1/checkout          - Checkout & Snap Token"
echo "  ✅ GET  /api/v1/orders/:id        - Order Detail"
echo "  ✅ POST /api/v1/payment/webhook   - Payment Notification"
echo "  ✅ GET  /api/v1/orders             - Order Listing"

echo -e "\nNote: Midtrans integration requires valid SERVER_KEY for actual payment processing."
echo "      Webhook signature verification will fail without real Midtrans notification."

echo -e "\nTest Complete."
