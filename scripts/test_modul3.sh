#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
echo "Starting Modul 3 Tests..."

# 1. Categories
echo -e "\n[1] Testing Categories..."
echo "GET /categories"
curl -s "$BASE_URL/categories"
echo ""

# 2. Upload (Mock) / Media
echo "test image content" > test_image.jpg

# 3. Products
echo -e "\n[2] Testing Products..."
TIMESTAMP=$(date +%s)
echo "Creating Product 'Baju Koko Modern $TIMESTAMP'..."
CREATE_RES=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Baju Koko Modern $TIMESTAMP\",
    \"sku\": \"SKU-$TIMESTAMP\",
    \"description\": \"Baju koko premium\",
    \"price\": 150000,
    \"category\": \"tops\",
    \"stock\": 100,
    \"brand\": \"Karima\",
    \"material\": \"Cotton\",
    \"weight\": 0.5
  }")
echo $CREATE_RES

PID=$(echo $CREATE_RES | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])")
echo "Created Product ID: $PID"

if [ "$PID" == "null" ] || [ -z "$PID" ]; then
    echo "Failed to create product. Exiting."
    exit 1
fi

echo -e "\n[3] Testing Image Upload..."
UPLOAD_RES=$(curl -s -X POST "$BASE_URL/products/$PID/media" \
  -F "product_id=$PID" \
  -F "file=@test_image.jpg" \
  -F "is_primary=true")
echo $UPLOAD_RES

echo -e "\n[4] Testing Variants..."
echo "Creating Variant 'Baju Koko Modern - L - Putih'..."
VAR_RES=$(curl -s -X POST "$BASE_URL/variants" \
  -H "Content-Type: application/json" \
  -d "{
    \"product_id\": $PID,
    \"name\": \"Baju Koko Modern - L - Putih\",
    \"size\": \"L\",
    \"color\": \"Putih\",
    \"price\": 150000,
    \"stock\": 50,
    \"sku\": \"BAJ-$TIMESTAMP-L-PUT\"
  }")
echo $VAR_RES

echo -e "\n[5] Testing Stock Update..."
STOCK_RES=$(curl -s -X PATCH "$BASE_URL/products/$PID/stock" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 10}')
echo $STOCK_RES

echo -e "\n[6] Verifying Product Details (Media & Stock)..."
curl -s "$BASE_URL/products/$PID"
echo ""

# Cleanup
rm test_image.jpg
