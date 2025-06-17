#!/bin/bash
#
# Authentication script to get token
#

echo "# ============================== Login system =============================="
url="http://localhost:7889/api/v1/auth/signin"

json=$(curl -X "POST" "$url" \
        -H "Content-Type: application/json; charset=utf-8" \
        -d '{
          "username": "admin@gfly.dev",
          "password": "P@seWor9"
}')

# Extract token
token=$(echo $json | jq -r ".access")

echo "Token: $token"
echo "export TOKEN=$token" > scripts/token.env

