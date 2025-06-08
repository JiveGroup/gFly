#!/bin/bash
#
# Get user profile
# Method: GET
# Path: /users/profile
#

source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi

# API URL
url="http://localhost:7889/api/v1/users/profile"

