#!/bin/bash
#
# Update user's status by ID
# Method: PUT
# Path: /users/{id}/status
#

source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi

# API URL
url="http://localhost:7889/api/v1/users/{id}/status"

# Parameter: id (User ID)
id=""  # TODO: Set value for id

# Request body for request.UpdateUserStatus
request_body='{
  "status": "active"
}'

# Replace path parameter in URL
url=${url/\{id\}/$id}

# Execute API call
echo "Calling API: $url"

response=$(curl -X "PUT" "$url" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d "$request_body")

# Display response
echo "$response" | jq .

