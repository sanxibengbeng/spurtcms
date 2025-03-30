#!/bin/bash

# Script to fix logger import paths in the codebase
# This script will replace relative imports with proper module imports

echo "Starting module import fixes..."

# Get the module name from go.mod
MODULE_NAME=$(grep "^module" go.mod | awk '{print $2}')
echo "Module name: $MODULE_NAME"

# Find all Go files
find . -name "*.go" -type f > /tmp/all_go_files.txt

# Process each file
total_files=$(wc -l < /tmp/all_go_files.txt)
processed=0

while IFS= read -r file; do
  processed=$((processed + 1))
  echo "[$processed/$total_files] Checking imports in $file"
  
  # Replace relative imports with module imports
  if grep -q "\"\\./logger\"" "$file"; then
    echo "  Fixing ./logger import in $file"
    sed -i '' "s|\"\\./logger\"|\"$MODULE_NAME/logger\"|g" "$file"
  fi
  
  if grep -q "\"\\.\\.\/logger\"" "$file"; then
    echo "  Fixing ../logger import in $file"
    sed -i '' "s|\"\\.\\.\/logger\"|\"$MODULE_NAME/logger\"|g" "$file"
  fi
  
  if grep -q "\"\\.\\.\/\\.\\.\/logger\"" "$file"; then
    echo "  Fixing ../../logger import in $file"
    sed -i '' "s|\"\\.\\.\/\\.\\.\/logger\"|\"$MODULE_NAME/logger\"|g" "$file"
  fi
  
  if grep -q "\"\\.\\.\/\\.\\.\/\\.\\.\/logger\"" "$file"; then
    echo "  Fixing ../../../logger import in $file"
    sed -i '' "s|\"\\.\\.\/\\.\\.\/\\.\\.\/logger\"|\"$MODULE_NAME/logger\"|g" "$file"
  fi
  
done < /tmp/all_go_files.txt

echo "Module import fixes complete. Processed $total_files files."
echo "Please verify the changes and ensure all imports are correctly resolved."
