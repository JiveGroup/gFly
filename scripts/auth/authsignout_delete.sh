#!/bin/bash
#
# de-authorize user and delete refresh token from Redis
# Method: DELETE
# Path: /auth/signout
#

source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi

# API URL
url="http://localhost:7889/api/v1/auth/signout"

