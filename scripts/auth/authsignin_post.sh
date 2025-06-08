#!/bin/bash
#
# authenticating user's credentials
# Method: POST
# Path: /auth/signin
#

# API URL
url="http://localhost:7889/api/v1/auth/signin"

# Request body for request.SignIn
request_body='{
  "password": "P@seWor9"
,
  "username": "admin@gfly.dev"
}'

# Execute API call
echo "Calling API: $url"

response=$(curl -X "POST" "$url" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d "$request_body")

# Display response
echo "$response" | jq .

