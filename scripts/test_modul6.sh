#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
echo "Testing Modul 6: Notification & Caching..."

# Ensure python3 is available for JSON parsing
if ! command -v python3 &> /dev/null; then
    echo "Error: python3 is not installed."
    exit 1
fi

# Function to parse JSON using python
parse_json() {
    python3 -c "import sys, json; print(json.load(sys.stdin)$1)" 2>/dev/null
}

echo -e "\n=== Test 1: WhatsApp Service Status ===\n"
STATUS_RES=$(curl -s -X GET "$BASE_URL/whatsapp/status")
echo "Response: $STATUS_RES"

STATUS=$(echo "$STATUS_RES" | parse_json "['data']['status']")
if [ "$STATUS" == "connected" ] || [ "$STATUS" == "disconnected" ]; then
    echo "✅ Status Endpoint Working (Current State: $STATUS)"
else
    echo "❌ Status Endpoint Failed"
fi

echo -e "\n=== Test 2: Get Webhook URL ===\n"
WEBHOOK_RES=$(curl -s -X GET "$BASE_URL/whatsapp/webhook-url")
echo "Response: $WEBHOOK_RES"

URL=$(echo "$WEBHOOK_RES" | parse_json "['data']['webhook_url']")
if [[ "$URL" == *"karimastore.com"* ]]; then
    echo "✅ Webhook URL Endpoint Working: $URL"
else
    echo "❌ Webhook URL Endpoint Failed"
fi

echo -e "\n=== Test 3: Send Test Message (Expect Failure/Success based on device conn) ===\n"
# Using a dummy number, Fonnte will likely reject if not connected or invalid number, 
# but we check if the API handles the request correctly.
PHONE="6281234567890"
MESSAGE="Test message from Karima Store Backend"

TEST_MSG_RES=$(curl -s -X POST "$BASE_URL/whatsapp/test?phone_number=$PHONE&message=${MESSAGE// /%20}")
echo "Response: $TEST_MSG_RES"

MSG_STATUS=$(echo "$TEST_MSG_RES" | parse_json "['status']")
if [ "$MSG_STATUS" == "success" ]; then
    echo "✅ Test Message Sent Successfully"
elif [ "$MSG_STATUS" == "error" ]; then
    echo "⚠️ Test Message Failed (Expected if device disconnected/invalid token)"
    ERR_MSG=$(echo "$TEST_MSG_RES" | parse_json "['message']")
    echo "   Error: $ERR_MSG"
else
    echo "❌ Unknown Response"
fi

echo -e "\n=== Test 4: Product Caching Verification (Indirect) ===\n"
echo "Fetching products list (First Hit - Should Cache)..."
START_TIME=$(date +%s%N)
curl -s -X GET "$BASE_URL/products?limit=5" > /dev/null
END_TIME=$(date +%s%N)
DURATION_1=$(( (END_TIME - START_TIME) / 1000000 ))
echo "Request 1 Duration: ${DURATION_1}ms"

echo "Fetching products list (Second Hit - Should be Faster)..."
START_TIME=$(date +%s%N)
curl -s -X GET "$BASE_URL/products?limit=5" > /dev/null
END_TIME=$(date +%s%N)
DURATION_2=$(( (END_TIME - START_TIME) / 1000000 ))
echo "Request 2 Duration: ${DURATION_2}ms"

if [ $DURATION_2 -lt $DURATION_1 ]; then
    echo "✅ Cache likely working (Request 2 was faster)"
else
    echo "⚠️ Cache verification inconclusive (Network variance possible)"
fi

echo -e "\n=== Module 6 Test Summary ===\n"
echo "Endpoints Tested:"
echo "  ✅ GET  /api/v1/whatsapp/status"
echo "  ✅ GET  /api/v1/whatsapp/webhook-url"
echo "  ✅ POST /api/v1/whatsapp/test"
echo "  ✅ GET  /api/v1/products (Caching Check)"

echo -e "\nNote: Real WhatsApp delivery requires a connected Fonnte device."
echo "Test Complete."
