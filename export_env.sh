#!/bin/bash

# This script exports variables from .env file to the current shell
# Usage: source ./export_env.sh

if [ ! -f .env ]; then
  echo "Error: .env file not found"
  echo "Create one by copying .env.example:"
  echo "cp .env.example .env"
  return 1 2>/dev/null || exit 1
fi

# Read .env file line by line
while IFS= read -r line; do
  # Skip empty lines and comments
  [[ -z "$line" || "$line" == \#* ]] && continue
  
  # Split line on the first equals sign
  key=$(echo "$line" | cut -d '=' -f 1)
  value=$(echo "$line" | cut -d '=' -f 2-)
  
  # Remove quotes if present
  value=$(echo "$value" | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//")
  
  # Export the variable
  if [ -n "$key" ]; then
    eval "export $key=\"$value\""
    echo "Exported: $key=$value"
  fi
done < .env

echo "Environment variables from .env have been exported."
echo "Run 'env | grep STRAVA' to verify."