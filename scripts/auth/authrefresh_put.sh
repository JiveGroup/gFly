#!/bin/bash
#
# refresh user token
# Method: PUT
# Path: /auth/refresh
#

source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi

# API URL
url="http://localhost:7889/api/v1/auth/refresh"

# Request body for request.RefreshToken
request_body='{
  "token": "d1a4216a226cbf75eaefc9107c2c64b6b2c0f18cd8634e3a6f495146c38e1324.1747914602"
}'

# Execute API call
echo "Calling API: $url"

response=$(curl -X "PUT" "$url" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d "$request_body")

# Display response
echo "$response" | jq .

