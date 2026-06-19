#!/bin/bash

BASE_URL="http://localhost:8080/v1"
EMAIL="testuser$(date +%s)@example.com"
PASSWORD="securepassword"

echo "=== Testing POST /v1/auth/register ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X POST "$BASE_URL/auth/register" \
     -H "Content-Type: application/json" \
     -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" > register_resp.txt
cat register_resp.txt
echo ""

echo "=== Testing POST /v1/auth/login ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X POST "$BASE_URL/auth/login" \
     -H "Content-Type: application/json" \
     -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" -c cookies.txt > login_resp.txt
cat login_resp.txt
echo ""

echo "=== Testing GET /v1/products ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X GET "$BASE_URL/products" > products_resp.txt
cat products_resp.txt
echo ""
