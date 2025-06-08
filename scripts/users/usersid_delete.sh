#!/bin/bash
#
# Delete user by given userID
# Method: DELETE
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

# Replace path parameter in URL
url=${url/\{id\}/$id}

# Execute API call
echo "Calling API: $url"

response=$(curl -X "DELETE" "$url" \
  -H "Authorization: Bearer $TOKEN")

# Display response
echo "$response" | jq .

