#!/bin/bash
#
# Sign up a new user
# Method: POST
# Path: /auth/signup
#

# API URL
url="http://localhost:7889/api/v1/auth/signup"

# Request body for request.SignUp
request_body='{
  "avatar": "https://i.pravatar.cc/32"
,
  "email": "john@jivecode.com"
,
  "fullname": "John Doe"
,
  "password": "M1PassW@s"
,
  "phone": "0989831911"
,
  "status": "pending"
}'

# Execute API call
echo "Calling API: $url"

response=$(curl -X "POST" "$url" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d "$request_body")

# Display response
echo "$response" | jq .

