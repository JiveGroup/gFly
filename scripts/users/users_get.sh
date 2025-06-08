#!/bin/bash
#
# Function list all users data
# Method: GET
# Path: /users
#

source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi

# API URL
url="http://localhost:7889/api/v1/users"

# Parameter: keyword (Keyword)
keyword=""  # TODO: Set value for keyword

# Parameter: order_by (Order By)
order_by=""  # TODO: Set value for order_by

# Parameter: page (Page)
page=""  # TODO: Set value for page

# Parameter: per_page (Items Per Page)
per_page=""  # TODO: Set value for per_page

# Add query parameters
params=""
if [ -n "$keyword" ]; then
  params="?$paramskeyword=$keyword"
fi
if [ -n "$order_by" ]; then
  params="$params&order_by=$order_by"
fi
if [ -n "$page" ]; then
  params="$params&page=$page"
fi
if [ -n "$per_page" ]; then
  params="$params&per_page=$per_page"
fi

url="$url$params"

# Execute API call
echo "Calling API: $url"

response=$(curl -X "GET" "$url" \
  -H "Authorization: Bearer $TOKEN")

# Display response
echo "$response" | jq .

