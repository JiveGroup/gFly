#!/bin/bash
#
# Function allows Administrator update an existing user
# Method: PUT
# Path: /users/{id}
#

source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi

# API URL
url="http://localhost:7889/api/v1/users/{id}"

# Parameter: id (User ID)
id=""  # TODO: Set value for id

# Request body for request.UpdateUser
request_body='{
  "avatar": "https://i.pravatar.cc/32"
,
  "fullname": "John Doe"
,
  "password": "M1PassW@s"
,
  "phone": "0989831911"
,
  "roles": []
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

