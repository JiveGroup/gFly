#!/bin/bash
#
# Reset password
# Method: POST
# Path: /password/reset
#

# API URL
url="http://localhost:7889/api/v1/password/reset"

# Request body for request.ResetPassword
request_body='{
  "password": "M1PassW@s"
,
  "token": "293r823or832eioj2eo9282o423"
}'

# Execute API call
echo "Calling API: $url"

response=$(curl -X "POST" "$url" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d "$request_body")

# Display response
echo "$response" | jq .

