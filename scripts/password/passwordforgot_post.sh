#!/bin/bash
#
# Forgot password
# Method: POST
# Path: /password/forgot
#

# API URL
url="http://localhost:7889/api/v1/password/forgot"

# Request body for request.ForgotPassword
request_body='{
  "username": "john@jivecode.com"
}'

# Execute API call
echo "Calling API: $url"

response=$(curl -X "POST" "$url" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d "$request_body")

# Display response
echo "$response" | jq .

