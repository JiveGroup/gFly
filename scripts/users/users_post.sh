#!/bin/bash
#
# Create a new user for Administrator
# Method: POST
# Path: /users
#

source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi

# API URL
url="http://localhost:7889/api/v1/users"

# Request body for request.CreateUser
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
  "roles": []
,
  "status": "pending"
}'

# Execute API call
echo "Calling API: $url"

response=$(curl -X "POST" "$url" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d "$request_body")

# Display response
echo "$response" | jq .

