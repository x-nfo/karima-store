#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
echo "Testing Modul 4: Pricing Engine & Shipping (Komerce Integration)..."

# Ensure python3 is available for JSON parsing
if ! command -v python3 &> /dev/null; then
    echo "Error: python3 is not installed."
    exit 1
fi

# Function to parse JSON using python
parse_json() {
    python3 -c "import sys, json; print(json.load(sys.stdin)$1)" 2>/dev/null
}

echo -e "\n1. Testing Shipping: Destination Search (Komerce API)..."
DEST=$(curl -s -X GET "$BASE_URL/shipping/destination/search?keyword=jakarta")
SUCCESS=$(echo "$DEST" | parse_json "['success']")
if [ "$SUCCESS" == "True" ]; then
    echo "✅ Destination Search Success"
    FIRST_DEST=$(echo "$DEST" | parse_json "['data'][0]['label']")
    FIRST_DEST_ID=$(echo "$DEST"  | parse_json "['data'][0]['id']")
    echo "   Sample Destination: $FIRST_DEST (ID: $FIRST_DEST_ID)"
else
    echo "❌ Destination Search Failed"
    echo "Response: $DEST"
fi

echo -e "\n2. Testing Shipping: Calculate Cost (Komerce API)..."
# Using example destination IDs from Komerce documentation
# shipper_destination_id=31597 (Example origin)
# receiver_destination_id=368 (Example destination)
COST_RES=$(curl -s -X GET "$BASE_URL/shipping/calculate?shipper_destination_id=31597&receiver_destination_id=368&weight=4&item_value=50000&cod=no")

SUCCESS=$(echo "$COST_RES" | parse_json "['success']")
if [ "$SUCCESS" == "True" ]; then
    echo "✅ Calculate Cost Success"
    # Try to extract first shipping option
    FIRST_SHIPPING=$(echo "$COST_RES" | parse_json "['data']['calculate_reguler'][0]['shipping_name']")
    FIRST_COST=$(echo "$COST_RES" | parse_json "['data']['calculate_reguler'][0]['shipping_cost']")
    echo "   Sample Courier: $FIRST_SHIPPING, Cost: $FIRST_COST"
else
    echo "❌ Calculate Cost Failed"
    echo "Response: $COST_RES"
fi

echo -e "\n3. Testing Pricing Engine: Calculate Price (Retail, Bulk Discount)..."
# Assuming Product ID 1 exists from previous tests
PRICE_PAYLOAD='{
  "product_id": 1,
  "quantity": 10,
  "customer_type": "retail"
}'

PRICE_RES=$(curl -s -X POST "$BASE_URL/pricing/calculate" \
  -H "Content-Type: application/json" \
  -d "$PRICE_PAYLOAD")

STATUS=$(echo "$PRICE_RES" | parse_json "['status']")

if [ "$STATUS" == "success" ]; then
    echo "✅ Calculate Price Success"
    FINAL_PRICE=$(echo "$PRICE_RES" | parse_json "['data']['final_price']")
    DISCOUNT_TYPE=$(echo "$PRICE_RES" | parse_json "['data']['discount_type']")
    echo "   Final Price: $FINAL_PRICE"
    echo "   Discount Type Applied: $DISCOUNT_TYPE"
else
    echo "⚠️ Calculate Price Failed (Product might not exist or other error)"
    echo "Response: $PRICE_RES"
fi

echo -e "\nTest Complete."
