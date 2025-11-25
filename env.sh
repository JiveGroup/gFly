#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    # Read .env file line by line, skip comments and empty lines
    while IFS= read -r line || [ -n "$line" ]; do
        # Skip empty lines and comments
        if [[ -z "$line" || "$line" =~ ^[[:space:]]*# ]]; then
            continue
        fi

        # Check if line contains a valid variable assignment (KEY=VALUE)
        if [[ "$line" =~ ^[a-zA-Z_][a-zA-Z0-9_]*= ]]; then
            # Export the variable
            export "$line"
        fi
    done < .env
fi
