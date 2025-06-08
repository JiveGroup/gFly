#!/bin/bash
#
# Generate API shell scripts from Swagger JSON file
#

# Configuration
SWAGGER_FILE="docs/swagger.json"
OUTPUT_DIR="scripts"
BASE_URL="http://localhost:7889/api/v1"
AUTH_USERNAME="admin@thietngon.com"
AUTH_PASSWORD="P@seWor9"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Function to get authentication token
get_auth_token() {
  echo '#!/bin/bash
#
# Authentication script to get token
#

echo "# ============================== Login system =============================="
url="'$BASE_URL'/auth/signin"

json=$(curl -X "POST" "$url" \
        -H "Content-Type: application/json; charset=utf-8" \
        -d '"'"'{
          "username": "'$AUTH_USERNAME'",
          "password": "'$AUTH_PASSWORD'"
}'"'"')

# Extract token
token=$(echo $json | jq -r ".access")

echo "Token: $token"
echo "export TOKEN=$token" > '$OUTPUT_DIR'/token.env
' > "$OUTPUT_DIR/auth.sh"
  chmod +x "$OUTPUT_DIR/auth.sh"
}

# Parse Swagger JSON and generate scripts
parse_swagger() {
  # Extract paths and methods
  paths=$(jq -r '.paths | keys[]' "$SWAGGER_FILE")

  for path in $paths; do
    # Get methods for this path (GET, POST, PUT, DELETE, etc.)
    methods=$(jq -r ".paths[\"$path\"] | keys[]" "$SWAGGER_FILE")

    for method in $methods; do
      # Skip non-HTTP methods
      if [[ "$method" != "get" && "$method" != "post" && "$method" != "put" && "$method" != "delete" && "$method" != "patch" ]]; then
        continue
      fi

      # Get summary and tags
      summary=$(jq -r ".paths[\"$path\"].$method.summary" "$SWAGGER_FILE")
      tags=$(jq -r ".paths[\"$path\"].$method.tags[0]" "$SWAGGER_FILE")

      # Create script name from path and method
      script_name=$(echo "${path}_${method}" | tr -d '/' | tr -d '{' | tr -d '}' | tr '[:upper:]' '[:lower:]')

      # Check if authentication is required
      requires_auth=$(jq -r ".paths[\"$path\"].$method.security != null" "$SWAGGER_FILE")

      # Get parameters
      parameters=$(jq -r ".paths[\"$path\"].$method.parameters" "$SWAGGER_FILE" 2>/dev/null)

      # Create script
      create_script "$path" "$method" "$summary" "$tags" "$script_name" "$requires_auth" "$parameters"
    done
  done
}

# Create individual API script
create_script() {
  local path="$1"
  local method="$2"
  local summary="$3"
  local tags="$4"
  local script_name="$5"
  local requires_auth="$6"
  local parameters="$7"

  # Create directory for tag if it doesn't exist
  if [ "$tags" != "null" ]; then
    # Convert to lowercase and replace spaces with underscores
    tag_dir=$(echo "$tags" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')
    mkdir -p "$OUTPUT_DIR/$tag_dir"
    script_path="$OUTPUT_DIR/$tag_dir/$script_name.sh"
  else
    script_path="$OUTPUT_DIR/$script_name.sh"
  fi

  # Start script content
  # Convert method to uppercase
  method_upper=$(echo "$method" | tr '[:lower:]' '[:upper:]')

  echo "#!/bin/bash
#
# $summary
# Method: $method_upper
# Path: $path
#
" > "$script_path"

  # Add source for token if authentication is required
  if [ "$requires_auth" == "true" ]; then
    echo 'source $(dirname "$0")/../token.env

if [ -z "$TOKEN" ]; then
  echo "No authentication token found. Please run auth.sh first."
  exit 1
fi
' >> "$script_path"
  fi

  # Add URL construction
  echo "# API URL
url=\"$BASE_URL$path\"
" >> "$script_path"

  # Add parameter handling
  add_parameters "$script_path" "$path" "$method" "$parameters"

  # Make script executable
  chmod +x "$script_path"

  echo "Generated script: $script_path"
}

# Add parameter handling to script
add_parameters() {
  local script_path="$1"
  local path="$2"
  local method="$3"
  local parameters="$4"

  # Extract query, path, and body parameters
  query_params=()
  path_params=()
  body_param=""
  body_schema_ref=""

  # Check if parameters is null or empty
  if [ "$parameters" == "null" ] || [ -z "$parameters" ]; then
    return
  fi

  # Parse parameters from JSON
  while read -r param_line; do
    if [ -n "$param_line" ] && [ "$param_line" != "null" ]; then
      param_name=$(echo "$param_line" | jq -r '.name // "null"')
      param_in=$(echo "$param_line" | jq -r '.in // "null"')
      param_type=$(echo "$param_line" | jq -r '.type // "object"')
      param_desc=$(echo "$param_line" | jq -r '.description // ""')

      if [ "$param_name" != "null" ] && [ -n "$param_name" ]; then
        if [ "$param_in" == "query" ]; then
          query_params+=("$param_name")
          # Get example value if available
          param_example=$(echo "$param_line" | jq -r '.example // ""')
          echo "# Parameter: $param_name ($param_desc)
${param_name}=\"$param_example\"  # TODO: Set value for $param_name
" >> "$script_path"
        elif [ "$param_in" == "path" ]; then
          path_params+=("$param_name")
          # Get example value if available
          param_example=$(echo "$param_line" | jq -r '.example // ""')
          echo "# Parameter: $param_name ($param_desc)
${param_name}=\"$param_example\"  # TODO: Set value for $param_name
" >> "$script_path"
        elif [ "$param_in" == "body" ]; then
          body_param="$param_name"
          # Extract schema reference
          body_schema_ref=$(echo "$param_line" | jq -r '.schema | if has("$ref") then ."$ref" else "" end')

          # If we have a schema reference, extract the definition name
          if [ -n "$body_schema_ref" ]; then
            # Extract definition name from reference (e.g., "#/definitions/request.CreateAddress" -> "request.CreateAddress")
            definition_name=${body_schema_ref#"#/definitions/"}

            # Get the definition from swagger.json
            definition=$(jq -r ".definitions[\"$definition_name\"]" "$SWAGGER_FILE")

            # Generate request body from definition
            generate_request_body "$script_path" "$definition" "$definition_name"
          else
            # Fallback to default request body if no schema reference
            echo "# Request body
request_body='{
  \"example\": \"value\"  # TODO: Set proper request body
}'
" >> "$script_path"
          fi
        fi
      fi
    fi
  done < <(jq -c '.[]?' <<<"$parameters" 2>/dev/null)

  # Replace path parameters in URL
  for param in "${path_params[@]}"; do
    echo "# Replace path parameter in URL
url=\${url/\{$param\}/\$$param}
" >> "$script_path"
  done

  # Add query parameters to URL
  if [ ${#query_params[@]} -gt 0 ]; then
    echo "# Add query parameters
params=\"\"" >> "$script_path"

    for i in "${!query_params[@]}"; do
      param="${query_params[$i]}"
      if [ $i -eq 0 ]; then
        echo "if [ -n \"\$$param\" ]; then
  params=\"?\$params$param=\$$param\"
fi" >> "$script_path"
      else
        echo "if [ -n \"\$$param\" ]; then
  params=\"\$params&$param=\$$param\"
fi" >> "$script_path"
      fi
    done

    echo "
url=\"\$url\$params\"
" >> "$script_path"
  fi

  # Add curl command
  echo "# Execute API call
echo \"Calling API: \$url\"
" >> "$script_path"

  # Build curl command based on method and parameters
  curl_cmd="curl -X \"$method_upper\" \"\$url\""

  # Add authentication header if required
  if grep -q "TOKEN" "$script_path"; then
    curl_cmd="$curl_cmd \\
  -H \"Authorization: Bearer \$TOKEN\""
  fi

  # Add content type header for methods with body
  if [ -n "$body_param" ]; then
    curl_cmd="$curl_cmd \\
  -H \"Content-Type: application/json; charset=utf-8\" \\
  -d \"\$request_body\""
  fi

  # Add response handling
  echo "response=\$($curl_cmd)

# Display response
echo \"\$response\" | jq .
" >> "$script_path"
}

# Function to generate request body from definition
generate_request_body() {
  local script_path="$1"
  local definition="$2"
  local definition_name="$3"

  # Start building the request body
  echo "# Request body for $definition_name" >> "$script_path"
  echo "request_body='{" >> "$script_path"

  # Get properties from definition
  properties=$(echo "$definition" | jq -r '.properties // {}')

  # Get required fields
  required_fields=$(echo "$definition" | jq -r '.required // []')

  # Process each property
  first_prop=true
  while read -r prop_name; do
    if [ -n "$prop_name" ] && [ "$prop_name" != "null" ]; then
      # Get property details
      prop_details=$(echo "$properties" | jq -r ".[\"$prop_name\"]")
      prop_type=$(echo "$prop_details" | jq -r '.type // "string"')
      prop_example=$(echo "$prop_details" | jq -r '.example // ""')

      # Check if property is required
      is_required=$(echo "$required_fields" | jq -r "contains([\"$prop_name\"])")

      # Add comma if not first property
      if [ "$first_prop" = true ]; then
        first_prop=false
      else
        echo "," >> "$script_path"
      fi

      # Add property to request body with example value
      if [ "$prop_type" == "string" ]; then
        echo "  \"$prop_name\": \"$prop_example\"" >> "$script_path"
      elif [ "$prop_type" == "integer" ] || [ "$prop_type" == "number" ]; then
        # Use 0 as default for numeric types if example is empty
        if [ -z "$prop_example" ]; then
          echo "  \"$prop_name\": 0" >> "$script_path"
        else
          echo "  \"$prop_name\": $prop_example" >> "$script_path"
        fi
      elif [ "$prop_type" == "boolean" ]; then
        # Use false as default for boolean types if example is empty
        if [ -z "$prop_example" ]; then
          echo "  \"$prop_name\": false" >> "$script_path"
        else
          echo "  \"$prop_name\": $prop_example" >> "$script_path"
        fi
      elif [ "$prop_type" == "array" ]; then
        echo "  \"$prop_name\": []" >> "$script_path"
      elif [ "$prop_type" == "object" ]; then
        echo "  \"$prop_name\": {}" >> "$script_path"
      else
        echo "  \"$prop_name\": null" >> "$script_path"
      fi
    fi
  done < <(echo "$properties" | jq -r 'keys[]')

  # Close the request body
  echo "}'
" >> "$script_path"
}

# Main execution
echo "Generating API scripts from $SWAGGER_FILE..."
get_auth_token
parse_swagger
echo "Done! Scripts generated in $OUTPUT_DIR"
